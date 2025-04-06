import React, { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import Iframe from "../components/iframe";
import {
  AppBar,
  Box,
  Container,
  IconButton,
  Toolbar,
  Typography,
} from "@mui/material";
import { CoffeeRounded, Folder, PlayArrow } from "@mui/icons-material";
import FileTree from "../components/FileTree";
import PTerminal from "../components/PTerminal";
const CodeEditor: React.FC = () => {
  const [code, setCode] = useState("//");
  const [highlightedCode, setHighlightedCode] = useState("");
  const [showoutput, setshowoutput] = useState(false);
  const navigate = useNavigate();
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const lineNumbersRef = useRef<HTMLDivElement>(null);
  const editorWrapperRef = useRef<HTMLDivElement>(null);

  const highlightSyntax = (code: string) => {
    const keywords =
      /\b(const|let|var|function|if|else|return|for|while|switch|case|break|default|try|catch|finally|throw|async|await|class|extends|super|import|export|from|this|new|delete|instanceof|typeof|void|in|of)\b/g;
    const stringLiterals = /(['"`].*?['"`])/g;
    const comments = /(\/\/.*?$|\/\*[\s\S]*?\*\/)/gm;
    return code
      .replace(comments, '<span style="color: green;">$&</span>')
      .replace(stringLiterals, '<span style="color: orange;">$&</span>')
      .replace(
        keywords,
        '<span style="color: blue; font-weight: bold;">$&</span>',
      );
  };

  const updateLineNumbers = () => {
    if (textareaRef.current && lineNumbersRef.current) {
      const lines = textareaRef.current.value.split("\n");
      const lineNumbers = Array.from({ length: lines.length }, (_, i) => i + 1)
        .map((num) => `<div class="line-number">${num}</div>`)
        .join("");
      lineNumbersRef.current.innerHTML = lineNumbers;

      // Sync scroll position
      if (editorWrapperRef.current) {
        lineNumbersRef.current.scrollTop = editorWrapperRef.current.scrollTop;
      }
    }
  };

  const handleInput = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newCode = e.target.value;
    setCode(newCode);
    updateLineNumbers();
  };

  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    if (lineNumbersRef.current) {
      lineNumbersRef.current.scrollTop = e.currentTarget.scrollTop;
    }
  };

  // Handle tab key
  //
  const handleoutput = () => {
    setshowoutput(true);
  };
  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Tab") {
      e.preventDefault();
      const start = e.currentTarget.selectionStart;
      const end = e.currentTarget.selectionEnd;
      const newCode = code.substring(0, start) + "    " + code.substring(end);
      setCode(newCode);

      // Reset cursor position
      setTimeout(() => {
        if (textareaRef.current) {
          textareaRef.current.selectionStart =
            textareaRef.current
              .selectionEnd =
            start + 4;
        }
      }, 0);
    }
  };

  useEffect(() => {
    updateLineNumbers();
    setHighlightedCode(highlightSyntax(code));
  }, [code]);

  useEffect(() => {
    updateLineNumbers();
  }, []);

  return (
    <>
      {!showoutput
        ? (
          <Container
            maxWidth={false}
            sx={{
              height: "100%",
              display: "flex",
              flexDirection: "row",
              padding: 0,
            }}
          >
            <AppBar position="absolute" sx={{ backgroundColor: "black" }}>
              <Toolbar>
                <IconButton edge="start" color="inherit" aria-label="menu">
                  <CoffeeRounded />
                </IconButton>
                <Typography variant="h6" sx={{ flexGrow: 1 }}>
                  Wasm IDE
                </Typography>
                <button onClick={() => handleCreateFile}>
                  +
                </button>
                <IconButton onClick={handleoutput} color="inherit">
                  <PlayArrow />
                </IconButton>
              </Toolbar>
            </AppBar>
            <Box sx={{ display: "flex", alignItems: "flex" }}>
              <Box
                sx={{
                  height: "auto",
                  width: "5rem",
                  marginRight: "1rem",
                  bgcolor: "silver",
                  scrollBehavior: "initial",
                  color: "lightgrey",
                  padding: "1px",
                  margin: "0",
                  boxSizing: "border-box",
                  border: "none",
                }}
              >
                <FileTree />
              </Box>
              <Box sx={{ boxSizing: "border-box" }}>
                <Box
                  sx={{ display: "flex", flexGrow: 1, position: "relative" }}
                >
                  <div
                    ref={lineNumbersRef}
                    style={{
                      width: "1rem",
                      padding: "1rem 0.5rem",
                      backgroundColor: "blueviolet",
                      color: "#666",
                      textAlign: "right",
                      userSelect: "none",
                      borderRight: "1px solid #3d3d3d",
                      fontFamily: "monospace",
                      fontSize: "14px",
                      lineHeight: "1.5",
                      whiteSpace: "pre",
                    }}
                  />
                </Box>
                <PTerminal />
              </Box>
            </Box>
          </Container>
        )
        : (
          <>
            <Iframe />
          </>
        )}
    </>
  );
};

export default CodeEditor;
