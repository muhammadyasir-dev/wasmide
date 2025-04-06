import { useEffect, useState } from "react";
import {
  AppBar,
  Box,
  Container,
  IconButton,
  Toolbar,
  Typography,
} from "@mui/material";
import { backendurl } from "../libs/Url";


function Login() {
  const [user, setUser] = useState<boolean>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<undefined | null>(null);

  // Check if we've been redirected back from OAuth provider
  useEffect(() => {
    if (window.location.pathname === "/") {
      setLoading(true);
      // Fetch user data from our backend
      fetch("http://localhost:8080/auth/user", {
        credentials: "include", // Important to include cookies
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Failed to get user data");
          }
          return response.json();
        })
        .then((userData) => {
          setUser(userData);
          console.log("User ID:", userData.id);
          console.log("User Email:", userData.email);
          console.log("User Name:", userData.name);
        })
        .catch((err) => {
          setError(err.message);
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, []);

  const handleLogin = () => {
    window.location.href = `${backendurl}/auth/login`;
  };

  return (
    <>
      <center>
        <div className="oauth-container">
          <h2>Login Now</h2>

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


              <div className="textbox">
                <input type="email" placeholder="Email" />
                <input type="password" placeholder="Password" />
              </div>

              {loading && <p>Loading...</p>}

              {error && <p className="error">Error: {error}</p>}

              <button onClick={handleLogin}>
                Login with Google
              </button>
            </div>
          </center>
        </>
        );
}

        export default Login;
