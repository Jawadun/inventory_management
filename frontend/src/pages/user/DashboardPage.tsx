import { useState, useEffect } from 'react';
import { api } from '../../services';

interface Stats {
  pending_requests: number;
  approved_requests: number;
  issued_items: number;
  returned_items: number;
}

export function DashboardPage() {
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get('/users/me').then((res) => {
      setStats(res.data);
      setLoading(false);
    });
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Dashboard</h1>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <StatCard label="Pending Requests" value={stats?.pending_requests || 0} color="yellow" />
        <StatCard label="Approved" value={stats?.approved_requests || 0} color="green" />
        <StatCard label="Issued Items" value={stats?.issued_items || 0} color="blue" />
        <StatCard label="Returned" value={stats?.returned_items || 0} color="gray" />
      </div>
      <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-4">
        <QuickAction
          title="Request Item"
          description="Submit a new item request"
          link="/requests/new"
          icon="📝"
        />
        <QuickAction
          title="View Inventory"
          description="Browse available items"
          link="/inventory"
          icon="📦"
        />
        <QuickAction
          title="My Requests"
          description="View request history"
          link="/requests"
          icon="📋"
        />
      </div>
    </div>
  );
}

function StatCard({ label, value, color }: { label: string; value: number; color: string }) {
  const colors: Record<string, string> = {
    yellow: 'bg-yellow-50 text-yellow-700 border-yellow-200',
    green: 'bg-green-50 text-green-700 border-green-200',
    blue: 'bg-blue-50 text-blue-700 border-blue-200',
    gray: 'bg-gray-50 text-gray-700 border-gray-200',
  };
  return (
    <div className={`p-6 rounded-lg border ${colors[color]}`}>
      <div className="text-3xl font-bold">{value}</div>
      <div className="text-sm">{label}</div>
    </div>
  );
}

function QuickAction({ title, description, link, icon }: { title: string; description: string; link: string; icon: string }) {
  return (
    <a
      href={link}
      className="block bg-white p-6 rounded-lg shadow hover:shadow-md transition"
    >
      <div className="text-3xl mb-2">{icon}</div>
      <h3 className="font-semibold">{title}</h3>
      <p className="text-sm text-gray-500">{description}</p>
    </a>
  );
}