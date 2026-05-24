import React, { useState } from "react";
import { Box, Typography, CircularProgress, Button, Fab } from "@mui/material";
import { useInfiniteQuery } from "@tanstack/react-query";
import { fetchApi } from "../api/client";
import { PostCard } from "../features/posts/PostCard";
import { CreatePostModal } from "../features/posts/CreatePostModal";
import type { paths } from "../api/types";

type GetPostsResponse = paths["/posts"]["get"]["responses"]["200"]["content"]["application/json"];

export const Timeline: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);

  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, status } = useInfiniteQuery<GetPostsResponse>({
    queryKey: ["posts"],
    queryFn: ({ pageParam }) => 
      fetchApi<GetPostsResponse>(`/posts?limit=20${pageParam ? `&cursor=${pageParam}` : ""}`),
    getNextPageParam: (lastPage) => lastPage.nextCursor || undefined,
    initialPageParam: "",
  });

  if (status === "pending") return <CircularProgress />;
  if (status === "error") return <Typography color="error">Failed to load timeline.</Typography>;

  return (
    <Box>
      <Typography variant="h5" gutterBottom sx={{ fontWeight: "bold" }}>Timeline</Typography>
      {data.pages.map((page, i) => (
        <React.Fragment key={i}>
          {(page.items || []).map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </React.Fragment>
      ))}
      {hasNextPage && (
        <Button onClick={() => fetchNextPage()} disabled={isFetchingNextPage} fullWidth sx={{ mt: 2, mb: 4 }}>
          {isFetchingNextPage ? "Loading..." : "Load More"}
        </Button>
      )}

      <Fab
        color="primary"
        aria-label="add"
        onClick={() => setIsModalOpen(true)}
        sx={{ position: "fixed", bottom: 32, right: 32 }}
      >
        +
      </Fab>

      <CreatePostModal open={isModalOpen} onClose={() => setIsModalOpen(false)} />
    </Box>
  );
};
