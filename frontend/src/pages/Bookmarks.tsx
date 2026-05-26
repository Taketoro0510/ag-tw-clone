import React from "react";
import { Box, Typography, CircularProgress, Button } from "@mui/material";
import { useInfiniteQuery } from "@tanstack/react-query";
import { fetchApi } from "../api/client";
import { PostCard } from "../features/posts/PostCard";
import type { paths } from "../api/types";

type GetPostsResponse = paths["/bookmarks"]["get"]["responses"]["200"]["content"]["application/json"];

export const Bookmarks: React.FC = () => {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, status } = useInfiniteQuery<GetPostsResponse>({
    queryKey: ["bookmarks"],
    queryFn: ({ pageParam }) => 
      fetchApi<GetPostsResponse>(`/bookmarks?limit=20${pageParam ? `&cursor=${pageParam}` : ""}`),
    getNextPageParam: (lastPage) => lastPage.nextCursor || undefined,
    initialPageParam: "",
  });

  if (status === "pending") return <CircularProgress />;
  if (status === "error") return <Typography color="error">Failed to load bookmarks.</Typography>;

  return (
    <Box>
      <Typography variant="h5" gutterBottom sx={{ fontWeight: "bold" }}>Bookmarks</Typography>
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
      {data.pages[0].items?.length === 0 && (
        <Typography color="textSecondary" sx={{ mt: 4, textAlign: 'center' }}>
          No bookmarks yet.
        </Typography>
      )}
    </Box>
  );
};
