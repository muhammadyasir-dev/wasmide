import React, { useState } from "react";
import {
  Backdrop,
  Box,
  Button,
  Fade,
  FormControl,
  InputLabel,
  MenuItem,
  Modal,
  Select,
  Snackbar,
  TextField,
  Typography,
} from "@mui/material";
import { useNavigate } from "react-router-dom";

const LandingPage: React.FC = () => {
  const [open, setOpen] = useState<boolean>(false);
  const [projectName, setProjectName] = useState<string>("");
  const [language, setLanguage] = useState<string>("");
  const [snackbarOpen, setSnackbarOpen] = useState<boolean>(false);
  const navigate = useNavigate();

  const handleOpen = () => setOpen(true);
  const handleClose = () => setOpen(false);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (projectName && language) {
      navigate(`/editor/?project=${projectName}&lang=${language}`);
    } else {
      setSnackbarOpen(true);
    }
  };

  const handleSnackbarClose = () => {
    setSnackbarOpen(false);
  };

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
        overflow: "hidden",
        width: "100vw",
        bgcolor: "lightgrey",
        p: 4,
      }}
    >
      <Box
        sx={{
          textAlign: "center",
          bgcolor: "white",
          borderRadius: 2,
          boxShadow: "0 4px 30px rgba(255, 255, 255, 0.2)",
          p: 4,
          width: "20rem",
        }}
      >
        <Typography
          variant="h4"
          component="h1"
          gutterBottom
          sx={{
            color: "#333",
          }}
        >
          Create Project
        </Typography>
        <Button
          variant="contained"
          color="default"
          onClick={handleOpen}
          sx={{
            borderRadius: "20px",
            boxShadow: "0 4px 20px rgba(255, 255, 255, 0.2)",
            bgcolor: "#555",
            color: "#fff",
            "&:hover": {
              boxShadow: "0 6px 30px rgba(255, 255, 255, 0.3)",
              bgcolor: "#777",
            },
            mb: 2, // Add margin bottom for spacing
          }}
        >
          Create New Project
        </Button>
        <Modal
          open={open}
          onClose={handleClose}
          closeAfterTransition
          BackdropComponent={Backdrop}
          BackdropProps={{
            timeout: 100,
          }}
        >
          <Fade in={open}>
            <Box
              sx={{
                bgcolor: "background.paper",
                borderRadius: 2,
                boxShadow: "0 4px 30px rgba(255, 255, 255, 0.2)",
                p: 4,
                width: "400px",
                mx: "auto",
                mt: "20vh",
                backdropFilter: "blur(10px)",
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
              }}
            >
              <Typography
                variant="h5"
                component="h2"
                gutterBottom
                sx={{ color: "#333" }}
              >
                Create a New Project
              </Typography>
              <form onSubmit={handleSubmit} style={{ width: "100%" }}>
                <TextField
                  label="Project Name"
                  variant="outlined"
                  fullWidth
                  margin="normal"
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value)}
                  required
                  sx={{
                    borderRadius: "10px",
                    "& .MuiOutlinedInput-root": {
                      "& fieldset": {
                        borderRadius: "10px",
                        borderColor: "#aaa",
                      },
                      "&:hover fieldset": {
                        borderColor: "#888",
                      },
                    },
                  }}
                />
                <FormControl fullWidth margin="normal" required>
                  <InputLabel sx={{ color: "#555" }}>
                    Programming Language
                  </InputLabel>
                  <Select
                    value={language}
                    onChange={(e) => setLanguage(e.target.value)}
                    label="Programming Language"
                    sx={{
                      borderRadius: "10px",
                      "& .MuiOutlinedInput-root": {
                        "& fieldset": {
                          borderRadius: "10px",
                          borderColor: "#aaa",
                        },
                      },
                    }}
                  >
                    <MenuItem value="c">C</MenuItem>
                    <MenuItem value="c++">C++</MenuItem>
                    <MenuItem value="rust">Rust</MenuItem>
                    <MenuItem value="golang">Go</MenuItem>
                    <MenuItem value="zig">Zig</MenuItem>
                  </Select>
                </FormControl>
                <Button
                  type="submit"
                  variant="contained"
                  color="default"
                  sx={{
                    mt: 2,
                    borderRadius: "20px",
                    boxShadow: "0 4px 20px rgba(255, 255, 255, 0.2)",
                    bgcolor: "#555",
                    color: "#fff",
                    "&:hover": {
                      boxShadow: "0 6px 30px rgba(255, 255, 255, 0.3)",
                      bgcolor: "#777",
                    },
                  }}
                >
                  Create Workspace
                </Button>
              </form>
            </Box>
          </Fade>
        </Modal>
        <Snackbar
          open={snackbarOpen}
          autoHideDuration={6000}
          onClose={handleSnackbarClose}
          message="Please fill in all fields"
        />
      </Box>
    </Box>
  );
};

export default LandingPage;