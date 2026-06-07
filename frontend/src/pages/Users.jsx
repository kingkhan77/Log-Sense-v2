import { useState, useEffect } from 'react';
import api, { jwtPayload } from '../api';

const BLANK_CREATE = { name: '', email: '', password: '' };
const BLANK_EDIT   = { name: '', email: '' };

export default function Users() {
  const [users,   setUsers]   = useState([]);
  const [modal,   setModal]   = useState(null); // null | 'new' | user-object
  const [form,    setForm]    = useState(BLANK_CREATE);
  const [saving,  setSaving]  = useState(false);
  const [loading, setLoading] = useState(true);
  const [error,   setError]   = useState('');
  const { role } = jwtPayload();
  const isAdmin = role === 'ADMIN';

  const load = () => {
    if (!isAdmin) return;
    setLoading(true);
    api.get('/admin/developers')
      .then(r => setUsers(r.data ?? []))
      .catch(() => setError('Failed to load users.'))
      .finally(() => setLoading(false));
  };

  useEffect(load, [isAdmin]);

  if (!isAdmin) {
    return (
      <div className="text-center py-16">
        <p className="text-gray-500">You don&apos;t have permission to view this page.</p>
      </div>
    );
  }

  const set = (k, v) => setForm(f => ({ ...f, [k]: v }));

  const openNew  = () => { setForm(BLANK_CREATE); setModal('new'); };
  const openEdit = u  => { setForm({ name: u.name, email: u.email }); setModal(u); };
  const closeModal = () => setModal(null);

  const save = async e => {
    e.preventDefault();
    setSaving(true);
    try {
      if (modal === 'new') {
        await api.post('/admin/developers', form);
      } else {
        await api.put(`/admin/developers/${modal.id}`, { name: form.name, email: form.email });
      }
      closeModal();
      load();
    } catch (err) {
      alert(err.response?.data?.error ?? 'Failed to save user.');
    } finally {
      setSaving(false);
    }
  };

  const remove = async u => {
    if (!confirm(`Delete ${u.name}? They will lose access immediately.`)) return;
    try {
      await api.delete(`/admin/developers/${u.id}`);
      load();
    } catch (err) {
      alert(err.response?.data?.error ?? 'Failed to deactivate user.');
    }
  };

  if (loading) return <p className="text-gray-400 text-sm animate-pulse">Loading…</p>;
  if (error)   return <p className="text-red-500 text-sm">{error}</p>;

  return (
    <div className="space-y-6 max-w-4xl">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Users</h1>
          <p className="text-sm text-gray-500 mt-1">{users.length} developer{users.length !== 1 ? 's' : ''} in this tenant</p>
        </div>
        <button onClick={openNew} className="btn-primary flex items-center gap-2">
          <span className="text-lg leading-none">+</span> New Developer
        </button>
      </div>

      <div className="bg-white rounded-xl border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              {['Name', 'Email', 'Role', 'Actions'].map(h => (
                <th key={h} className="px-5 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {users.length === 0 ? (
              <tr>
                <td colSpan={4} className="px-5 py-12 text-center text-gray-400">
                  No developers yet. Click &ldquo;New Developer&rdquo; to invite one.
                </td>
              </tr>
            ) : users.map(u => (
              <tr key={u.id} className="hover:bg-gray-50 transition-colors">
                <td className="px-5 py-3.5">
                  <p className="font-medium text-gray-900">{u.name}</p>
                  <p className="text-xs text-gray-400 font-mono mt-0.5">{u.id}</p>
                </td>
                <td className="px-5 py-3.5 text-gray-600">{u.email}</td>
                <td className="px-5 py-3.5">
                  <span className={`inline-block text-xs font-semibold px-2.5 py-1 rounded-full ${
                    u.role === 'ADMIN'
                      ? 'bg-blue-100 text-blue-700'
                      : 'bg-gray-100 text-gray-600'
                  }`}>
                    {u.role}
                  </span>
                </td>
                <td className="px-5 py-3.5">
                  <div className="flex gap-2">
                    <button
                      onClick={() => openEdit(u)}
                      className="text-xs px-2.5 py-1 rounded-md bg-gray-50 text-gray-700 border border-gray-200 hover:bg-gray-100 font-medium"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => remove(u)}
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

      <div className="bg-amber-50 border border-amber-200 rounded-xl p-5">
        <h3 className="text-sm font-semibold text-amber-800 mb-1">About roles</h3>
        <p className="text-sm text-amber-700">
          Developers created here are assigned the <code className="bg-amber-100 px-1 rounded font-mono text-xs">DEVELOPER</code> role.
          Deactivating a user revokes access immediately without deleting their audit history.
        </p>
      </div>

      {modal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md">
            <div className="px-6 pt-6 pb-4 border-b border-gray-100">
              <h2 className="text-lg font-semibold text-gray-900">
                {modal === 'new' ? 'New Developer' : `Edit: ${modal.name}`}
              </h2>
              {modal === 'new' && (
                <p className="text-sm text-gray-500 mt-1">They will be able to log in with these credentials.</p>
              )}
            </div>
            <form onSubmit={save} className="px-6 py-5 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Full name</label>
                <input
                  required
                  className="input"
                  placeholder="Jane Smith"
                  value={form.name}
                  onChange={e => set('name', e.target.value)}
                  autoFocus
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Email</label>
                <input
                  type="email"
                  required
                  className="input"
                  placeholder="jane@company.com"
                  value={form.email}
                  onChange={e => set('email', e.target.value)}
                />
              </div>
              {modal === 'new' && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1.5">Password</label>
                  <input
                    type="password"
                    required
                    minLength={6}
                    className="input"
                    placeholder="Min. 6 characters"
                    value={form.password}
                    onChange={e => set('password', e.target.value)}
                  />
                </div>
              )}
              <div className="flex justify-end gap-3 pt-2">
                <button type="button" onClick={closeModal} className="btn-ghost">Cancel</button>
                <button type="submit" disabled={saving} className="btn-primary">
                  {saving ? 'Saving…' : modal === 'new' ? 'Create developer' : 'Save changes'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
