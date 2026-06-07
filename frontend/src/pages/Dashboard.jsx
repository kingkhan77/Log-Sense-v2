import { useState, useEffect } from 'react';
import api from '../api';

const SEV_STYLE = {
  CRITICAL: { badge: 'bg-red-100 text-red-700',    bar: 'bg-red-500'    },
  WARNING:  { badge: 'bg-yellow-100 text-yellow-700', bar: 'bg-yellow-500' },
  INFO:     { badge: 'bg-blue-100 text-blue-700',   bar: 'bg-blue-500'   },
};

function StatCard({ label, value, sub, accent }) {
  const accents = {
    red:    'border-l-red-500',
    yellow: 'border-l-yellow-500',
    blue:   'border-l-blue-500',
    green:  'border-l-green-500',
  };
  return (
    <div className={`bg-white rounded-xl border border-l-4 ${accents[accent] ?? 'border-l-gray-300'} p-6`}>
      <p className="text-xs font-medium text-gray-500 uppercase tracking-wide">{label}</p>
      <p className="text-3xl font-bold text-gray-900 mt-1">{value ?? '—'}</p>
      {sub && <p className="text-xs text-gray-400 mt-1">{sub}</p>}
    </div>
  );
}

export default function Dashboard() {
  const [summary, setSummary] = useState(null);
  const [error,   setError]   = useState('');

  useEffect(() => {
    api.get('/dashboard/summary')
      .then(r => setSummary(r.data))
      .catch(() => setError('Failed to load dashboard data.'));
  }, []);

  if (error)   return <p className="text-red-500 text-sm">{error}</p>;
  if (!summary) return <p className="text-gray-400 text-sm animate-pulse">Loading…</p>;

  const sevEntries = Object.entries(summary.by_severity ?? {});
  const total = sevEntries.reduce((s, [, n]) => s + n, 0) || 1;

  return (
    <div className="space-y-8 max-w-5xl">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-sm text-gray-500 mt-1">Live overview of your observability platform</p>
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard label="Open Alerts"   value={summary.open_alerts}     accent="red"    sub="currently active" />
        <StatCard label="Alerts (24 h)" value={summary.alerts_last_24h} accent="yellow" sub="last 24 hours" />
        <StatCard label="Enabled Rules" value={summary.enabled_rules}   accent="blue"   sub="actively evaluated" />
        <StatCard label="Services"      value={summary.services}        accent="green"  sub="registered" />
      </div>

      {/* Severity breakdown */}
      <div className="bg-white rounded-xl border p-6">
        <h2 className="text-sm font-semibold text-gray-700 mb-5">Open Alerts by Severity</h2>
        {sevEntries.length === 0 ? (
          <p className="text-sm text-gray-400">No open alerts — all clear.</p>
        ) : (
          <div className="space-y-4">
            {sevEntries.map(([sev, count]) => {
              const s = SEV_STYLE[sev] ?? SEV_STYLE.INFO;
              return (
                <div key={sev} className="flex items-center gap-4">
                  <span className={`text-xs font-semibold px-2.5 py-1 rounded-full w-24 text-center ${s.badge}`}>
                    {sev}
                  </span>
                  <div className="flex-1 bg-gray-100 rounded-full h-2 overflow-hidden">
                    <div
                      className={`h-2 rounded-full transition-all ${s.bar}`}
                      style={{ width: `${Math.round((count / total) * 100)}%` }}
                    />
                  </div>
                  <span className="text-sm font-medium text-gray-700 w-6 text-right">{count}</span>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
