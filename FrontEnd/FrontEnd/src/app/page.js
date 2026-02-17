"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import api from "@/services/api";
import "bootstrap/dist/css/bootstrap.min.css";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const handleLogin = async (e) => {
    e.preventDefault();
    setError("");

    try {
      // Chiamata al backend Go
      const response = await api.post("/login", { username, password });

      // Salvataggio Token (Requisito punto 5)
      localStorage.setItem("access_token", response.data.access_token);
      localStorage.setItem("refresh_token", response.data.refresh_token);

      // Reindirizzamento alla Dashboard
      router.push("/dashboard");
    } catch (err) {
      setError(err.response?.data?.error || "Errore di connessione");
    }
  };

  return (
    <div className="container">
      <div className="row justify-content-center">
        <div className="col-12 col-lg-6">
          <h1 className="text-center mt-5">Accesso Gestionale</h1>

          <form
            onSubmit={handleLogin}
            className="m-5 border border-2 rounded-3 p-5 shadow-sm"
          >
            {error && <div className="alert alert-danger">{error}</div>}

            <div className="form-group mb-3">
              <label htmlFor="username">Username</label>
              <input
                type="text"
                className="form-control"
                id="username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
            </div>

            <div className="form-group mb-3">
              <label htmlFor="password">Password</label>
              <input
                type="password"
                className="form-control"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>

            <div className="d-flex justify-content-center mb-4">
              <a href="/forgot-password">Password dimenticata?</a>
            </div>

            <div className="d-flex justify-content-center">
              <button type="submit" className="btn btn-primary px-5">
                Accedi
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
