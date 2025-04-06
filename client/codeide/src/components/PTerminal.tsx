import React, { useEffect, useRef, useState } from "react";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "xterm/css/xterm.css";
import { backendurl } from "../libs/Url";

const PTerminal: React.FC = () => {
  const terminalRef = useRef<HTMLDivElement | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const xterm = useRef<Terminal | null>(null);
  const fitAddon = useRef<FitAddon>(new FitAddon());
  const inputBuffer = useRef<string>("");

  // Update the backend URL to match the Golang backend

  const getQueryParams = () => {
    const params = new URLSearchParams(window.location.search);
    return {
      project: params.get("project") || "default-project",
    };
  };

  const executeCommand = async (command: string) => {
    const { project } = getQueryParams();

    try {
      const response = await fetch(
        `${backendurl}/execute?project=${project}`,
        {
          method: "POST",
          body: command,
          headers: {
            "Content-Type": "text/plain",
          },
        },
      );

      if (!response.ok) {
        const errorText = await response.text();
        xterm.current?.writeln(`\r\nError: ${errorText}`);
        return;
      }

      const output = await response.text();
      xterm.current?.writeln(`\r\n${output}`);
    } catch (error) {
      xterm.current?.writeln(
        `\r\nError: ${error instanceof Error ? error.message : "Unknown error"
        }`,
      );
    }

    xterm.current?.write("$ ");
  };

  const scrollToBottom = () => {
    if (containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
    }
  };

  useEffect(() => {
    if (terminalRef.current) {
      xterm.current = new Terminal({
        cursorBlink: true,
        theme: {
          background: "#1e1e1e",
          foreground: "#ffffff",
          cursor: "#ffffff",
        },
        fontFamily: "Courier New, monospace",
        fontSize: 14,
        convertEol: true,
        scrollback: 1000,
      });

      xterm.current.loadAddon(fitAddon.current);
      xterm.current.open(terminalRef.current);
      fitAddon.current.fit();

      // Initial prompt
      xterm.current.write("$ ");

      xterm.current.onData((data) => {
        if (data.charCodeAt(0) === 13) { // Enter
          const command = inputBuffer.current.trim();
          if (command) {
            xterm.current?.writeln(""); // Move to next line
            executeCommand(command);
          } else {
            xterm.current?.write("\r\n$ ");
          }
          inputBuffer.current = "";
        } else if (data.charCodeAt(0) === 127) { // Backspace
          if (inputBuffer.current.length > 0) {
            inputBuffer.current = inputBuffer.current.slice(0, -1);
            xterm.current?.write("\b \b");
          }
        } else {
          inputBuffer.current += data;
          xterm.current?.write(data);
        }
      });

      const resizeObserver = new ResizeObserver(() => {
        fitAddon.current.fit();
        scrollToBottom();
      });

      if (containerRef.current) {
        resizeObserver.observe(containerRef.current);
      }

      return () => {
        resizeObserver.disconnect();
        xterm.current?.dispose();
      };
    }
  }, []);

  return (
    <div
      ref={containerRef}
      className="terminal-container"
      style={{
        position: "fixed",
        bottom: 0,
        left: 0,
        width: "100vw",
        height: "30vh",
        backgroundColor: "rgba(30, 30, 30, 0.95)",
        padding: "1rem",
        zIndex: 1000,
        overflow: "auto",
        display: "flex",
        flexDirection: "column",
      }}
    >
      <div
        ref={terminalRef}
        style={{
          flex: 1,
          minHeight: 0,
          width: "100%",
        }}
      />
    </div>
  );
};

export default PTerminal;
