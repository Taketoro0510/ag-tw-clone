import React from "react";
import { AppBar, Toolbar, Typography, Button, Container, Box, Drawer, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Paper } from "@mui/material";
import { useAuth } from "../features/auth/AuthContext";
import { Outlet, useNavigate, useLocation } from "react-router-dom";
import HomeIcon from "@mui/icons-material/Home";
import BookmarkIcon from "@mui/icons-material/Bookmark";
import PersonIcon from "@mui/icons-material/Person";

const drawerWidth = { xs: 72, md: 240 };

export const Layout: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <Box sx={{ display: 'flex' }}>
      <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1, color: "inherit", fontWeight: "bold" }}>
            CloudCode SNS
          </Typography>
          {user && (
            <Button color="inherit" onClick={logout}>
              Logout
            </Button>
          )}
        </Toolbar>
      </AppBar>
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: { 
            width: drawerWidth, 
            boxSizing: 'border-box',
            transition: 'width 0.2s ease-in-out',
            overflowX: 'hidden'
          },
        }}
      >
        <Toolbar />
        <Box sx={{ overflow: 'auto', py: 1 }}>
          <List>
            <ListItem disablePadding>
              <ListItemButton 
                selected={location.pathname === '/'} 
                onClick={() => navigate('/')}
                sx={{ 
                  py: 1.5,
                  px: { xs: 1, md: 3 },
                  justifyContent: { xs: 'center', md: 'flex-start' },
                  borderRadius: 2,
                  mx: 1,
                  mb: 0.5,
                }}
              >
                <ListItemIcon sx={{ 
                  minWidth: { xs: 0, md: 40 },
                  mr: { xs: 0, md: 0 },
                  justifyContent: 'center',
                  color: location.pathname === '/' ? 'primary.main' : 'text.secondary',
                }}>
                  <HomeIcon />
                </ListItemIcon>
                <ListItemText 
                  primary="Timeline" 
                  sx={{ 
                    display: { xs: 'none', md: 'block' },
                    '& .MuiTypography-root': { fontWeight: location.pathname === '/' ? 'bold' : 'normal' }
                  }} 
                />
              </ListItemButton>
            </ListItem>
            <ListItem disablePadding>
              <ListItemButton 
                selected={location.pathname === '/bookmarks'} 
                onClick={() => navigate('/bookmarks')}
                sx={{ 
                  py: 1.5,
                  px: { xs: 1, md: 3 },
                  justifyContent: { xs: 'center', md: 'flex-start' },
                  borderRadius: 2,
                  mx: 1,
                  mb: 0.5,
                }}
              >
                <ListItemIcon sx={{ 
                  minWidth: { xs: 0, md: 40 },
                  mr: { xs: 0, md: 0 },
                  justifyContent: 'center',
                  color: location.pathname === '/bookmarks' ? 'primary.main' : 'text.secondary',
                }}>
                  <BookmarkIcon />
                </ListItemIcon>
                <ListItemText 
                  primary="Bookmarks" 
                  sx={{ 
                    display: { xs: 'none', md: 'block' },
                    '& .MuiTypography-root': { fontWeight: location.pathname === '/bookmarks' ? 'bold' : 'normal' }
                  }} 
                />
              </ListItemButton>
            </ListItem>
            <ListItem disablePadding>
              <ListItemButton 
                selected={location.pathname === '/profile'} 
                onClick={() => navigate('/profile')}
                sx={{ 
                  py: 1.5,
                  px: { xs: 1, md: 3 },
                  justifyContent: { xs: 'center', md: 'flex-start' },
                  borderRadius: 2,
                  mx: 1,
                  mb: 0.5,
                }}
              >
                <ListItemIcon sx={{ 
                  minWidth: { xs: 0, md: 40 },
                  mr: { xs: 0, md: 0 },
                  justifyContent: 'center',
                  color: location.pathname === '/profile' ? 'primary.main' : 'text.secondary',
                }}>
                  <PersonIcon />
                </ListItemIcon>
                <ListItemText 
                  primary="Profile" 
                  sx={{ 
                    display: { xs: 'none', md: 'block' },
                    '& .MuiTypography-root': { fontWeight: location.pathname === '/profile' ? 'bold' : 'normal' }
                  }} 
                />
              </ListItemButton>
            </ListItem>
          </List>
        </Box>
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, p: 3, bgcolor: "background.default", minHeight: "100vh", width: { xs: 'calc(100% - 72px)', md: 'calc(100% - 240px)' } }}>
        <Toolbar />
        <Container maxWidth="lg">
          <Box sx={{ display: "flex", gap: 4 }}>
            <Box sx={{ flexGrow: 1, width: { xs: "100%", sm: "calc(100% - 320px)" } }}>
              <Outlet />
            </Box>
            <Box sx={{ width: 288, flexShrink: 0, display: { xs: "none", sm: "block" } }}>
              <Box sx={{ position: "sticky", top: 88 }}>
                <Paper elevation={1} sx={{ p: 2, borderRadius: 3, mb: 3 }}>
                  <Typography variant="h6" sx={{ fontWeight: "bold" }} gutterBottom>
                    Trending Topics
                  </Typography>
                  <List dense>
                    <ListItem disablePadding><ListItemButton sx={{ borderRadius: 1 }}><ListItemText primary="#CloudCode" secondary="10.5k posts" /></ListItemButton></ListItem>
                    <ListItem disablePadding><ListItemButton sx={{ borderRadius: 1 }}><ListItemText primary="#GoLang" secondary="8.2k posts" /></ListItemButton></ListItem>
                    <ListItem disablePadding><ListItemButton sx={{ borderRadius: 1 }}><ListItemText primary="#ReactMUI" secondary="3.1k posts" /></ListItemButton></ListItem>
                  </List>
                </Paper>
                
                <Paper elevation={1} sx={{ p: 2, borderRadius: 3 }}>
                  <Typography variant="h6" sx={{ fontWeight: "bold" }} gutterBottom>
                    Who to follow
                  </Typography>
                  <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                    Suggestions will appear here based on your activity.
                  </Typography>
                  <Button variant="outlined" size="small" fullWidth sx={{ borderRadius: 2, mt: 1 }}>
                    Find People
                  </Button>
                </Paper>
              </Box>
            </Box>
          </Box>
        </Container>
      </Box>
    </Box>
  );
};
