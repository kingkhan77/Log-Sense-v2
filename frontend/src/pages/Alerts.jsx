import { useState, useEffect } from 'react';
import api from '../api';

const STATUS = {
  OPEN:         'bg-red-100 text-red-700 border-red-200',
  ACKNOWLEDGED: 'bg-yellow-100 text-yellow-700 border-yellow-200',
  RESOLVED:     'bg-green-100 text-green-700 border-green-200',
};
const SEV = {
  CRITICAL: 'bg-red-100 text-red-700',
  WARNING:  'bg-yellow-100 text-yellow-700',
  INFO:     'bg-blue-100 text-blue-700',
};

const FILTERS = ['All', 'OPEN', 'ACKNOWLEDGED', 'RESOLVED'];

const PAGE_SIZE = 50;

export default function Alerts() {
  const [alerts,  setAlerts]  = useState([]);
  const [total,   setTotal]   = useState(0);
  const [page,    setPage]    = useState(0);
  const [filter,  setFilter]  = useState('All');
  const [loading, setLoading] = useState(true);
  const [error,   setError]   = useState('');

  const load = (p = page) => {
    setLoading(true);
    api.get(`/alerts?limit=${PAGE_SIZE}&offset=${p * PAGE_SIZE}`)
      .then(r => {
        setAlerts(r.data?.alerts ?? r.data ?? []);
        setTotal(r.data?.total ?? (r.data?.length ?? 0));
      })
      .catch(() => setError('Failed to load alerts.'))
      .finally(() => setLoading(false));
  };

  useEffect(() => { load(page); }, [page]);

  const action = async (id, verb) => {
    try {
      await api.post(`/alerts/${id}/${verb}`);
      load(page);
    } catch (err) {
      alert(err.response?.data?.error ?? `Failed to ${verb} alert.`);
    }
  };

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));
  const visible = filter === 'All' ? alerts : alerts.filter(a => a.status === filter);

  return (
    <div className="space-y-6 max-w-6xl">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Alerts</h1>
          <p className="text-sm text-gray-500 mt-1">{total} total · {alerts.filter(a => a.status === 'OPEN').length} open on this page</p>
        </div>

        {/* Filter tabs */}
        <div className="flex gap-1 bg-gray-100 p-1 rounded-lg">
          {FILTERS.map(f => (
            <button
              key={f}
              onClick={() => setFilter(f)}
              className={`px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
                filter === f
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              {f}
            </button>
          ))}
        </div>
      </div>

      {loading && <p className="text-gray-400 text-sm animate-pulse">Loading…</p>}
      {error   && <p className="text-red-500 text-sm">{error}</p>}

      {!loading && !error && (
        <div className="bg-white rounded-xl border overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                {['Title', 'Severity', 'Status', 'Count / Threshold', 'Triggered', 'Actions'].map(h => (
                  <th key={h} className="px-5 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {visible.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-5 py-12 text-center text-gray-400">
                    No alerts match this filter.
                  </td>
                </tr>
              ) : visible.map(a => (
                <tr key={a.id} className="hover:bg-gray-50 transition-colors">
                  <td className="px-5 py-3.5">
                    <p className="font-medium text-gray-900">{a.title}</p>
                    {a.description && <p className="text-xs text-gray-400 mt-0.5 truncate max-w-xs">{a.description}</p>}
                  </td>
                  <td className="px-5 py-3.5">
                    <Badge cls={SEV[a.severity]}>{a.severity}</Badge>
                  </td>
                  <td className="px-5 py-3.5">
                    <Badge cls={`border ${STATUS[a.status]}`}>{a.status}</Badge>
                  </td>
                  <td className="px-5 py-3.5 text-gray-600 tabular-nums">
                    {a.current_count} / {a.threshold}
                  </td>
                  <td className="px-5 py-3.5 text-gray-500 whitespace-nowrap">
                    {fmt(a.triggered_at)}
                  </td>
                  <td className="px-5 py-3.5">
                    <div className="flex gap-2">
                      {a.status === 'OPEN' && (
                        <button
                          onClick={() => action(a.id, 'ack')}
                          className="text-xs px-2.5 py-1 rounded-md bg-yellow-50 text-yellow-700 border border-yellow-200 hover:bg-yellow-100 font-medium"
                        >
                          Ack
                        </button>
                      )}
                      {a.status !== 'RESOLVED' && (
                        <button
                          onClick={() => action(a.id, 'resolve')}
                          className="text-xs px-2.5 py-1 rounded-md bg-green-50 text-green-700 border border-green-200 hover:bg-green-100 font-medium"
                        >
                          Resolve
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-gray-500">
            Page {page + 1} of {totalPages}
          </p>
          <div className="flex gap-2">
            <button
              disabled={page === 0}
              onClick={() => setPage(p => p - 1)}
              className="text-sm px-3 py-1.5 rounded-md border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
            >
              ← Previous
            </button>
            <button
              disabled={page >= totalPages - 1}
              onClick={() => setPage(p => p + 1)}
              className="text-sm px-3 py-1.5 rounded-md border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
            >
              Next →
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

function Badge({ cls, children }) {
  return (
    <span className={`inline-block text-xs font-semibold px-2 py-0.5 rounded-full ${cls}`}>
      {children}
    </span>
  );
}

function fmt(ts) {
  if (!ts) return '—';
  try { return new Date(ts).toLocaleString(); } catch { return ts; }
}
