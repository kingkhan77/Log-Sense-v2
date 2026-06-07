import { useState } from 'react';
import api, { jwtPayload } from '../api';

export default function Profile() {
  const { role, user_id } = jwtPayload();

  const [form,    setForm]    = useState({ current_password: '', new_password: '', confirm: '' });
  const [saving,  setSaving]  = useState(false);
  const [success, setSuccess] = useState('');
  const [error,   setError]   = useState('');

  const set = (k, v) => { setForm(f => ({ ...f, [k]: v })); setSuccess(''); setError(''); };

  const submit = async e => {
    e.preventDefault();
    if (form.new_password !== form.confirm) {
      setError('New passwords do not match.');
      return;
    }
    setSaving(true);
    try {
      await api.put('/me/password', {
        current_password: form.current_password,
        new_password:     form.new_password,
      });
      setSuccess('Password changed successfully.');
      setForm({ current_password: '', new_password: '', confirm: '' });
    } catch (err) {
      setError(err.response?.data?.error ?? 'Failed to change password.');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="max-w-lg space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Profile</h1>
        <p className="text-sm text-gray-500 mt-1">Manage your account settings</p>
      </div>

      {/* Identity card */}
      <div className="bg-white rounded-xl border p-5 space-y-3">
        <h2 className="text-sm font-semibold text-gray-700">Account</h2>
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center text-blue-700 font-bold text-sm">
            {role?.[0] ?? '?'}
          </div>
          <div>
            <p className="text-sm font-medium text-gray-900">{user_id ?? '—'}</p>
            <span className={`inline-block text-xs font-semibold px-2 py-0.5 rounded-full ${
              role === 'ADMIN' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-600'
            }`}>
              {role ?? 'DEVELOPER'}
            </span>
          </div>
        </div>
      </div>

      {/* Change password */}
      <div className="bg-white rounded-xl border p-5">
        <h2 className="text-sm font-semibold text-gray-700 mb-4">Change password</h2>

        {success && (
          <div className="mb-4 text-sm text-green-700 bg-green-50 border border-green-200 rounded-lg px-4 py-2.5">
            {success}
          </div>
        )}
        {error && (
          <div className="mb-4 text-sm text-red-700 bg-red-50 border border-red-200 rounded-lg px-4 py-2.5">
            {error}
          </div>
        )}

        <form onSubmit={submit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">Current password</label>
            <input
              type="password"
              required
              className="input"
              value={form.current_password}
              onChange={e => set('current_password', e.target.value)}
              autoComplete="current-password"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">New password</label>
            <input
              type="password"
              required
              minLength={6}
              className="input"
              placeholder="Min. 6 characters"
              value={form.new_password}
              onChange={e => set('new_password', e.target.value)}
              autoComplete="new-password"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">Confirm new password</label>
            <input
              type="password"
              required
              minLength={6}
              className="input"
              value={form.confirm}
              onChange={e => set('confirm', e.target.value)}
              autoComplete="new-password"
            />
          </div>
          <div className="pt-1">
            <button type="submit" disabled={saving} className="btn-primary">
              {saving ? 'Saving…' : 'Change password'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
