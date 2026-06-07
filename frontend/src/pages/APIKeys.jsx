import { useState, useEffect } from 'react';
import api from '../api';
import { jwtPayload } from '../api';

export default function APIKeys() {
  const { role } = jwtPayload();
  const isAdmin = role === 'ADMIN';

  const [keys,     setKeys]     = useState([]);
  const [services, setServices] = useState([]);
  const [loading,  setLoading]  = useState(true);
  const [error,    setError]    = useState('');
  const [modal,    setModal]    = useState(false);
  const [form,     setForm]     = useState({ service_id: '', name: '' });
  const [saving,   setSaving]   = useState(false);
  const [newKey,   setNewKey]   = useState(null); // revealed key after creation

  const load = async () => {
    setLoading(true);
    try {
      const [keysRes, svcRes] = await Promise.all([
        api.get('/admin/api-keys'),
        api.get('/services'),
      ]);
      setKeys(keysRes.data ?? []);
      setServices(svcRes.data ?? []);
    } catch {
      setError('Failed to load API keys.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  if (!isAdmin) {
    return (
      <div className="text-center py-16">
        <p className="text-gray-500">You don&apos;t have permission to view this page.</p>
      </div>
    );
  }

  const set = (k, v) => setForm(f => ({ ...f, [k]: v }));

  const create = async e => {
    e.preventDefault();
    setSaving(true);
    try {
      const { data } = await api.post('/admin/api-keys', form);
      setNewKey(data.key);
      setModal(false);
      setForm({ service_id: '', name: '' });
      load();
    } catch (err) {
      alert(err.response?.data?.error ?? 'Failed to create API key.');
    } finally {
      setSaving(false);
    }
  };

  const revoke = async key => {
    if (!confirm(`Revoke key "${key.name}"? Any application using it will lose access immediately.`)) return;
    try {
      await api.delete(`/admin/api-keys/${key.id}`);
      load();
    } catch (err) {
      alert(err.response?.data?.error ?? 'Failed to revoke key.');
    }
  };

  const svcName = id => services.find(s => s.id === id)?.name ?? id?.slice(0, 8) ?? '—';

  if (loading) return <p className="text-gray-400 text-sm animate-pulse">Loading…</p>;
  if (error)   return <p className="text-red-500 text-sm">{error}</p>;

  return (
    <div className="space-y-6 max-w-4xl">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">API Keys</h1>
          <p className="text-sm text-gray-500 mt-1">Keys used by services to ingest logs via <code className="bg-gray-100 px-1 rounded text-xs">POST /api/v1/logs</code></p>
        </div>
        <button onClick={() => setModal(true)} className="btn-primary flex items-center gap-2">
          <span className="text-lg leading-none">+</span> New API Key
        </button>
      </div>

      {/* Newly created key banner */}
      {newKey && (
        <div className="bg-green-50 border border-green-300 rounded-xl p-5">
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1 min-w-0">
              <p className="text-sm font-semibold text-green-800 mb-1">API key created — copy it now</p>
              <p className="text-xs text-green-700 mb-3">This key will not be shown again. Store it securely.</p>
              <code className="block bg-white border border-green-200 rounded-lg px-4 py-2.5 text-sm font-mono text-gray-800 break-all select-all">
                {newKey}
              </code>
            </div>
            <button
              onClick={() => { navigator.clipboard.writeText(newKey); }}
              className="shrink-0 text-xs px-3 py-1.5 rounded-md bg-green-100 text-green-800 border border-green-300 hover:bg-green-200 font-medium"
            >
              Copy
            </button>
          </div>
          <button onClick={() => setNewKey(null)} className="mt-3 text-xs text-green-600 hover:underline">
            I&apos;ve saved it — dismiss
          </button>
        </div>
      )}

      <div className="bg-white rounded-xl border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              {['Name', 'Service', 'Status', 'Created', 'Actions'].map(h => (
                <th key={h} className="px-5 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {keys.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-5 py-12 text-center text-gray-400">
                  No API keys yet. Click &ldquo;New API Key&rdquo; to create one.
                </td>
              </tr>
            ) : keys.map(k => (
              <tr key={k.id} className={`hover:bg-gray-50 transition-colors ${!k.is_active ? 'opacity-50' : ''}`}>
                <td className="px-5 py-3.5">
                  <p className="font-medium text-gray-900">{k.name}</p>
                  <p className="text-xs text-gray-400 font-mono mt-0.5">{k.id.slice(0, 8)}…</p>
                </td>
                <td className="px-5 py-3.5 text-gray-600">{svcName(k.service_id)}</td>
                <td className="px-5 py-3.5">
                  <span className={`inline-block text-xs font-semibold px-2.5 py-1 rounded-full ${
                    k.is_active
                      ? 'bg-green-100 text-green-700 border border-green-200'
                      : 'bg-gray-100 text-gray-500 border border-gray-200'
                  }`}>
                    {k.is_active ? 'Active' : 'Revoked'}
                  </span>
                </td>
                <td className="px-5 py-3.5 text-gray-500 text-xs">
                  {new Date(k.created_at).toLocaleDateString()}
                </td>
                <td className="px-5 py-3.5">
                  {k.is_active && (
                    <button
                      onClick={() => revoke(k)}
                      className="text-xs px-2.5 py-1 rounded-md bg-red-50 text-red-700 border border-red-200 hover:bg-red-100 font-medium"
                    >
                      Revoke
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="bg-amber-50 border border-amber-200 rounded-xl p-5">
        <h3 className="text-sm font-semibold text-amber-800 mb-1">Usage</h3>
        <p className="text-sm text-amber-700 mb-2">Send logs with the key in the <code className="bg-amber-100 px-1 rounded text-xs">X-API-KEY</code> header:</p>
        <pre className="text-xs bg-white border border-amber-200 rounded-lg p-3 text-gray-700 overflow-x-auto">{`curl -X POST https://your-domain/api/v1/logs \\
  -H "X-API-KEY: ls_..." \\
  -H "Content-Type: application/json" \\
  -d '{"level":"ERROR","message":"payment failed"}'`}</pre>
      </div>

      {/* Create modal */}
      {modal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md">
            <div className="px-6 pt-6 pb-4 border-b border-gray-100">
              <h2 className="text-lg font-semibold text-gray-900">New API Key</h2>
              <p className="text-sm text-gray-500 mt-1">The key will be shown once after creation.</p>
            </div>
            <form onSubmit={create} className="px-6 py-5 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Name</label>
                <input
                  required
                  className="input"
                  placeholder="e.g. payment-service-prod"
                  value={form.name}
                  onChange={e => set('name', e.target.value)}
                  autoFocus
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Service</label>
                <select
                  required
                  className="input"
                  value={form.service_id}
                  onChange={e => set('service_id', e.target.value)}
                >
                  <option value="">Select a service…</option>
                  {services.map(s => (
                    <option key={s.id} value={s.id}>{s.name}</option>
                  ))}
                </select>
              </div>
              <div className="flex justify-end gap-3 pt-2">
                <button type="button" onClick={() => setModal(false)} className="btn-ghost">Cancel</button>
                <button type="submit" disabled={saving} className="btn-primary">
                  {saving ? 'Creating…' : 'Create key'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
