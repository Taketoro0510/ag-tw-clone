import React from "react";
import { Button, Typography, Box, Paper } from "@mui/material";
import { useAuth } from "../features/auth/AuthContext";
import { useNavigate, useLocation } from "react-router-dom";

export const Login: React.FC = () => {
  const { loginWithGoogle, user, token } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const from = location.state?.from?.pathname || "/";

  React.useEffect(() => {
    if (user && token) {
      navigate(from, { replace: true });
    }
  }, [user, token, navigate, from]);

  const handleLogin = async () => {
    try {
      await loginWithGoogle();
      navigate(from, { replace: true });
    } catch (err) {
      console.error(err);
      alert("Login failed. Please try again.");
    }
  };

  return (
    <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", minHeight: "100vh", bgcolor: "background.default" }}>
      <Paper elevation={3} sx={{ p: 4, textAlign: "center", borderRadius: 2 }}>
        <Typography variant="h4" component="h1" gutterBottom sx={{ fontWeight: "bold" }}>
          CloudCode SNS
        </Typography>
        <Typography variant="body1" color="textSecondary" sx={{ mb: 4 }}>
          Sign in to share your thoughts with the world.
        </Typography>
        <Button variant="contained" color="primary" size="large" onClick={handleLogin}>
          Sign in with Google
        </Button>
      </Paper>
    </Box>
  );
};
