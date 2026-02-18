"use client";
import "bootstrap/dist/css/bootstrap.min.css";
import "./globals.css";
import Link from "next/link";

export default function RootLayout({ children }) {
  const handleLogout = () => {
    localStorage.clear();
    window.location.href = "/";
  };
  return (
    <html lang="it">
      <body suppressHydrationWarning={true}>
        <nav
          className="navbar navbar-expand-lg navbar-dark shadow-sm mb-4"
          style={{ backgroundColor: "#0066CC" }}
        >
          <div className="container">
            <Link className="navbar-brand fw-bold" href="/dashboard">
              Gestionale Eramus
            </Link>
            <div className="collapse navbar-collapse" id="navbarNav">
              <ul className="navbar-nav me-auto">
                <li className="nav-item">
                  <Link className="nav-link" href="/dashboard">
                    Dashboard
                  </Link>
                </li>
                <li className="nav-item">
                  <Link className="nav-link" href="/utenti">
                    Gestione Utenti
                  </Link>
                </li>
                <li className="nav-item">
                  <Link className="nav-link" href="/inventario">
                    Inventario
                  </Link>
                </li>
              </ul>
              <div className="d-flex">
                <button
                  className="btn btn-outline-light btn-sm"
                  onClick={() => {
                    localStorage.clear();
                    window.location.href = "/";
                  }}
                >
                  Logout
                </button>
              </div>
            </div>
          </div>
        </nav>
        <main>{children}</main>
      </body>
    </html>
  );
}
