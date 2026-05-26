import React from "react";
import { Card, CardHeader, CardContent, Typography, Avatar, Box, CardActions, IconButton, Badge } from "@mui/material";
import FavoriteIcon from "@mui/icons-material/Favorite";
import FavoriteBorderIcon from "@mui/icons-material/FavoriteBorder";
import BookmarkIcon from "@mui/icons-material/Bookmark";
import BookmarkBorderIcon from "@mui/icons-material/BookmarkBorder";
import ShareIcon from "@mui/icons-material/Share";
import DeleteIcon from "@mui/icons-material/Delete";
import ChatBubbleOutlineIcon from "@mui/icons-material/ChatBubbleOutlined";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchApi } from "../../api/client";
import type { paths } from "../../api/types";
import { useAuth } from "../auth/AuthContext";
import { CommentSection } from "./CommentSection";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

type Post = NonNullable<paths["/posts"]["get"]["responses"]["200"]["content"]["application/json"]["items"]>[0];

export const PostCard: React.FC<{ post: Post; defaultShowComments?: boolean }> = ({ post, defaultShowComments = false }) => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const [showComments, setShowComments] = useState(defaultShowComments);

  const likeMutation = useMutation({
    mutationFn: () => fetchApi(`/posts/${post.id}/likes`, { method: post.likedByMe ? "DELETE" : "POST" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["bookmarks"] });
    }
  });

  const bookmarkMutation = useMutation({
    mutationFn: () => fetchApi(`/posts/${post.id}/bookmarks`, { method: post.bookmarkedByMe ? "DELETE" : "POST" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["bookmarks"] });
    }
  });

  const deleteMutation = useMutation({
    mutationFn: () => fetchApi(`/posts/${post.id}`, { method: "DELETE" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["bookmarks"] });
    }
  });

  const handleShare = () => {
    navigator.clipboard.writeText(`${window.location.origin}/posts/${post.id}`);
    alert("URL copied to clipboard!");
  };

  const handleDelete = () => {
    if (window.confirm("Are you sure you want to delete this post?")) {
      deleteMutation.mutate();
    }
  };

  const isMyPost = user && post.author?.firebaseUid === user.uid;

  const authorName = post.author?.displayName || (post.authorId ? `User ${post.authorId.substring(0, 5)}` : "Unknown");
  const authorAvatar = post.author?.avatarUrl;

  return (
    <Card sx={{ 
      mb: 3, 
      borderRadius: 3, 
      boxShadow: "0 8px 24px rgba(0,0,0,0.06)",
      transition: "transform 0.2s ease-in-out",
      "&:hover": {
        transform: "translateY(-2px)",
        boxShadow: "0 12px 28px rgba(0,0,0,0.08)",
      }
    }}>
      <CardHeader
        avatar={
          <Avatar 
            sx={{ bgcolor: "primary.main", cursor: "pointer" }} 
            src={authorAvatar || undefined}
            onClick={() => navigate(`/users/${post.authorId}`)}
          >
            {!authorAvatar ? authorName.substring(0, 2).toUpperCase() : ""}
          </Avatar>
        }
        title={
          <Typography 
            variant="subtitle1" 
            sx={{ fontWeight: "bold", cursor: "pointer", "&:hover": { textDecoration: "underline" } }}
            onClick={() => navigate(`/users/${post.authorId}`)}
          >
            {authorName}
          </Typography>
        }
        subheader={<Typography variant="caption" color="textSecondary">{post.createdAt ? new Date(post.createdAt).toLocaleString() : ""}</Typography>}
        action={
          isMyPost && (
            <IconButton onClick={handleDelete} color="error" size="small">
              <DeleteIcon />
            </IconButton>
          )
        }
      />
      <CardContent sx={{ pt: 0, pb: 1 }}>
        <Typography variant="body1" sx={{ whiteSpace: "pre-wrap", lineHeight: 1.6 }}>{post.body}</Typography>
        {post.mediaUrl && (
          <Box sx={{ mt: 2, borderRadius: 2, overflow: "hidden" }}>
            {post.mediaType === "video" ? (
              <video src={post.mediaUrl} controls width="100%" />
            ) : (
              <img src={post.mediaUrl} alt="post media" style={{ width: "100%", maxHeight: 450, objectFit: "cover", display: "block" }} />
            )}
          </Box>
        )}
      </CardContent>
      <CardActions disableSpacing sx={{ borderTop: "1px solid", borderColor: "divider", px: 2 }}>
        <IconButton aria-label="like" onClick={() => likeMutation.mutate()} color={post.likedByMe ? "error" : "default"}>
          <Badge badgeContent={post.likeCount} color="error">
            {post.likedByMe ? <FavoriteIcon /> : <FavoriteBorderIcon />}
          </Badge>
        </IconButton>
        <IconButton aria-label="comment" onClick={() => setShowComments(!showComments)}>
          <ChatBubbleOutlineIcon />
        </IconButton>
        <IconButton aria-label="bookmark" onClick={() => bookmarkMutation.mutate()} color={post.bookmarkedByMe ? "primary" : "default"}>
          <Badge badgeContent={post.bookmarkCount} color="primary">
            {post.bookmarkedByMe ? <BookmarkIcon /> : <BookmarkBorderIcon />}
          </Badge>
        </IconButton>
        <Box sx={{ flexGrow: 1 }} />
        <IconButton aria-label="share" onClick={handleShare}>
          <ShareIcon />
        </IconButton>
      </CardActions>
      {showComments && post.id && <CommentSection postId={post.id} />}
    </Card>
  );
};
