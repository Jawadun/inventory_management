import React, { useState, useEffect } from 'react';
import api from '../../services/api';
import { PublicStats } from '../../types';

export function StatsPage() {
  const [stats, setStats] = useState<PublicStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  console.log('StatsPage rendering', { loading, error, stats });

  useEffect(() => {
    console.log('StatsPage useEffect running');
    api.get<PublicStats>('/public/stats')
      .then((res) => {
        console.log('Stats API response:', res.data);
        setStats(res.data);
        setError(null);
      })
      .catch((err) => {
        console.error('Failed to load stats:', err);
        setError('Unable to connect to server. Please try again later.');
      })
      .finally(() => {
        console.log('StatsPage loading complete');
        setLoading(false);
      });
  }, []);

  console.log('StatsPage render state:', { loading, error, stats });

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold mb-6">Inventory Statistics</h1>
      
      {error && (
        <div className="bg-yellow-50 border-l-4 border-yellow-400 p-4 mb-6">
          <p className="text-sm text-yellow-700">{error}</p>
        </div>
      )}
      
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <StatCard label="Total Items" value={stats?.total_items ?? 0} />
        <StatCard label="Available" value={stats?.available_items ?? 0} />
        <StatCard label="Issued" value={stats?.issued_items ?? 0} />
        <StatCard label="Categories" value={stats?.total_categories ?? 0} />
        <StatCard label="Suppliers" value={stats?.total_suppliers ?? 0} />
        <StatCard label="Pending Requests" value={stats?.pending_requests ?? 0} />
        <StatCard label="Active Notices" value={stats?.active_notices ?? 0} />
        <StatCard label="Low Stock" value={stats?.low_stock_items ?? 0} />
      </div>
    </div>
  );
}

function StatCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <div className="text-2xl font-bold text-blue-600">{value}</div>
      <div className="text-sm text-gray-500">{label}</div>
    </div>
  );
}