import React, { useState } from "react";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Box } from "@mui/material";
import { fetchApi } from "../../api/client";
import { storage } from "../../lib/firebase";
import { ref, uploadBytes, getDownloadURL } from "firebase/storage";
import { useAuth } from "../auth/AuthContext";
import { useQueryClient } from "@tanstack/react-query";
import { v7 as uuidv7 } from "uuid";

interface CreatePostModalProps {
  open: boolean;
  onClose: () => void;
}

export const CreatePostModal: React.FC<CreatePostModalProps> = ({ open, onClose }) => {
  const [body, setBody] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { user } = useAuth();
  const queryClient = useQueryClient();

  const handleClose = () => {
    setBody("");
    setFile(null);
    onClose();
  };

  const handleSubmit = async () => {
    if (!user) return;
    if (!body && !file) return;

    setIsSubmitting(true);
    try {
      let mediaType: "image" | "video" | undefined = undefined;
      let mediaPath: string | undefined = undefined;
      let mediaUrl: string | undefined = undefined;

      if (file) {
        const isVideo = file.type.startsWith("video/");
        mediaType = isVideo ? "video" : "image";
        
        const postId = uuidv7(); 
        mediaPath = `users/${user.uid}/posts/${postId}/${file.name}`;
        const storageRef = ref(storage, mediaPath);
        
        await uploadBytes(storageRef, file);
        mediaUrl = await getDownloadURL(storageRef);
      }

      await fetchApi("/posts", {
        method: "POST",
        body: JSON.stringify({
          body,
          mediaType,
          mediaPath,
          mediaUrl,
        }),
      });

      queryClient.invalidateQueries({ queryKey: ["posts"] });
      handleClose();
    } catch (err) {
      console.error(err);
      alert("Failed to create post.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} fullWidth maxWidth="sm">
      <DialogTitle>Create Post</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          label="What's happening?"
          fullWidth
          multiline
          rows={4}
          value={body}
          onChange={(e) => setBody(e.target.value.slice(0, 140))}
          helperText={`${body.length}/140`}
        />
        <Box sx={{ mt: 2 }}>
          <input
            accept="image/*,video/*"
            style={{ display: "none" }}
            id="raised-button-file"
            type="file"
            onChange={(e) => setFile(e.target.files?.[0] || null)}
          />
          <label htmlFor="raised-button-file">
            <Button variant="outlined" component="span">
              Upload Media
            </Button>
          </label>
          {file && <Box sx={{ mt: 1 }}>{file.name}</Box>}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} disabled={isSubmitting}>Cancel</Button>
        <Button onClick={handleSubmit} variant="contained" disabled={isSubmitting || (!body && !file)}>
          Post
        </Button>
      </DialogActions>
    </Dialog>
  );
};
