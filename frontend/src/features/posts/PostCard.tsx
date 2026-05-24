import React from "react";
import { Card, CardHeader, CardContent, Typography, Avatar, Box } from "@mui/material";
import type { paths } from "../../api/types";

type Post = NonNullable<paths["/posts"]["get"]["responses"]["200"]["content"]["application/json"]["items"]>[0];

export const PostCard: React.FC<{ post: Post }> = ({ post }) => {
  return (
    <Card sx={{ mb: 2 }}>
      <CardHeader
        avatar={<Avatar>{post.authorId ? post.authorId.substring(0, 2) : "?"}</Avatar>}
        title={post.authorId ? `User ${post.authorId.substring(0, 5)}` : "Unknown"}
        subheader={post.createdAt ? new Date(post.createdAt).toLocaleString() : ""}
      />
      <CardContent>
        <Typography variant="body1">{post.body}</Typography>
        {post.mediaUrl && (
          <Box sx={{ mt: 2 }}>
            {post.mediaType === "video" ? (
              <video src={post.mediaUrl} controls width="100%" />
            ) : (
              <img src={post.mediaUrl} alt="post media" style={{ width: "100%", maxHeight: 400, objectFit: "cover" }} />
            )}
          </Box>
        )}
      </CardContent>
    </Card>
  );
};
