import { useState, useEffect } from 'react';
import api from '../api';

const SEV = {
  CRITICAL: 'bg-red-100 text-red-700',
  WARNING:  'bg-yellow-100 text-yellow-700',
  INFO:     'bg-blue-100 text-blue-700',
};

const BLANK = {
  service_id: '', name: '', description: '',
  severity: 'WARNING', query: '{}',
  threshold: 5, window_minutes: 5, is_enabled: true,
};

export default function Rules() {
  const [rules,    setRules]    = useState([]);
  const [services, setServices] = useState([]);
  const [modal,    setModal]    = useState(null); // null | 'new' | rule
  const [form,     setForm]     = useState(BLANK);
  const [saving,   setSaving]   = useState(false);
  const [loading,  setLoading]  = useState(true);
  const [error,    setError]    = useState('');

  const load = () => {
    setLoading(true);
    Promise.all([api.get('/rules'), api.get('/services')])
      .then(([r, s]) => { setRules(r.data ?? []); setServices(s.data ?? []); })
      .catch(() => setError('Failed to load data.'))
      .finally(() => setLoading(false));
  };

  useEffect(load, []);

  const svcName = id => services.find(s => s.id === id)?.name ?? id.slice(0, 8) + '…';

  const openNew = () => {
    setForm({ ...BLANK, service_id: services[0]?.id ?? '' });
    setModal('new');
  };

  const openEdit = rule => {
    setForm({
      service_id:     rule.service_id,
      name:           rule.name,
      description:    rule.description ?? '',
      severity:       rule.severity,
      query:          safeStringify(rule.query),
      threshold:      rule.threshold,
      window_minutes: rule.window_minutes,
      is_enabled:     rule.is_enabled,
    });
    setModal(rule);
  };

  const closeModal = () => setModal(null);

  const set = (k, v) => setForm(f => ({ ...f, [k]: v }));

  const save = async e => {
    e.preventDefault();
    let query;
    try { query = JSON.parse(form.query); }
    catch { alert('Query must be valid JSON (e.g. {} or {"level":"ERROR"})'); return; }

    setSaving(true);
    try {
      const payload = {
        ...form,
        query,
        threshold:      Number(form.threshold),
        window_minutes: Number(form.window_minutes),
      };
      if (modal === 'new') {
        await api.post('/rules', payload);
      } else {
        await api.put(`/rules/${modal.id}`, payload);
      }
      closeModal();
      load();
    } catch (err) {
      alert(err.response?.data?.error ?? 'Failed to save rule.');
    } finally {
      setSaving(false);
    }
  };

  const remove = async id => {
    if (!confirm('Delete this rule permanently?')) return;
    try { await api.delete(`/rules/${id}`); load(); }
    catch { alert('Failed to delete rule.'); }
  };

  const toggle = async rule => {
    try { await api.put(`/rules/${rule.id}`, { is_enabled: !rule.is_enabled }); load(); }
    catch { alert('Failed to update rule.'); }
  };

  if (loading) return <p className="text-gray-400 text-sm animate-pulse">Loading…</p>;
  if (error)   return <p className="text-red-500 text-sm">{error}</p>;

  return (
    <div className="space-y-6 max-w-6xl">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Alert Rules</h1>
          <p className="text-sm text-gray-500 mt-1">{rules.filter(r => r.is_enabled).length} enabled · {rules.length} total</p>
        </div>
        <button onClick={openNew} className="btn-primary flex items-center gap-2">
          <span className="text-lg leading-none">+</span> New Rule
        </button>
      </div>

      <div className="bg-white rounded-xl border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              {['Name', 'Service', 'Severity', 'Threshold', 'Window', 'Status', 'Actions'].map(h => (
                <th key={h} className="px-5 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {rules.length === 0 ? (
              <tr><td colSpan={7} className="px-5 py-12 text-center text-gray-400">No rules configured yet.</td></tr>
            ) : rules.map(r => (
              <tr key={r.id} className="hover:bg-gray-50 transition-colors">
                <td className="px-5 py-3.5">
                  <p className="font-medium text-gray-900">{r.name}</p>
                  {r.description && <p className="text-xs text-gray-400 mt-0.5 truncate max-w-xs">{r.description}</p>}
                </td>
                <td className="px-5 py-3.5 text-gray-600">{svcName(r.service_id)}</td>
                <td className="px-5 py-3.5">
                  <span className={`text-xs font-semibold px-2 py-0.5 rounded-full ${SEV[r.severity] ?? ''}`}>{r.severity}</span>
                </td>
                <td className="px-5 py-3.5 text-gray-600 tabular-nums">{r.threshold}</td>
                <td className="px-5 py-3.5 text-gray-600">{r.window_minutes} min</td>
                <td className="px-5 py-3.5">
                  <button
                    onClick={() => toggle(r)}
                    className={`text-xs font-semibold px-2.5 py-1 rounded-full border transition-colors ${
                      r.is_enabled
                        ? 'bg-green-50 text-green-700 border-green-200 hover:bg-green-100'
                        : 'bg-gray-50 text-gray-500 border-gray-200 hover:bg-gray-100'
                    }`}
                  >
                    {r.is_enabled ? 'Enabled' : 'Disabled'}
                  </button>
                </td>
                <td className="px-5 py-3.5">
                  <div className="flex gap-2">
                    <button
                      onClick={() => openEdit(r)}
                      className="text-xs px-2.5 py-1 rounded-md bg-gray-50 text-gray-700 border border-gray-200 hover:bg-gray-100 font-medium"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => remove(r.id)}
                      className="text-xs px-2.5 py-1 rounded-md bg-red-50 text-red-700 border border-red-200 hover:bg-red-100 font-medium"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Create / Edit modal */}
      {modal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
            <div className="px-6 pt-6 pb-4 border-b border-gray-100">
              <h2 className="text-lg font-semibold text-gray-900">
                {modal === 'new' ? 'New Alert Rule' : `Edit: ${modal.name}`}
              </h2>
            </div>

            <form onSubmit={save} className="px-6 py-5 space-y-4">
              <Field label="Rule name">
                <input required className="input" value={form.name}
                  placeholder="High error rate"
                  onChange={e => set('name', e.target.value)} />
              </Field>

              <Field label="Service">
                <select required className="input" value={form.service_id}
                  onChange={e => set('service_id', e.target.value)}>
                  {services.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
                </select>
              </Field>

              <div className="grid grid-cols-2 gap-4">
                <Field label="Severity">
                  <select className="input" value={form.severity} onChange={e => set('severity', e.target.value)}>
                    {['INFO', 'WARNING', 'CRITICAL'].map(s => <option key={s}>{s}</option>)}
                  </select>
                </Field>
                <Field label="Threshold (log count)">
                  <input type="number" min="1" required className="input" value={form.threshold}
                    onChange={e => set('threshold', e.target.value)} />
                </Field>
              </div>

              <Field label="Window (minutes)">
                <input type="number" min="1" required className="input" value={form.window_minutes}
                  onChange={e => set('window_minutes', e.target.value)} />
              </Field>

              <Field label='Query (JSON — e.g. {"level":"ERROR","message_contains":"timeout"})'>
                <textarea
                  rows={3}
                  className="input font-mono text-xs"
                  value={form.query}
                  onChange={e => set('query', e.target.value)}
                  placeholder='{"level":"ERROR"}'
                />
              </Field>

              <Field label="Description (optional)">
                <input className="input" value={form.description}
                  placeholder="Brief description of this rule"
                  onChange={e => set('description', e.target.value)} />
              </Field>

              <label className="flex items-center gap-2.5 cursor-pointer">
                <input type="checkbox" className="rounded" checked={form.is_enabled}
                  onChange={e => set('is_enabled', e.target.checked)} />
                <span className="text-sm text-gray-700">Enable rule immediately</span>
              </label>

              <div className="flex justify-end gap-3 pt-2">
                <button type="button" onClick={closeModal} className="btn-ghost">Cancel</button>
                <button type="submit" disabled={saving} className="btn-primary">
                  {saving ? 'Saving…' : 'Save rule'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

function Field({ label, children }) {
  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-1.5">{label}</label>
      {children}
    </div>
  );
}

function safeStringify(q) {
  if (!q) return '{}';
  if (typeof q === 'string') return q;
  try { return JSON.stringify(q, null, 2); } catch { return '{}'; }
}
