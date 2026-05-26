import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Box, Typography, Avatar, Button, Paper, Tabs, Tab, CircularProgress, List, ListItem, ListItemAvatar, ListItemText } from "@mui/material";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchApi } from "../api/client";
import { PostCard } from "../features/posts/PostCard";
import type { paths } from "../api/types";

type UserDTO = paths["/users/{id}"]["get"]["responses"]["200"]["content"]["application/json"];
type PostDTO = NonNullable<paths["/users/{id}/posts"]["get"]["responses"]["200"]["content"]["application/json"]["items"]>[number];

export const Profile: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState(0);

  // Fetch current user /me to get own ID if "id" param is not set
  const { data: me, isLoading: loadingMe } = useQuery<UserDTO>({
    queryKey: ["me"],
    queryFn: () => fetchApi<UserDTO>("/me"),
  });

  const targetUserId = id || me?.id;

  // Fetch target user profile
  const { data: profile, isLoading: loadingProfile, error: profileError } = useQuery<UserDTO>({
    queryKey: ["user", targetUserId],
    queryFn: () => fetchApi<UserDTO>(`/users/${targetUserId}`),
    enabled: !!targetUserId,
  });

  // Fetch target user's posts
  const { data: postsResponse, isLoading: loadingPosts } = useQuery<{ items: PostDTO[] }>({
    queryKey: ["user-posts", targetUserId],
    queryFn: () => fetchApi<{ items: PostDTO[] }>(`/users/${targetUserId}/posts`),
    enabled: !!targetUserId,
  });

  // Fetch followers
  const { data: followers, isLoading: loadingFollowers } = useQuery<UserDTO[]>({
    queryKey: ["followers", targetUserId],
    queryFn: () => fetchApi<UserDTO[]>(`/users/${targetUserId}/followers`),
    enabled: !!targetUserId && activeTab === 1,
  });

  // Fetch following
  const { data: following, isLoading: loadingFollowing } = useQuery<UserDTO[]>({
    queryKey: ["following", targetUserId],
    queryFn: () => fetchApi<UserDTO[]>(`/users/${targetUserId}/following`),
    enabled: !!targetUserId && activeTab === 2,
  });

  // Follow/Unfollow Mutation
  const followMutation = useMutation({
    mutationFn: () => fetchApi(`/users/${targetUserId}/follow`, { method: profile?.followedByMe ? "DELETE" : "POST" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user", targetUserId] });
      queryClient.invalidateQueries({ queryKey: ["followers", targetUserId] });
      queryClient.invalidateQueries({ queryKey: ["following", me?.id] });
    },
  });

  const isLoading = loadingMe || (targetUserId && loadingProfile);

  if (isLoading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (profileError || !profile) {
    return (
      <Box sx={{ py: 4, textAlign: "center" }}>
        <Typography variant="h6" color="error">
          User profile not found or failed to load.
        </Typography>
        <Button variant="contained" onClick={() => navigate("/")} sx={{ mt: 2, borderRadius: 5 }}>
          Back to Timeline
        </Button>
      </Box>
    );
  }

  const isOwnProfile = me?.id === profile.id;
  const authorName = profile.displayName || (profile.id ? `User ${profile.id.substring(0, 5)}` : "Unknown");
  const posts = postsResponse?.items || [];

  return (
    <Box>
      {/* Profile Header Card */}
      <Paper elevation={0} sx={{ border: "1px solid", borderColor: "divider", borderRadius: 4, overflow: "hidden", mb: 4 }}>
        {/* Cover Banner with smooth modern gradient */}
        <Box sx={{ height: 160, background: "linear-gradient(135deg, #1e3c72 0%, #2a5298 100%)" }} />
        <Box sx={{ px: 3, pb: 3, pt: 0, position: "relative" }}>
          {/* Overlapping Avatar */}
          <Avatar
            src={profile.avatarUrl || undefined}
            sx={{
              width: 100,
              height: 100,
              border: "4px solid white",
              mt: -6,
              mb: 2,
              bgcolor: "primary.main",
              fontSize: "2rem",
              boxShadow: "0 8px 16px rgba(0,0,0,0.1)",
            }}
          >
            {!profile.avatarUrl ? authorName.substring(0, 2).toUpperCase() : ""}
          </Avatar>

          <Box sx={{ display: "flex", flexWrap: "wrap", justifyContent: "space-between", alignItems: "flex-start" }}>
            <Box sx={{ width: { xs: "100%", sm: "66%" } }}>
              <Typography variant="h4" sx={{ fontWeight: "bold" }}>
                {authorName}
              </Typography>
              <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
                {profile.email}
              </Typography>

              {/* Stats Block */}
              <Box sx={{ display: "flex", gap: 3, mb: 1 }}>
                <Box onClick={() => setActiveTab(1)} sx={{ cursor: "pointer", "&:hover": { opacity: 0.8 } }}>
                  <Typography variant="subtitle1" sx={{ fontWeight: "bold" }} component="span">
                    {profile.followersCount}
                  </Typography>
                  <Typography variant="body2" color="textSecondary" component="span" sx={{ ml: 0.5 }}>
                    Followers
                  </Typography>
                </Box>
                <Box onClick={() => setActiveTab(2)} sx={{ cursor: "pointer", "&:hover": { opacity: 0.8 } }}>
                  <Typography variant="subtitle1" sx={{ fontWeight: "bold" }} component="span">
                    {profile.followingCount}
                  </Typography>
                  <Typography variant="body2" color="textSecondary" component="span" sx={{ ml: 0.5 }}>
                    Following
                  </Typography>
                </Box>
              </Box>
            </Box>

            <Box sx={{ width: { xs: "100%", sm: "33%" }, display: "flex", justifyContent: { xs: "flex-start", sm: "flex-end" }, mt: { xs: 2, sm: 0 } }}>
              {!isOwnProfile && (
                <Button
                  variant={profile.followedByMe ? "outlined" : "contained"}
                  color="primary"
                  onClick={() => followMutation.mutate()}
                  disabled={followMutation.isPending}
                  sx={{ borderRadius: 5, px: 4, textTransform: "none", fontWeight: "bold" }}
                >
                  {profile.followedByMe ? "Unfollow" : "Follow"}
                </Button>
              )}
            </Box>
          </Box>
        </Box>
      </Paper>

      {/* Tabs list */}
      <Box sx={{ borderBottom: 1, borderColor: "divider", mb: 3 }}>
        <Tabs value={activeTab} onChange={(_, val) => setActiveTab(val)} aria-label="profile tabs">
          <Tab label="Posts" sx={{ textTransform: "none", fontWeight: "bold" }} />
          <Tab label="Followers" sx={{ textTransform: "none", fontWeight: "bold" }} />
          <Tab label="Following" sx={{ textTransform: "none", fontWeight: "bold" }} />
        </Tabs>
      </Box>

      {/* Tab Panels */}
      {activeTab === 0 && (
        <Box>
          {loadingPosts ? (
            <CircularProgress />
          ) : posts.length > 0 ? (
            posts.map((post) => <PostCard key={post.id} post={post} />)
          ) : (
            <Typography color="textSecondary" align="center" sx={{ py: 4 }}>
              No posts yet.
            </Typography>
          )}
        </Box>
      )}

      {activeTab === 1 && (
        <Box>
          {loadingFollowers ? (
            <CircularProgress />
          ) : followers && followers.length > 0 ? (
            <List sx={{ bgcolor: "background.paper", borderRadius: 3, border: "1px solid", borderColor: "divider" }}>
              {followers.map((usr) => {
                const uName = usr.displayName || (usr.id ? `User ${usr.id.substring(0, 5)}` : "Unknown");
                return (
                  <ListItem
                    key={usr.id}
                    sx={{ borderBottom: "1px solid", borderColor: "divider", "&:last-child": { borderBottom: 0 }, cursor: "pointer" }}
                    onClick={() => {
                      if (usr.id) {
                        navigate(`/users/${usr.id}`);
                        setActiveTab(0);
                      }
                    }}
                  >
                    <ListItemAvatar>
                      <Avatar src={usr.avatarUrl || undefined}>{!usr.avatarUrl ? uName.substring(0, 2).toUpperCase() : ""}</Avatar>
                    </ListItemAvatar>
                    <ListItemText primary={uName} secondary={usr.email} />
                  </ListItem>
                );
              })}
            </List>
          ) : (
            <Typography color="textSecondary" align="center" sx={{ py: 4 }}>
              No followers yet.
            </Typography>
          )}
        </Box>
      )}

      {activeTab === 2 && (
        <Box>
          {loadingFollowing ? (
            <CircularProgress />
          ) : following && following.length > 0 ? (
            <List sx={{ bgcolor: "background.paper", borderRadius: 3, border: "1px solid", borderColor: "divider" }}>
              {following.map((usr) => {
                const uName = usr.displayName || (usr.id ? `User ${usr.id.substring(0, 5)}` : "Unknown");
                return (
                  <ListItem
                    key={usr.id}
                    sx={{ borderBottom: "1px solid", borderColor: "divider", "&:last-child": { borderBottom: 0 }, cursor: "pointer" }}
                    onClick={() => {
                      if (usr.id) {
                        navigate(`/users/${usr.id}`);
                        setActiveTab(0);
                      }
                    }}
                  >
                    <ListItemAvatar>
                      <Avatar src={usr.avatarUrl || undefined}>{!usr.avatarUrl ? uName.substring(0, 2).toUpperCase() : ""}</Avatar>
                    </ListItemAvatar>
                    <ListItemText primary={uName} secondary={usr.email} />
                  </ListItem>
                );
              })}
            </List>
          ) : (
            <Typography color="textSecondary" align="center" sx={{ py: 4 }}>
              Not following anyone yet.
            </Typography>
          )}
        </Box>
      )}
    </Box>
  );
};
