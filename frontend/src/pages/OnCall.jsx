import { useState, useEffect, useMemo } from 'react';
import api, { jwtPayload } from '../api';

const BLANK = { service_id: '', user_id: '', start_time: '', end_time: '' };

function toLocalInput(iso) {
  if (!iso) return '';
  const d = new Date(iso);
  const pad = n => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

function toRFC3339(localStr) {
  return new Date(localStr).toISOString();
}

function scheduleStatus(s) {
  const now = new Date();
  const start = new Date(s.start_time);
  const end = new Date(s.end_time);
  if (now >= start && now <= end) return 'active';
  if (now < start) return 'upcoming';
  return 'ended';
}

const STATUS_BADGE = {
  active:   'bg-green-100 text-green-700 border border-green-300',
  upcoming: 'bg-blue-100 text-blue-700 border border-blue-200',
  ended:    'bg-gray-100 text-gray-500 border border-gray-200',
};

export default function OnCall() {
  const { role } = jwtPayload();
  const isAdmin = role === 'ADMIN';

  const [services,   setServices]   = useState([]);
  const [users,      setUsers]      = useState([]);
  const [schedules,  setSchedules]  = useState([]);
  const [filterSvc,  setFilterSvc]  = useState('all');
  const [modal,      setModal]      = useState(null); // null | 'new' | schedule-object
  const [form,       setForm]       = useState(BLANK);
  const [saving,     setSaving]     = useState(false);
  const [loading,    setLoading]    = useState(true);
  const [error,      setError]      = useState('');

  const load = async () => {
    setLoading(true);
    try {
      const [svcRes, schedRes] = await Promise.all([
        api.get('/services'),
        api.get('/admin/oncall/schedules'),
      ]);
      setServices(svcRes.data ?? []);
      setSchedules(schedRes.data ?? []);
      if (isAdmin) {
        const usrRes = await api.get('/admin/developers');
        setUsers(usrRes.data ?? []);
      }
    } catch {
      setError('Failed to load on-call data.');
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

  const openNew = () => {
    setForm(BLANK);
    setModal('new');
  };

  const openEdit = s => {
    setForm({
      service_id: s.service_id,
      user_id:    s.user_id,
      start_time: toLocalInput(s.start_time),
      end_time:   toLocalInput(s.end_time),
    });
    setModal(s);
  };

  const closeModal = () => setModal(null);

  const save = async e => {
    e.preventDefault();
    setSaving(true);
    try {
      const payload = {
        service_id: form.service_id,
        user_id:    form.user_id,
        start_time: toRFC3339(form.start_time),
        end_time:   toRFC3339(form.end_time),
      };
      if (modal === 'new') {
        await api.post('/admin/oncall/schedules', payload);
      } else {
        await api.put(`/admin/oncall/schedules/${modal.id}`, payload);
      }
      closeModal();
      load();
    } catch (err) {
      alert(err.response?.data?.error ?? 'Failed to save schedule.');
    } finally {
      setSaving(false);
    }
  };

  const remove = async s => {
    const user = users.find(u => u.id === s.user_id);
    const svc  = services.find(sv => sv.id === s.service_id);
    if (!confirm(`Delete on-call schedule for ${user?.name ?? 'this user'} on ${svc?.name ?? 'this service'}?`)) return;
    try {
      await api.delete(`/admin/oncall/schedules/${s.id}`);
      load();
    } catch (err) {
      alert(err.response?.data?.error ?? 'Failed to delete schedule.');
    }
  };

  const serviceMap = useMemo(() => Object.fromEntries(services.map(s => [s.id, s])), [services]);
  const userMap    = useMemo(() => Object.fromEntries(users.map(u => [u.id, u])),    [users]);

  const visible = filterSvc === 'all'
    ? schedules
    : schedules.filter(s => s.service_id === filterSvc);

  const sorted = [...visible].sort((a, b) => new Date(b.start_time) - new Date(a.start_time));

  if (loading) return <p className="text-gray-400 text-sm animate-pulse">Loading…</p>;
  if (error)   return <p className="text-red-500 text-sm">{error}</p>;

  const activeByService = {};
  schedules.forEach(s => {
    if (scheduleStatus(s) === 'active') {
      if (!activeByService[s.service_id]) activeByService[s.service_id] = [];
      activeByService[s.service_id].push(s);
    }
  });

  return (
    <div className="space-y-6 max-w-5xl">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">On-Call Schedules</h1>
          <p className="text-sm text-gray-500 mt-1">Manage who receives alert emails for each service</p>
        </div>
        <button onClick={openNew} className="btn-primary flex items-center gap-2">
          <span className="text-lg leading-none">+</span> New Schedule
        </button>
      </div>

      {/* Current on-call summary cards */}
      {services.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {services.map(svc => {
            const actives = activeByService[svc.id] ?? [];
            return (
              <div key={svc.id} className={`rounded-xl border p-4 ${actives.length > 0 ? 'bg-green-50 border-green-200' : 'bg-gray-50 border-gray-200'}`}>
                <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">{svc.name}</p>
                {actives.length > 0 ? (
                  <div className="space-y-2">
                    {actives.map(a => {
                      const u = userMap[a.user_id];
                      return (
                        <div key={a.id}>
                          <p className="text-sm font-semibold text-green-800">{u?.name ?? a.user_id}</p>
                          <p className="text-xs text-green-600">{u?.email}</p>
                          <p className="text-xs text-gray-500">Until {new Date(a.end_time).toLocaleString()}</p>
                        </div>
                      );
                    })}
                  </div>
                ) : (
                  <p className="text-sm text-gray-400 italic">No one on-call</p>
                )}
              </div>
            );
          })}
        </div>
      )}

      {/* Filter by service */}
      <div className="flex items-center gap-2 flex-wrap">
        <button
          onClick={() => setFilterSvc('all')}
          className={`text-xs px-3 py-1.5 rounded-full font-medium border transition-colors ${
            filterSvc === 'all' ? 'bg-blue-600 text-white border-blue-600' : 'bg-white text-gray-600 border-gray-300 hover:border-gray-400'
          }`}
        >
          All Services
        </button>
        {services.map(svc => (
          <button
            key={svc.id}
            onClick={() => setFilterSvc(svc.id)}
            className={`text-xs px-3 py-1.5 rounded-full font-medium border transition-colors ${
              filterSvc === svc.id ? 'bg-blue-600 text-white border-blue-600' : 'bg-white text-gray-600 border-gray-300 hover:border-gray-400'
            }`}
          >
            {svc.name}
          </button>
        ))}
      </div>

      {/* Schedule table */}
      <div className="bg-white rounded-xl border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              {['Service', 'Developer', 'Start', 'End', 'Status', 'Actions'].map(h => (
                <th key={h} className="px-5 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {sorted.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-5 py-12 text-center text-gray-400">
                  No schedules yet. Click &ldquo;New Schedule&rdquo; to add one.
                </td>
              </tr>
            ) : sorted.map(s => {
              const status = scheduleStatus(s);
              const svc    = serviceMap[s.service_id];
              const user   = userMap[s.user_id];
              return (
                <tr key={s.id} className={`hover:bg-gray-50 transition-colors ${status === 'active' ? 'bg-green-50/30' : ''}`}>
                  <td className="px-5 py-3.5 font-medium text-gray-900">{svc?.name ?? s.service_id}</td>
                  <td className="px-5 py-3.5">
                    <p className="font-medium text-gray-900">{user?.name ?? s.user_id}</p>
                    <p className="text-xs text-gray-400">{user?.email}</p>
                  </td>
                  <td className="px-5 py-3.5 text-gray-600 text-xs">{new Date(s.start_time).toLocaleString()}</td>
                  <td className="px-5 py-3.5 text-gray-600 text-xs">{new Date(s.end_time).toLocaleString()}</td>
                  <td className="px-5 py-3.5">
                    <span className={`inline-block text-xs font-semibold px-2.5 py-1 rounded-full capitalize ${STATUS_BADGE[status]}`}>
                      {status}
                    </span>
                  </td>
                  <td className="px-5 py-3.5">
                    <div className="flex gap-2">
                      <button
                        onClick={() => openEdit(s)}
                        className="text-xs px-2.5 py-1 rounded-md bg-gray-50 text-gray-700 border border-gray-200 hover:bg-gray-100 font-medium"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => remove(s)}
                        className="text-xs px-2.5 py-1 rounded-md bg-red-50 text-red-700 border border-red-200 hover:bg-red-100 font-medium"
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      <div className="bg-blue-50 border border-blue-200 rounded-xl p-5">
        <h3 className="text-sm font-semibold text-blue-800 mb-1">How alerts work</h3>
        <p className="text-sm text-blue-700">
          When an alert fires, the system looks up who is currently on-call for that service and sends an email to their address.
          Make sure every service has an active schedule at all times to avoid missed notifications.
        </p>
      </div>

      {/* Create / Edit modal */}
      {modal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md">
            <div className="px-6 pt-6 pb-4 border-b border-gray-100">
              <h2 className="text-lg font-semibold text-gray-900">
                {modal === 'new' ? 'New On-Call Schedule' : 'Edit Schedule'}
              </h2>
              <p className="text-sm text-gray-500 mt-1">
                Alert emails will be sent to the selected developer when they are on-call.
              </p>
            </div>
            <form onSubmit={save} className="px-6 py-5 space-y-4">
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
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">On-Call Developer</label>
                <select
                  required
                  className="input"
                  value={form.user_id}
                  onChange={e => set('user_id', e.target.value)}
                >
                  <option value="">Select a developer…</option>
                  {users.map(u => (
                    <option key={u.id} value={u.id}>{u.name} ({u.email})</option>
                  ))}
                </select>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1.5">Start</label>
                  <input
                    type="datetime-local"
                    required
                    className="input"
                    value={form.start_time}
                    onChange={e => set('start_time', e.target.value)}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1.5">End</label>
                  <input
                    type="datetime-local"
                    required
                    className="input"
                    value={form.end_time}
                    onChange={e => set('end_time', e.target.value)}
                  />
                </div>
              </div>
              <div className="flex justify-end gap-3 pt-2">
                <button type="button" onClick={closeModal} className="btn-ghost">Cancel</button>
                <button type="submit" disabled={saving} className="btn-primary">
                  {saving ? 'Saving…' : modal === 'new' ? 'Create schedule' : 'Save changes'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
