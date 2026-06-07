import { useState, useEffect } from 'react';
import api from '../api';

const LEVELS = ['', 'DEBUG', 'INFO', 'WARNING', 'ERROR', 'CRITICAL'];

const toLocalInput = d => new Date(d).toISOString().slice(0, 16);

const LEV_CLS = {
  ERROR:    'bg-red-100 text-red-700',
  CRITICAL: 'bg-red-200 text-red-800',
  WARNING:  'bg-yellow-100 text-yellow-700',
  INFO:     'bg-blue-100 text-blue-700',
  DEBUG:    'bg-gray-100 text-gray-600',
};

export default function Logs() {
  const [services, setServices] = useState([]);
  const [logs,     setLogs]     = useState([]);
  const [total,    setTotal]    = useState(null);
  const [offset,   setOffset]   = useState(0);
  const [loading,  setLoading]  = useState(false);
  const [more,     setMore]     = useState(false);
  const [error,    setError]    = useState('');
  const [expanded, setExpanded] = useState(null);

  const PAGE = 50;

  const [form, setForm] = useState({
    service_id: '',
    level:      '',
    message:    '',
    from:       toLocalInput(Date.now() - 3_600_000),
    to:         toLocalInput(Date.now()),
  });
  const set = (k, v) => setForm(f => ({ ...f, [k]: v }));

  useEffect(() => {
    api.get('/services').then(r => setServices(r.data ?? []));
  }, []);

  const buildParams = (off = 0) => {
    const params = new URLSearchParams();
    if (form.service_id) params.set('service_id', form.service_id);
    if (form.level)      params.set('level',      form.level);
    if (form.message)    params.set('message',    form.message);
    if (form.from)       params.set('from', new Date(form.from).toISOString());
    if (form.to)         params.set('to',   new Date(form.to).toISOString());
    params.set('limit',  String(PAGE));
    params.set('offset', String(off));
    return params;
  };

  const search = async e => {
    e?.preventDefault();
    setLoading(true);
    setError('');
    setExpanded(null);
    setOffset(0);
    try {
      const { data } = await api.get(`/logs?${buildParams(0)}`);
      setLogs(data.logs ?? []);
      setTotal(data.total ?? 0);
    } catch {
      setError('Search failed. Verify OpenSearch is running and logs have been ingested.');
    } finally {
      setLoading(false);
    }
  };

  const loadMore = async () => {
    setMore(true);
    const nextOffset = offset + PAGE;
    try {
      const { data } = await api.get(`/logs?${buildParams(nextOffset)}`);
      setLogs(prev => [...prev, ...(data.logs ?? [])]);
      setOffset(nextOffset);
    } catch {
      setError('Failed to load more logs.');
    } finally {
      setMore(false);
    }
  };

  const svcName = id => services.find(s => s.id === id)?.name ?? id?.slice(0, 8) ?? '—';

  return (
    <div className="space-y-6 max-w-6xl">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Logs</h1>
        <p className="text-sm text-gray-500 mt-1">Search indexed logs in OpenSearch</p>
      </div>

      <form onSubmit={search} className="bg-white rounded-xl border p-5 space-y-4">
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1.5">From</label>
            <input
              type="datetime-local"
              className="input"
              value={form.from}
              onChange={e => set('from', e.target.value)}
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1.5">To</label>
            <input
              type="datetime-local"
              className="input"
              value={form.to}
              onChange={e => set('to', e.target.value)}
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1.5">Service</label>
            <select className="input" value={form.service_id} onChange={e => set('service_id', e.target.value)}>
              <option value="">All services</option>
              {services.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
            </select>
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1.5">Level</label>
            <select className="input" value={form.level} onChange={e => set('level', e.target.value)}>
              {LEVELS.map(l => <option key={l} value={l}>{l || 'All levels'}</option>)}
            </select>
          </div>
        </div>
        <div className="flex gap-3">
          <div className="flex-1">
            <label className="block text-xs font-medium text-gray-600 mb-1.5">Message contains</label>
            <input
              className="input"
              value={form.message}
              placeholder="Search in message text…"
              onChange={e => set('message', e.target.value)}
            />
          </div>
          <div className="flex items-end">
            <button type="submit" disabled={loading} className="btn-primary h-[38px] px-6">
              {loading ? 'Searching…' : 'Search'}
            </button>
          </div>
        </div>
      </form>

      {error && <p className="text-red-500 text-sm">{error}</p>}

      {total !== null && !loading && (
        <div className="flex items-center gap-3">
          <p className="text-sm text-gray-500">
            Showing <span className="font-semibold text-gray-900">{logs.length}</span> of{' '}
            <span className="font-semibold text-gray-900">{total.toLocaleString()}</span> matching logs
          </p>
        </div>
      )}

      {logs.length > 0 && (
        <div className="bg-white rounded-xl border overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                {['Timestamp', 'Level', 'Service', 'Message', 'Metadata'].map(h => (
                  <th key={h} className="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {logs.map((log, i) => (
                <tr key={i} className="hover:bg-gray-50 transition-colors">
                  <td className="px-4 py-3 text-gray-500 whitespace-nowrap text-xs font-mono">
                    {log.timestamp ? new Date(log.timestamp).toLocaleString() : '—'}
                  </td>
                  <td className="px-4 py-3">
                    <span className={`text-xs font-semibold px-2 py-0.5 rounded-full ${LEV_CLS[log.level] ?? 'bg-gray-100 text-gray-600'}`}>
                      {log.level ?? '—'}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-gray-600 text-xs">{svcName(log.service_id)}</td>
                  <td className="px-4 py-3 text-gray-800 max-w-sm">
                    <p className="truncate">{log.message}</p>
                  </td>
                  <td className="px-4 py-3">
                    {log.metadata && Object.keys(log.metadata).length > 0 ? (
                      <div>
                        <button
                          onClick={() => setExpanded(expanded === i ? null : i)}
                          className="text-xs text-blue-600 hover:underline"
                        >
                          {expanded === i ? 'hide' : 'view'}
                        </button>
                        {expanded === i && (
                          <pre className="mt-1 text-xs bg-gray-50 border rounded p-2 max-w-xs overflow-auto text-gray-700 whitespace-pre-wrap">
                            {JSON.stringify(log.metadata, null, 2)}
                          </pre>
                        )}
                      </div>
                    ) : (
                      <span className="text-gray-300 text-xs">—</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {total === 0 && !loading && (
        <div className="bg-white rounded-xl border p-12 text-center">
          <p className="text-gray-400">No logs found for the selected filters.</p>
          <p className="text-sm text-gray-400 mt-1">Try expanding the time range or clearing filters.</p>
        </div>
      )}

      {/* Load more */}
      {logs.length > 0 && logs.length < total && (
        <div className="text-center">
          <button
            onClick={loadMore}
            disabled={more}
            className="text-sm px-5 py-2 rounded-lg border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-40"
          >
            {more ? 'Loading…' : `Load more (${total - logs.length} remaining)`}
          </button>
        </div>
      )}
    </div>
  );
}
