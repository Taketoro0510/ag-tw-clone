import React, { createContext, useContext, useEffect, useState } from "react";
import { onAuthStateChanged, signOut as firebaseSignOut, GoogleAuthProvider, signInWithRedirect, getRedirectResult } from "firebase/auth";
import type { User } from "firebase/auth";
import { auth } from "../../lib/firebase";

interface AuthContextType {
  user: User | null;
  loading: boolean;
  loginWithGoogle: () => Promise<void>;
  logout: () => Promise<void>;
  token: string | null;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem("jwt"));
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (currentUser) => {
      setUser(currentUser);
      if (currentUser) {
        try {
          const idToken = await currentUser.getIdToken();
          const res = await fetch("http://localhost:8080/api/v1/auth/sessions", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ idToken }),
          });
          if (res.ok) {
            const data = await res.json();
            localStorage.setItem("jwt", data.token);
            setToken(data.token);
          }
        } catch (err) {
          console.error("Failed to sync session JWT on auth state change", err);
        }
      } else {
        localStorage.removeItem("jwt");
        setToken(null);
      }
      setLoading(false);
    });

    const handleRedirect = async () => {
      try {
        await getRedirectResult(auth);
      } catch (err) {
        console.error("Redirect login failed", err);
      }
    };
    handleRedirect();

    return unsubscribe;
  }, []);

  const loginWithGoogle = async () => {
    const provider = new GoogleAuthProvider();
    await signInWithRedirect(auth, provider);
  };

  const logout = async () => {
    await firebaseSignOut(auth);
    localStorage.removeItem("jwt");
    setToken(null);
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, loginWithGoogle, logout, token }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
