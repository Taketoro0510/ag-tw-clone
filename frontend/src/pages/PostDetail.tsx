import React from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Box, IconButton, Typography, CircularProgress, Paper, Button } from "@mui/material";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import { useQuery } from "@tanstack/react-query";
import { fetchApi } from "../api/client";
import { PostCard } from "../features/posts/PostCard";
import type { paths } from "../api/types";

type Post = paths["/posts/{id}"]["get"]["responses"]["200"]["content"]["application/json"];

export const PostDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: post, isLoading, error } = useQuery<Post>({
    queryKey: ["post", id],
    queryFn: () => fetchApi<Post>(`/posts/${id}`),
    enabled: !!id,
  });

  if (isLoading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error || !post) {
    return (
      <Box sx={{ py: 4, textAlign: "center" }}>
        <Typography variant="h6" color="error" gutterBottom>
          Post not found or failed to load.
        </Typography>
        <Button variant="contained" onClick={() => navigate("/")} sx={{ mt: 2, borderRadius: 5 }}>
          Back to Timeline
        </Button>
      </Box>
    );
  }

  return (
    <Box>
      <Box sx={{ display: "flex", alignItems: "center", mb: 3 }}>
        <IconButton onClick={() => navigate(-1)} sx={{ mr: 2 }}>
          <ArrowBackIcon />
        </IconButton>
        <Typography variant="h5" sx={{ fontWeight: "bold" }}>
          Post
        </Typography>
      </Box>
      <Paper elevation={0} sx={{ border: "1px solid", borderColor: "divider", borderRadius: 4, overflow: "hidden" }}>
        <PostCard post={post} defaultShowComments={true} />
      </Paper>
    </Box>
  );
};
