import { ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { lightTheme } from "./theme/lightTheme";
import { AuthProvider } from "./features/auth/AuthContext";
import { RequireAuth } from "./features/auth/RequireAuth";
import { Layout } from "./components/Layout";
import { Login } from "./pages/Login";
import { Timeline } from "./pages/Timeline";
import { Bookmarks } from "./pages/Bookmarks";
import { PostDetail } from "./pages/PostDetail";
import { Profile } from "./pages/Profile";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={lightTheme}>
        <CssBaseline />
        <AuthProvider>
          <BrowserRouter>
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route element={<RequireAuth><Layout /></RequireAuth>}>
                <Route path="/" element={<Timeline />} />
                <Route path="/bookmarks" element={<Bookmarks />} />
                <Route path="/posts/:id" element={<PostDetail />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="/users/:id" element={<Profile />} />
              </Route>
            </Routes>
          </BrowserRouter>
        </AuthProvider>
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App;
