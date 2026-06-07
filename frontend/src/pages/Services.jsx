import { useState, useEffect } from 'react';
import api, { jwtPayload } from '../api';

const BLANK = { name: '', description: '' };

export default function Services() {
  const [services, setServices] = useState([]);
  const [modal,    setModal]    = useState(null); // null | 'new' | service
  const [form,     setForm]     = useState(BLANK);
  const [saving,   setSaving]   = useState(false);
  const [loading,  setLoading]  = useState(true);
  const [error,    setError]    = useState('');
  const { role } = jwtPayload();
  const isAdmin = role === 'ADMIN';

  const load = () => {
    setLoading(true);
    api.get('/services')
      .then(r => setServices(r.data ?? []))
      .catch(() => setError('Failed to load services.'))
      .finally(() => setLoading(false));
  };

  useEffect(load, []);

  const openNew = () => { setForm(BLANK); setModal('new'); };
  const openEdit = s => { setForm({ name: s.name, description: s.description ?? '' }); setModal(s); };
  const closeModal = () => setModal(null);
  const set = (k, v) => setForm(f => ({ ...f, [k]: v }));

  const save = async e => {
    e.preventDefault();
    setSaving(true);
    try {
      if (modal === 'new') {
        await api.post('/services', form);
      } else {
        await api.put(`/services/${modal.id}`, form);
      }
      closeModal();
      load();
    } catch (err) {
      alert(err.response?.data?.error ?? 'Failed to save service.');
    } finally {
      setSaving(false);
    }
  };

  const remove = async id => {
    if (!confirm('Delete this service? This will also remove associated rules and alerts.')) return;
    try { await api.delete(`/services/${id}`); load(); }
    catch (err) { alert(err.response?.data?.error ?? 'Failed to delete service.'); }
  };

  if (loading) return <p className="text-gray-400 text-sm animate-pulse">Loading…</p>;
  if (error)   return <p className="text-red-500 text-sm">{error}</p>;

  return (
    <div className="space-y-6 max-w-5xl">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Services</h1>
          <p className="text-sm text-gray-500 mt-1">
            {services.filter(s => s.is_active).length} active · {services.length} total registered
          </p>
        </div>
        {isAdmin && (
          <button onClick={openNew} className="btn-primary flex items-center gap-2">
            <span className="text-lg leading-none">+</span> New Service
          </button>
        )}
      </div>

      <div className="bg-white rounded-xl border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              {['Name', 'Description', 'Status', 'Registered', ...(isAdmin ? ['Actions'] : [])].map(h => (
                <th key={h} className="px-5 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {services.length === 0 ? (
              <tr>
                <td colSpan={isAdmin ? 5 : 4} className="px-5 py-12 text-center text-gray-400">
                  No services registered yet.{isAdmin && ' Click "New Service" to create one.'}
                </td>
              </tr>
            ) : services.map(s => (
              <tr key={s.id} className="hover:bg-gray-50 transition-colors">
                <td className="px-5 py-3.5">
                  <p className="font-medium text-gray-900">{s.name}</p>
                  <p className="text-xs text-gray-400 font-mono mt-0.5">{s.id}</p>
                </td>
                <td className="px-5 py-3.5 text-gray-600">{s.description || '—'}</td>
                <td className="px-5 py-3.5">
                  <span className={`inline-block text-xs font-semibold px-2.5 py-1 rounded-full ${
                    s.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'
                  }`}>
                    {s.is_active ? 'Active' : 'Inactive'}
                  </span>
                </td>
                <td className="px-5 py-3.5 text-gray-500">
                  {s.created_at ? new Date(s.created_at).toLocaleDateString() : '—'}
                </td>
                {isAdmin && (
                  <td className="px-5 py-3.5">
                    <div className="flex gap-2">
                      <button
                        onClick={() => openEdit(s)}
                        className="text-xs px-2.5 py-1 rounded-md bg-gray-50 text-gray-700 border border-gray-200 hover:bg-gray-100 font-medium"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => remove(s.id)}
                        className="text-xs px-2.5 py-1 rounded-md bg-red-50 text-red-700 border border-red-200 hover:bg-red-100 font-medium"
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Log ingestion info */}
      <div className="bg-blue-50 border border-blue-200 rounded-xl p-5">
        <h3 className="text-sm font-semibold text-blue-800 mb-1">Log Ingestion</h3>
        <p className="text-sm text-blue-700">
          Each service uses a scoped{' '}
          <code className="bg-blue-100 px-1 rounded font-mono text-xs">X-API-KEY</code> header
          for log ingestion via{' '}
          <code className="bg-blue-100 px-1 rounded font-mono text-xs">POST /api/v1/logs</code>.
          API keys are provisioned by an admin using the seed script.
        </p>
      </div>

      {/* Create / Edit modal */}
      {modal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md">
            <div className="px-6 pt-6 pb-4 border-b border-gray-100">
              <h2 className="text-lg font-semibold text-gray-900">
                {modal === 'new' ? 'New Service' : `Edit: ${modal.name}`}
              </h2>
            </div>
            <form onSubmit={save} className="px-6 py-5 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Service name</label>
                <input
                  required
                  className="input"
                  placeholder="payment-service"
                  value={form.name}
                  onChange={e => set('name', e.target.value)}
                  autoFocus
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Description <span className="text-gray-400 font-normal">(optional)</span></label>
                <input
                  className="input"
                  placeholder="Handles payment processing"
                  value={form.description}
                  onChange={e => set('description', e.target.value)}
                />
              </div>
              <div className="flex justify-end gap-3 pt-2">
                <button type="button" onClick={closeModal} className="btn-ghost">Cancel</button>
                <button type="submit" disabled={saving} className="btn-primary">
                  {saving ? 'Saving…' : modal === 'new' ? 'Create service' : 'Save changes'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
