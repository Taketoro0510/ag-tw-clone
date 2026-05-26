import React, { useState } from "react";
import { Box, Typography, Avatar, IconButton, TextField, Button, CircularProgress, Badge, Divider } from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import FavoriteIcon from "@mui/icons-material/Favorite";
import FavoriteBorderIcon from "@mui/icons-material/FavoriteBorder";
import BookmarkIcon from "@mui/icons-material/Bookmark";
import BookmarkBorderIcon from "@mui/icons-material/BookmarkBorder";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchApi } from "../../api/client";
import type { paths } from "../../api/types";
import { useAuth } from "../auth/AuthContext";
import { useNavigate } from "react-router-dom";

type Comment = NonNullable<paths["/posts/{id}/comments"]["get"]["responses"]["200"]["content"]["application/json"]["items"]>[0];

const CommentItem: React.FC<{ comment: Comment; postId: string }> = ({ comment, postId }) => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const isMyComment = user && comment.author?.firebaseUid === user.uid;
  const authorName = comment.author?.displayName || (comment.authorId ? `User ${comment.authorId.substring(0, 5)}` : "Unknown");
  const authorAvatar = comment.author?.avatarUrl;

  const likeMutation = useMutation({
    mutationFn: () => fetchApi(`/comments/${comment.id}/likes`, { method: comment.likedByMe ? "DELETE" : "POST" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["comments", postId] });
    }
  });

  const bookmarkMutation = useMutation({
    mutationFn: () => fetchApi(`/comments/${comment.id}/bookmarks`, { method: comment.bookmarkedByMe ? "DELETE" : "POST" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["comments", postId] });
    }
  });

  const deleteMutation = useMutation({
    mutationFn: () => fetchApi(`/posts/${postId}/comments/${comment.id}`, { method: "DELETE" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["comments", postId] });
    }
  });

  const handleDelete = () => {
    if (window.confirm("Are you sure you want to delete this comment?")) {
      deleteMutation.mutate();
    }
  };

  return (
    <Box sx={{ display: "flex", mb: 2, alignItems: "flex-start" }}>
      <Avatar 
        sx={{ width: 32, height: 32, mr: 1.5, bgcolor: "secondary.main", fontSize: "0.875rem", cursor: "pointer" }} 
        src={authorAvatar || undefined}
        onClick={() => navigate(`/users/${comment.authorId}`)}
      >
        {!authorAvatar ? authorName.substring(0, 2).toUpperCase() : ""}
      </Avatar>
      <Box sx={{ flexGrow: 1, bgcolor: "background.default", p: 1.5, borderRadius: 2 }}>
        <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 0.5 }}>
          <Typography 
            variant="subtitle2" 
            sx={{ fontWeight: "bold", cursor: "pointer", "&:hover": { textDecoration: "underline" } }}
            onClick={() => navigate(`/users/${comment.authorId}`)}
          >
            {authorName}
          </Typography>
          <Typography variant="caption" color="textSecondary">
            {comment.createdAt ? new Date(comment.createdAt).toLocaleString() : ""}
          </Typography>
        </Box>
        <Typography variant="body2" sx={{ whiteSpace: "pre-wrap", mb: 1 }}>
          {comment.body}
        </Typography>
        <Box sx={{ display: "flex", gap: 1, alignItems: "center", ml: -1 }}>
          <IconButton size="small" onClick={() => likeMutation.mutate()} color={comment.likedByMe ? "error" : "default"}>
            <Badge badgeContent={comment.likeCount} color="error" sx={{ "& .MuiBadge-badge": { fontSize: "0.6rem", minWidth: "14px", height: "14px" } }}>
              {comment.likedByMe ? <FavoriteIcon fontSize="small" /> : <FavoriteBorderIcon fontSize="small" />}
            </Badge>
          </IconButton>
          <IconButton size="small" onClick={() => bookmarkMutation.mutate()} color={comment.bookmarkedByMe ? "primary" : "default"}>
            <Badge badgeContent={comment.bookmarkCount} color="primary" sx={{ "& .MuiBadge-badge": { fontSize: "0.6rem", minWidth: "14px", height: "14px" } }}>
              {comment.bookmarkedByMe ? <BookmarkIcon fontSize="small" /> : <BookmarkBorderIcon fontSize="small" />}
            </Badge>
          </IconButton>
          {isMyComment && (
            <IconButton size="small" color="error" onClick={handleDelete} sx={{ ml: "auto" }}>
              <DeleteIcon fontSize="small" />
            </IconButton>
          )}
        </Box>
      </Box>
    </Box>
  );
};

export const CommentSection: React.FC<{ postId: string }> = ({ postId }) => {
  const [newComment, setNewComment] = useState("");
  const queryClient = useQueryClient();

  const { data: response, isLoading } = useQuery({
    queryKey: ["comments", postId],
    queryFn: () => fetchApi(`/posts/${postId}/comments`) as Promise<paths["/posts/{id}/comments"]["get"]["responses"]["200"]["content"]["application/json"]>
  });

  const createMutation = useMutation({
    mutationFn: (body: string) => fetchApi(`/posts/${postId}/comments`, {
      method: "POST",
      body: JSON.stringify({ body })
    }),
    onSuccess: () => {
      setNewComment("");
      queryClient.invalidateQueries({ queryKey: ["comments", postId] });
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (newComment.trim()) {
      createMutation.mutate(newComment);
    }
  };

  const comments = response?.items || [];

  return (
    <Box sx={{ px: 2, pb: 2, pt: 1 }}>
      <Divider sx={{ mb: 2 }} />
      {isLoading ? (
        <Box sx={{ display: "flex", justifyContent: "center", p: 2 }}>
          <CircularProgress size={24} />
        </Box>
      ) : comments.length > 0 ? (
        comments.map((comment) => (
          <CommentItem key={comment.id} comment={comment} postId={postId} />
        ))
      ) : (
        <Typography variant="body2" color="textSecondary" sx={{ mb: 2, textAlign: "center" }}>
          No comments yet. Be the first to comment!
        </Typography>
      )}

      <Box component="form" onSubmit={handleSubmit} sx={{ display: "flex", gap: 1, mt: 1 }}>
        <TextField
          fullWidth
          size="small"
          placeholder="Write a comment..."
          value={newComment}
          onChange={(e) => setNewComment(e.target.value)}
          disabled={createMutation.isPending}
          variant="outlined"
          sx={{ "& .MuiOutlinedInput-root": { borderRadius: 5 } }}
        />
        <Button 
          type="submit" 
          variant="contained" 
          disabled={!newComment.trim() || createMutation.isPending}
          sx={{ borderRadius: 5, px: 3, textTransform: "none" }}
        >
          Post
        </Button>
      </Box>
    </Box>
  );
};
