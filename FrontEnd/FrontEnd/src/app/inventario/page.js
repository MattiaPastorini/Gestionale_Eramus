"use client";
import { useEffect, useState } from "react";
import api from "@/services/api";

export default function Inventario() {
  const [prodotti, setProdotti] = useState([]);
  const [filtroTipo, setFiltroTipo] = useState("[];");
  const [search, setSearch] = useState("");

  const fetchProdotti = async () => {
    try {
      const res = await api.get(
        `/inventario/prodotti?nome=${search}&tipo=${filtroTipo}`,
      );
      setProdotti(res.data || []);
    } catch (err) {
      console.error("Errore nel caricameto dell'inventario", err);
    }
  };

  useEffect(() => {
    fetchProdotti();
  }, [search, filtroTipo]);

  const handleAggiornamentoStock = async (id, attuale) => {
    const nuova = prompt("Inserisci la nuova quantità totale:", attuale);
    if (nuova !== null) {
      try {
        await api.put(`/inventario/prodotti/${id}/stock`, {
          nuova_quantita: parseInt(nuova),
          note: "Aggiornamento manuale dalla dashboard",
        });
        fetchProdotti();
      } catch (err) {
        alert("Errore nell'aggiornamento stock");
      }
    }
  };

  return (
    <div className="container">
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h2 className="fw-bold text-primary">Gestione Inventario</h2>
        <button className="btn btn-success">Aggiungi prodotto</button>
      </div>

      <div className="row g-3 mb-4">
        <div className="col-md-6">
          <input
            type="text"
            className="form-control"
            placeholder="Cerca prodotto per nome..."
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        <div className="col-md-6">
          <select
            className="form-select"
            onChange={(e) => setFiltroTipo(e.target.value)}
          >
            <option value="">Tutte le categorie</option>
            <option value="Buste">Buste</option>
            <option value="Carta">Carta</option>
            <option value="Toner">Toner</option>
          </select>
        </div>
      </div>

      <div className="card shadow-sm border-0">
        <div className="table-responsive">
          <table className="table table-hover align-middle">
            <thead className="table-light">
              <tr>
                <th>Nome</th>
                <th>Categoria</th>
                <th>Prezzo Unitario</th>
                <th>Giacenza</th>
                <th>Stato</th>
                <th className="text-end">Azioni</th>
              </tr>
            </thead>
            <tbody>
              {prodotti.map((p) => (
                <tr key={p.Id}>
                  <td>
                    <strong>{p.nome_oggetto}</strong>
                    <br />
                    <small className="text-muted">{p.descrizione}</small>
                  </td>
                  <td>
                    <span className="badge bg-secondary">
                      {p.TipoProdotto?.corpo_messaggio}
                    </span>
                  </td>
                  <td>€ {p.prezzo_unitario.toFixed(2)}</td>
                  <td>
                    <span
                      className={`fw-bold ${p.quantita_disponibile < p.soglia_minima_di_magazzino ? "text-danger" : ""}`}
                    >
                      {p.quantita_disponibile}
                    </span>
                  </td>
                  <td>
                    {p.quantita_disponibile < p.soglia_minima_di_magazzino ? (
                      <span className="badge bg-danger">
                        Sotto la soglia minima
                      </span>
                    ) : (
                      <span className="badge bg-success">OK</span>
                    )}
                  </td>
                  <td className="text-end">
                    <button
                      className="btn btn-sm btn-outline-primary me-2"
                      onClick={() =>
                        handleAggiornamentoStock(p.Id, p.quantita_disponibile)
                      }
                    >
                      Aggiorna Stock
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
