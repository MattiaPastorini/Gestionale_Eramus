"use client";
import { useEffect, useState } from "react";
import api from "@/services/api";

export default function UsersPage() {
  const [utenti, setUtenti] = useState([]);
  const [filtroRuolo, setFiltroRuolo] = useState("");
  const [search, setSearch] = useState("");
  const [showModal, setShowModal] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [formData, setFormData] = useState({
    username: "",
    email: "",
    nome: "",
    cognome: "",
    data_nascita: "",
    ruolo_id: "",
    password: "",
  });
  const [ruoli, setRuoli] = useState([]);

  useEffect(() => {
    const fetchRuoli = async () => {
      try {
        const res = await api.get("/utenti/ruoli");
        console.log("API Response completa:", res);
        console.log("res.data:", res.data);

        const ruoliArray = Array.isArray(res.data.ruoli)
          ? res.data.ruoli
          : res.data.ruoli?.data || [];

        console.log(
          "Ruoli SETTATI:",
          ruoliArray,
          "Lunghezza:",
          ruoliArray.length,
        );
        setRuoli(ruoliArray);
      } catch (err) {
        console.error("Errore caricamento ruoli", err);
        setRuoli([]);
      }
    };
    fetchRuoli();
  }, []);

  const fetchUtenti = async () => {
    try {
      const params = new URLSearchParams({
        search: search || "",
        ruolo: filtroRuolo || "",
      });
      const res = await api.get(`/utenti?${params}`);
      setUtenti(res.data.data || res.data || []);
    } catch (err) {
      console.error("Errore caricamento utenti", err);
    }
  };

  useEffect(() => {
    fetchUtenti();
  }, [search, filtroRuolo]);

  const handleSave = async (e) => {
    e.preventDefault();
    try {
      if (editingUser) {
        await api.put(`/utenti/${editingUser.ID}`, formData);
      } else {
        const payload = {
          ...formData,
          data_nascita: formData.data_nascita || null,
        };
        await api.post("/utenti", payload);
      }
      setShowModal(false);
      fetchUtenti();
      alert(editingUser ? "Utente aggiornato!" : "Utente creato!");
    } catch (err) {
      alert("Errore: " + (err.response?.data?.error || "Controlla i dati"));
    }
  };

  const handleDelete = async (id) => {
    if (confirm("Disattivare questo utente?")) {
      try {
        await api.delete(`/utenti/${id}`);
        fetchUtenti();
        alert("Utente disattivato!");
      } catch (err) {
        alert("Errore disattivazione");
      }
    }
  };

  const openEditModal = (user) => {
    setEditingUser(user);
    setFormData({
      username: user.Username || "",
      email: user.Email || "",
      nome: user.Nome || "",
      cognome: user.Cognome || "",
      data_nascita: user.DataNascita
        ? new Date(user.DataNascita).toISOString().split("T")[0]
        : "",
      ruolo_id: user.RuoloID || "",
      password: "",
    });
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
    setEditingUser(null);
    setFormData({
      username: "",
      email: "",
      nome: "",
      cognome: "",
      data_nascita: "",
      ruolo_id: "",
      password: "",
    });
  };

  return (
    <div className="container">
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h2 className="fw-bold text-primary">Gestione Utenti</h2>
        <button
          className="btn btn-success"
          onClick={() => {
            closeModal();
            setShowModal(true);
          }}
        >
          Aggiungi Utente
        </button>
      </div>

      {/* Filtri */}
      <div className="row g-3 mb-4">
        <div className="col-md-6">
          <input
            type="text"
            className="form-control"
            placeholder="Cerca per username o email..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        <div className="col-md-6">
          <select
            className="form-select"
            value={filtroRuolo}
            onChange={(e) => setFiltroRuolo(e.target.value)}
          >
            <option value="">Tutti i ruoli</option>
            <option value="Admin">Admin</option>
            <option value="Operatore">Operatore</option>
          </select>
        </div>
      </div>

      {/* Tabella */}
      <div className="card shadow-sm border-0">
        <div className="table-responsive">
          <table className="table table-hover align-middle">
            <thead className="table-light">
              <tr>
                <th>Username</th>
                <th>Email</th>
                <th>Nome Completo</th>
                <th>Ruolo</th>
                <th>Stato</th>
                <th className="text-end">Azioni</th>
              </tr>
            </thead>
            <tbody>
              {utenti.length === 0 ? (
                <tr>
                  <td colSpan="6" className="text-center py-4 text-muted">
                    Nessun utente trovato
                  </td>
                </tr>
              ) : (
                utenti.map((u) => (
                  <tr key={u.ID}>
                    <td>
                      <strong>{u.Username}</strong>
                    </td>
                    <td>{u.Email}</td>
                    <td>
                      {u.Nome} {u.Cognome}
                      <br />
                      <small className="text-muted">
                        {u.DataNascita
                          ? new Date(u.DataNascita).toLocaleDateString("it-IT")
                          : "N/D"}
                      </small>
                    </td>
                    <td>
                      <span
                        className={`badge fs-6 ${
                          u.Ruolo?.NomeRuolo === "Admin"
                            ? "bg-danger"
                            : "bg-primary"
                        }`}
                      >
                        {u.Ruolo?.NomeRuolo || "N/D"}
                      </span>
                    </td>
                    <td>
                      <span
                        className={`badge fs-6 ${
                          u.StatoAccount === "Attivo"
                            ? "bg-success"
                            : "bg-secondary"
                        }`}
                      >
                        {u.StatoAccount}
                      </span>
                    </td>
                    <td className="text-end">
                      <div className="btn-group btn-group-sm">
                        <button
                          className="btn btn-outline-warning me-1"
                          onClick={() => openEditModal(u)}
                          title="Modifica"
                        >
                          Modifica
                        </button>
                        <button
                          className="btn btn-outline-danger"
                          onClick={() => handleDelete(u.ID)}
                          title="Disattiva"
                        >
                          Disattiva
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Modal Crea/Modifica */}
      {showModal && (
        <div
          className="modal d-block"
          style={{ backgroundColor: "rgba(0,0,0,0.5)" }}
        >
          <div className="modal-dialog modal-lg">
            <div className="modal-content border-0 shadow">
              <div className="modal-header bg-primary text-white">
                <h5 className="modal-title">
                  {editingUser ? "Modifica Utente" : "Nuovo Utente"}
                </h5>
                <button
                  type="button"
                  className="btn-close btn-close-white"
                  onClick={closeModal}
                ></button>
              </div>
              <form onSubmit={handleSave}>
                <div className="modal-body">
                  <div className="row g-3">
                    <div className="col-md-6">
                      <label className="form-label">
                        Username{" "}
                        {editingUser && <small>(non modificabile)</small>}
                      </label>
                      <input
                        type="text"
                        className="form-control"
                        required
                        value={formData.username}
                        onChange={(e) =>
                          setFormData({ ...formData, username: e.target.value })
                        }
                        disabled={!!editingUser}
                      />
                    </div>
                    <div className="col-md-6">
                      <label className="form-label">Email</label>
                      <input
                        type="email"
                        className="form-control"
                        required
                        value={formData.email}
                        onChange={(e) =>
                          setFormData({ ...formData, email: e.target.value })
                        }
                      />
                    </div>
                    <div className="col-md-6">
                      <label className="form-label">
                        Password {editingUser && "(opzionale)"}
                      </label>
                      <input
                        type="password"
                        className="form-control"
                        value={formData.password}
                        onChange={(e) =>
                          setFormData({ ...formData, password: e.target.value })
                        }
                        required={!editingUser}
                      />
                    </div>
                    <div className="col-md-6">
                      <label className="form-label">Ruolo</label>
                      <select
                        className="form-select"
                        required
                        value={formData.ruolo_id}
                        onChange={(e) =>
                          setFormData({ ...formData, ruolo_id: e.target.value })
                        }
                      >
                        <option value="">Seleziona ruolo...</option>
                        {Array.isArray(ruoli) &&
                          ruoli.map((r) => (
                            <option key={r.ID} value={r.ID}>
                              {r.NomeRuolo}
                            </option>
                          ))}
                      </select>
                    </div>
                    <div className="col-md-6">
                      <label className="form-label">Nome</label>
                      <input
                        type="text"
                        className="form-control"
                        value={formData.nome}
                        onChange={(e) =>
                          setFormData({ ...formData, nome: e.target.value })
                        }
                      />
                    </div>
                    <div className="col-md-6">
                      <label className="form-label">Cognome</label>
                      <input
                        type="text"
                        className="form-control"
                        value={formData.cognome}
                        onChange={(e) =>
                          setFormData({ ...formData, cognome: e.target.value })
                        }
                      />
                    </div>
                    <div className="col-12">
                      <label className="form-label">Data di Nascita</label>
                      <input
                        type="date"
                        className="form-control"
                        value={formData.data_nascita}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            data_nascita: e.target.value,
                          })
                        }
                      />
                    </div>
                  </div>
                </div>
                <div className="modal-footer">
                  <button
                    type="button"
                    className="btn btn-secondary"
                    onClick={closeModal}
                  >
                    Annulla
                  </button>
                  <button type="submit" className="btn btn-primary">
                    {editingUser ? "Aggiorna Utente" : "Crea Utente"}
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
