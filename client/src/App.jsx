// src/App.jsx
import "./index.css";
import Header from "./components/Header";
import ModernTabs from "./components/ModernTabs";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./auth/AuthProvider";
import LoginPage from "./pages/LoginPage";
import RequireAuth from "./components/RequireAuth";

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <main className="min-h-screen bg-gray-50">
          <Routes>
            {/* Login page (no header, no padding) */}
            <Route path="/login" element={<LoginPage />} />

            {/* Protected main app (with header + content layout) */}
            <Route
              path="/"
              element={
                <RequireAuth>
                  <>
                    <Header />
                    <div className="mx-auto max-w-7xl px-4 py-2 sm:px-6 lg:px-8">
                      <ModernTabs />
                    </div>
                  </>
                </RequireAuth>
              }
            />

            {/* Redirect unknown routes */}
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </main>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;
