import React, { useState } from "react";
import { backendurl } from "../libs/Url";
import {
  Alert,
  Box,
  Button,
  Fade,
  Modal,
  TextField,
  Typography,
} from "@mui/material";
import GoogleIcon from "@mui/icons-material/Google";

const style = {
  position: "absolute",
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
  width: 400,
  bgcolor: "background.paper",
  border: "2px solid #000",
  boxShadow: 24,
  p: 4,
  borderRadius: "8px",
};

interface LoginFormData {
  email: string;
  password: string;
}

const Login: React.FC = () => {
  const [open, setOpen] = useState(true);
  const [formData, setFormData] = useState<LoginFormData>({
    email: "",
    password: "",
  });
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
    setError(""); // Clear error when user types
  };

  const handleRegularLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include", // Important for cookies
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message || "Login failed");
      }

      // Login successful
      window.location.href = "/home"; // Redirect to home page
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  };

  const handleGoogleLogin = () => {
    // Redirect to Google OAuth endpoint
    window.location.href = `${backendurl}/auth/login`;
  };

  return (
    <div
      style={{
        backgroundColor: "#f0f0f0",
        height: "100vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <Modal
        open={open}
        onClose={() => { }}
        closeAfterTransition
        BackdropProps={{
          timeout: 500,
        }}
        disableEscapeKeyDown
      >
        <Fade in={open}>
          <Box sx={style}>
            <Typography
              variant="h6"
              component="h2"
              style={{
                marginBottom: "20px",
                textAlign: "center",
              }}
            >
              Log In
            </Typography>

            {error && (
              <Alert
                severity="error"
                style={{ marginBottom: "20px" }}
              >
                {error}
              </Alert>
            )}

            <form onSubmit={handleRegularLogin}>
              <TextField
                fullWidth
                label="Email"
                name="email"
                variant="outlined"
                margin="normal"
                onChange={handleChange}
                disabled={loading}
                required
              />
              <TextField
                fullWidth
                label="Password"
                name="password"
                type="password"
                variant="outlined"
                margin="normal"
                onChange={handleChange}
                disabled={loading}
                required
              />

              <Button
                type="submit"
                variant="contained"
                fullWidth
                disabled={loading}
                style={{
                  backgroundColor: "#007aff",
                  color: "#fff",
                  marginTop: "20px",
                }}
              >
                {loading ? "Logging in..." : "Log In"}
              </Button>

              <Box
                sx={{
                  display: "flex",
                  alignItems: "center",
                  margin: "20px 0",
                }}
              >
                <Box
                  sx={{
                    flex: 1,
                    borderBottom: "1px solid #ccc",
                  }}
                />
                <Typography
                  sx={{
                    margin: "0 10px",
                    color: "#666",
                  }}
                >
                  or
                </Typography>
                <Box
                  sx={{
                    flex: 1,
                    borderBottom: "1px solid #ccc",
                  }}
                />
              </Box>

              <Button
                onClick={handleGoogleLogin}
                variant="outlined"
                fullWidth
                startIcon={<GoogleIcon />}
                disabled={loading}
                style={{ marginTop: "10px" }}
              >
                Login with Google
              </Button>
            </form>
          </Box>
        </Fade>
      </Modal>
    </div>
  );
};

export default Login;
