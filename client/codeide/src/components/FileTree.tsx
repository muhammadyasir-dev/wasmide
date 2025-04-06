import React, { useEffect, useState } from "react";
import { fileserver } from "../libs/Url";
import { colors } from "@mui/material";
const styles = {
  container: {
    display: "flex",
    height: "100vw",
    fontFamily: "Arial, sans-serif",
  },
  sidebar: {
    width: "250px",
    borderRight: "0px solid #ccc",
    padding: "0px",
  },
  mainContent: {
    flex: 1,
    padding: "20px",
  },
  header: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    margin: ".5rem 0",
  },
  title: {
    margin: 0,
    fontSize: "1.2em",
    fontWeight: "bold",
  },
  button: {
    padding: ".3rem .7rem",
    backgroundColor: "#007bff",
    color: "white",
    border: "none",
    borderRadius: "1px",
    cursor: "pointer",
  },
  fileList: {
    listStyle: "none",
    padding: 0,
    margin: 0,
    color: "red",
  },
  fileItem: {
    padding: "4px",
    cursor: "pointer",
    marginBottom: "2px",
    color: "black",
  },
  selectedFile: {
    fontWeight: "bold",
    color: "white",
    backgroundColor: "brown",
  },
  textarea: {
    caretColor: "red",
    "positiion": "absolute",
    width: "50rem",
    height: "calc(100vh - 300px)",
    padding: "1rem",
    border: "none",
    borderRadius: "4px",
    resize: "none",
    fontSize: "1rem",
    fontFamily: "monospace",
    cursor: "arrow",
    color: "white",
    background: "black",
  },
  error: {
    color: "red",
    padding: "10px",
  },
  loading: {
    textAlign: "center",
    padding: "20px",
    color: "black",
  },
};

const FileTree: React.FC = () => {
  const [files, setFiles] = useState<string[]>([]);
  const [selectedFile, setSelectedFile] = useState<string | null>(null);
  const [fileContent, setFileContent] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // Fetch the list of files when the component mounts
  useEffect(() => {
    const fetchFiles = async () => {
      try {
        const response = await fetch(
          `${fileserver}/list-files?project=myproject`,
        );
        if (!response.ok) {
          throw new Error("Failed to fetch files");
        }
        const fileList = await response.json();
        setFiles(fileList);
      } catch (error) {
        setError("No files availible");
        console.error("Error:", error);
      }
    };

    fetchFiles();
  }, []);

  // Load the content of the selected file
  useEffect(() => {
    const loadFileContent = async () => {
      if (!selectedFile) return;

      setIsLoading(true);
      setError(null);

      try {
        const response = await fetch(
          `${fileserver}/files/${selectedFile}?project=myproject`,
        );
        if (!response.ok) {
          throw new Error(`Failed to fetch ${selectedFile}`);
        }
        const result = await response.json();
        setFileContent(result.content || "");
      } catch (error) {
        setError("Failed to load file content");
        console.error("Error:", error);
      } finally {
        setIsLoading(false);
      }
    };

    loadFileContent();
  }, [selectedFile]);

  // Handle file selection
  const handleFileClick = (fileName: string) => {
    setSelectedFile(fileName);
    setError(null);
  };

  // Handle changes in the textarea
  const handleChange = async (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value;
    setFileContent(newContent);

    try {
      const response = await fetch(
        `${fileserver}/files/${selectedFile}?project=myproject`,
        {
          method: "POST",
          headers: {
            "Content-Type": "text/plain",
          },
          body: newContent,
        },
      );

      if (!response.ok) {
        throw new Error("Failed to save changes");
      }
    } catch (error) {
      setError("Failed to save changes");
      console.error("Error:", error);
    }
  };

  // Handle file creation
  const handleCreateFile = async () => {
    const fileName = prompt("");
    if (!fileName) return;

    try {
      const response = await fetch(
        `${fileserver}/create-file?project=myproject`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(fileName),
        },
      );

      if (!response.ok) {
        throw new Error("Failed to create file");
      }

      // Refresh the file list after creating a new file
      const listResponse = await fetch(
        `${fileserver}/list-files?project=myproject`,
      );
      if (listResponse.ok) {
        const fileList = await listResponse.json();
        setFiles(fileList);
      }
    } catch (error) {
      setError("Failed to create file");
      console.error("Error:", error);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.sidebar}>
        <div style={styles.header}>
          {
            /* <button style={styles.button} onClick={handleCreateFile}>
            +
          </button> */
          }
        </div>
        <ul style={styles.fileList}>
          {files.map((file) => (
            <li
              key={file}
              onClick={() => handleFileClick(file)}
              style={{
                ...styles.fileItem,
                ...(selectedFile === file ? styles.selectedFile : {}),
              }}
            >
              {file}
            </li>
          ))}
        </ul>
      </div>

      <div style={styles.mainContent}>
        {isLoading
          ? <div style={styles.loading}>Loading...</div>
          : error
          ? <div style={styles.error}>{error}</div>
          : selectedFile
          ? (
            <div>
              <h3 style={styles.title}>{selectedFile}</h3>
              <textarea
                value={fileContent}
                onChange={handleChange}
                style={styles.textarea}
                placeholder=""
              />
            </div>
          )
          : (
            <div style={styles.loading}>
              Wasm ide
            </div>
          )}
      </div>
    </div>
  );
};

export default FileTree;
