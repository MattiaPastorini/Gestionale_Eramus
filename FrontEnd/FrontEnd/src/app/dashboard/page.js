"use client";
import { useEffect, useState } from "react";
import api from "@/services/api";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Pie } from "react-chartjs-2";

ChartJS.register(ArcElement, Tooltip, Legend);

export default function Dashboard() {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await api.get("/dashboard/statistiche");
        setStats(res.data);
      } catch (err) {
        console.error("Errore caricamento dashboard", err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  if (loading)
    return <div className="p-5 text-center">Caricamento in corso...</div>;

  const movimenti = stats?.ultimi_movimenti || [];
  const categorie = stats?.grafico_categorie || [];

  const chartData = {
    labels:
      categorie.length > 0 ? categorie.map((c) => c.nome) : ["Nessun dato"],
    datasets: [
      {
        data: categorie.length > 0 ? categorie.map((c) => c.quantita) : [1],
        backgroundColor: [
          "#0066CC",
          "#2ECC71",
          "#F1C40F",
          "#E74C3C",
          "#9B59B6",
        ],
      },
    ],
  };

  return (
    <div
      className="container-fluid p-4"
      style={{ backgroundColor: "#f8f9fa", minHeight: "100vh" }}
    >
      <h1 className="mb-4 fw-bold" style={{ color: "#0066CC" }}>
        Dashboard Amministrativa
      </h1>

      <div className="row g-4 mb-5">
        <div className="col-md-4">
          <div className="card shadow-sm border-0 border-start border-primary border-4 p-3">
            <small className="text-muted text-uppercase fw-bold">
              Totale Utenti
            </small>
            <h2 className="fw-bold">{stats?.total_utenti || 0}</h2>
          </div>
        </div>
        <div className="col-md-4">
          <div className="card shadow-sm border-0 border-start border-success border-4 p-3">
            <small className="text-muted text-uppercase fw-bold">
              Valore Inventario
            </small>
            <h2 className="fw-bold">
              €{" "}
              {stats?.valore_inventario
                ? stats.valore_inventario.toFixed(2)
                : "0.00"}
            </h2>
          </div>
        </div>
        <div className="col-md-4">
          <div className="card shadow-sm border-0 border-start border-warning border-4 p-3">
            <small className="text-muted text-uppercase fw-bold">
              Totale Prodotti
            </small>
            <h2 className="fw-bold">{stats?.total_prodotti || 0}</h2>
          </div>
        </div>
      </div>

      <div className="row">
        <div className="col-lg-8 mb-4">
          <div className="card shadow-sm border-0 p-4 h-100">
            <h5 className="fw-bold mb-4">Ultimi 5 movimenti</h5>
            <div className="table-responsive">
              <table className="table table-hover">
                <thead className="table-light">
                  <tr>
                    <th>Data</th>
                    <th>Prodotto</th>
                    <th>Tipo</th>
                    <th>Quantità</th>
                  </tr>
                </thead>
                <tbody>
                  {movimenti.length > 0 ? (
                    movimenti.map((m, index) => (
                      <tr key={m.Id || index}>
                        <td>
                          {m.data_movimento
                            ? new Date(m.data_movimento).toLocaleDateString()
                            : "-"}
                        </td>
                        <td>
                          {m.Prodotto?.nome_oggetto || "Prodotto eliminato"}
                        </td>
                        <td>
                          <span
                            className={`badge ${m.tipo_movimento?.includes("Carico") ? "bg-success" : "bg-danger"}`}
                          >
                            {m.tipo_movimento}
                          </span>
                        </td>
                        <td>{m.quantita}</td>
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <td colSpan="4" className="text-center text-muted">
                        Nessun movimento recente
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <div className="col-lg-4 mb-4">
          <div className="card shadow-sm border-0 p-4 h-100">
            <h5 className="fw-bold mb-4 text-center">Prodotti per categoria</h5>
            <div style={{ maxWidth: "300px", margin: "0 auto" }}>
              <Pie data={chartData} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
