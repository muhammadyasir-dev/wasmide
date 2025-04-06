import React, { useState } from "react";
import { Box, Button, Fade, Modal, TextField, Typography } from "@mui/material";
import { backendurl } from "../libs/Url";
const style = {
    position: "absolute",
    top: "50%",
    left: "50%",
    transform: "translate(-50%, -50%)",
    width: 400,
    bgcolor: "background.paper",
    boxShadow: 24, // Material-UI shadow
    p: 4,
    borderRadius: "8px",
};

const Signup: React.FC = () => {
    const [open, setOpen] = useState(true); // Open modal by default
    const [formData, setFormData] = useState({
        email: "",
        password: "",
    });
    const [error, setError] = useState<string | null>(null);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData({ ...formData, [name]: value });
    };

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        // Basic validation
        if (!formData.email || !formData.password) {
            setError("Email and password are required.");
            return;
        }
        console.log(formData);
        // Reset error if submission is successful
        setError(null);
        const signupdata = async () => {
            const signupdata = await fetch(`${backendurl}`, {
                methods: "POST",
                headers: {
                    "Content-Type": "appplication/json",
                },
                credentials: "include",
                body: JSON.stringiufy()
            });
        };
        signupdata();
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
                onClose={() => setOpen(false)} // Allow closing the modal
                closeAfterTransition
                BackdropProps={{
                    timeout: 500,
                }}
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
                            Signup
                        </Typography>
                        {error && (
                            <Typography
                                color="error"
                                style={{ textAlign: "center" }}
                            >
                                {error}
                            </Typography>
                        )}
                        <form onSubmit={handleSubmit}>
                            <TextField
                                fullWidth
                                label="Email"
                                name="email"
                                variant="outlined"
                                margin="normal"
                                onChange={handleChange}
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
                                required
                            />
                            <Button
                                type="submit"
                                variant="contained"
                                style={{
                                    backgroundColor: "#007aff",
                                    color: "#fff",
                                    marginTop: "20px",
                                }}
                            >
                                Sign Up
                            </Button>
                        </form>
                    </Box>
                </Fade>
            </Modal>
        </div>
    );
};

export default Signup;

