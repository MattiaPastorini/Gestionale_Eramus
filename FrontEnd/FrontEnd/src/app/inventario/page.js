"use client";
import { useEffect, useState } from "react";
import api from "@/services/api";

export default function Inventario() {
  const [prodotti, setProdotti] = useState([]);
  const [filtroTipo, setFiltroTipo] = useState("");
  const [search, setSearch] = useState("");
  const [showModal, setShowModal] = useState(false);
  const [nuovoProdotto, setNuovoProdotto] = useState({
    nome_oggetto: "",
    descrizione: "",
    prezzo_unitario: 0,
    quantita_disponibile: 0,
    soglia_minima: 5,
    tipo_prodotto_id: "",
  });
  const [tipi, setTipi] = useState([]);

  useEffect(() => {
    const fetchTipi = async () => {
      try {
        const res = await api.get("/inventario/tipi");
        console.log("Dati arrivati:", res.data);
        setTipi(res.data || []);
      } catch (err) {
        console.error("Errore caricamento tipi", err);
      }
    };
    fetchTipi();
  }, []);

  const handleSave = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        ...nuovoProdotto,
        prezzo_unitario: parseFloat(nuovoProdotto.prezzo_unitario),
        quantita_disponibile: parseInt(nuovoProdotto.quantita_disponibile),
        soglia_minima: parseInt(nuovoProdotto.soglia_minima),
      };
      await api.post("/inventario/prodotti", payload);
      setShowModal(false);
      fetchProdotti();
      alert("Salvato!");
    } catch (err) {
      alert("Errore 400: Controlla che la categoria sia selezionata");
    }
  };

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

  const handleElimina = async (id, nome) => {
    if (!confirm(`Eliminare "${nome}" dal magazzino?`)) return;

    try {
      await api.delete(`/inventario/prodotti/${id}`);
      fetchProdotti();
      alert("Prodotto eliminato!");
    } catch (err) {
      console.error("Errore eliminazione:", err);
      alert("Errore eliminazione");
    }
  };

  return (
    <div className="container">
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h2 className="fw-bold text-primary">Gestione Inventario</h2>
        <button className="btn btn-success" onClick={() => setShowModal(true)}>
          Aggiungi prodotto
        </button>
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
                <tr key={p.ID}>
                  <td>
                    <strong>{p.NomeOggetto}</strong>
                    <br />
                    <small className="text-muted">{p.Descrizione}</small>
                  </td>
                  <td>
                    <span className="badge bg-secondary">
                      {p.TipoProdotto.CorpoMessaggio}
                    </span>
                  </td>
                  <td>
                    € {p.PrezzoUnitario ? p.PrezzoUnitario.toFixed(2) : "0.00"}
                  </td>
                  <td>
                    <span
                      className={`fw-bold ${p.quantita_disponibile < p.soglia_minima ? "text-danger" : ""}`}
                    >
                      {p.QuantitaDisponibile}
                    </span>
                  </td>
                  <td>
                    {p.QuantitaDisponibile < p.SogliaMinimaDiMagazzino ? (
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
                        handleAggiornamentoStock(
                          p.id || p.ID,
                          p.QuantitaDisponibile || p.quantita_disponibile,
                        )
                      }
                    >
                      Aggiorna Stock
                    </button>

                    <button
                      className="btn btn-sm btn-outline-danger"
                      onClick={() =>
                        handleElimina(
                          p.id || p.ID,
                          p.NomeOggetto || p.nome_oggetto,
                        )
                      }
                    >
                      Elimina
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
      {showModal && (
        <div
          className="modal d-block"
          style={{ backgroundColor: "rgba(0,0,0,0.5)" }}
        >
          <div className="modal-dialog modal-lg">
            <div className="modal-content border-0 shadow">
              <div className="modal-header bg-primary text-white">
                <h5 className="modal-title">Nuovo Prodotto</h5>
                <button
                  type="button"
                  className="btn-close btn-close-white"
                  onClick={() => setShowModal(false)}
                ></button>
              </div>
              <form onSubmit={handleSave}>
                <div className="modal-body">
                  <div className="row g-3">
                    <div className="col-md-6">
                      <label className="form-label">Nome Oggetto</label>
                      <input
                        type="text"
                        className="form-control"
                        required
                        onChange={(e) =>
                          setNuovoProdotto({
                            ...nuovoProdotto,
                            nome_oggetto: e.target.value,
                          })
                        }
                      />
                    </div>
                    <div className="col-md-6">
                      <label className="form-label">Prezzo Unitario (€)</label>
                      <input
                        type="number"
                        step="0.01"
                        className="form-control"
                        required
                        onChange={(e) =>
                          setNuovoProdotto({
                            ...nuovoProdotto,
                            prezzo_unitario: parseFloat(e.target.value),
                          })
                        }
                      />
                    </div>
                    <div className="col-md-4">
                      <label className="form-label">Quantità Iniziale</label>
                      <input
                        type="number"
                        className="form-control"
                        required
                        onChange={(e) =>
                          setNuovoProdotto({
                            ...nuovoProdotto,
                            quantita_disponibile: parseInt(e.target.value),
                          })
                        }
                      />
                    </div>
                    <div className="col-md-4">
                      <label className="form-label">Soglia Minima</label>
                      <input
                        type="number"
                        className="form-control"
                        defaultValue="5"
                        onChange={(e) =>
                          setNuovoProdotto({
                            ...nuovoProdotto,
                            soglia_minima: parseInt(e.target.value),
                          })
                        }
                      />
                    </div>
                    <div className="col-md-4">
                      <label className="form-label">Categoria</label>
                      <select
                        className="form-select"
                        required
                        value={nuovoProdotto.tipo_prodotto_id}
                        onChange={(e) =>
                          setNuovoProdotto({
                            ...nuovoProdotto,
                            tipo_prodotto_id: e.target.value,
                          })
                        }
                      >
                        <option value="">Seleziona categoria...</option>

                        {tipi.map((t) => (
                          <option key={t.ID} value={t.ID}>
                            {t.CorpoMessaggio}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div className="col-12">
                      <label className="form-label">Descrizione</label>
                      <textarea
                        className="form-control"
                        rows="2"
                        onChange={(e) =>
                          setNuovoProdotto({
                            ...nuovoProdotto,
                            descrizione: e.target.value,
                          })
                        }
                      ></textarea>
                    </div>
                  </div>
                </div>
                <div className="modal-footer">
                  <button
                    type="button"
                    className="btn btn-secondary"
                    onClick={() => setShowModal(false)}
                  >
                    Annulla
                  </button>
                  <button type="submit" className="btn btn-primary">
                    Salva Prodotto
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
