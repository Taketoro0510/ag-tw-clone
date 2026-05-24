import { getAuth } from "firebase/auth";

const API_BASE = "http://localhost:8080/api/v1";

export async function fetchApi<T = unknown>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const getHeaders = () => {
    const headers = new Headers(options.headers);
    if (options.body && !headers.has("Content-Type") && typeof options.body === "string") {
      headers.set("Content-Type", "application/json");
    }
    const token = localStorage.getItem("jwt");
    if (token) {
      headers.set("Authorization", `Bearer ${token}`);
    }
    return headers;
  };

  let response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: getHeaders(),
  });

  if (response.status === 401) {
    const auth = getAuth();
    if (auth.currentUser) {
      try {
        const idToken = await auth.currentUser.getIdToken(true);
        const res = await fetch(`${API_BASE}/auth/sessions`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ idToken }),
        });
        if (res.ok) {
          const data = await res.json();
          localStorage.setItem("jwt", data.token);
          
          // Retry original request
          response = await fetch(`${API_BASE}${endpoint}`, {
            ...options,
            headers: getHeaders(),
          });
        }
      } catch (err) {
        console.error("Failed to refresh token", err);
      }
    }
  }

  if (!response.ok) {
    const errorData = await response.json().catch(() => null);
    throw new Error(errorData?.error?.message || `API Error: ${response.status}`);
  }

  if (response.status === 204) {
    return null as any;
  }
  return response.json();
}
