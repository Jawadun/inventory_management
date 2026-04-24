import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { adminService } from '../../services';
import { AdminOverview, AnalyticsData } from '../../types';

export function AdminDashboardPage() {
  const [overview, setOverview] = useState<AdminOverview | null>(null);
  const [analytics, setAnalytics] = useState<AnalyticsData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      adminService.getOverview(),
      adminService.getAnalytics(),
    ]).then(([ov, an]) => {
      setOverview(ov);
      setAnalytics(an);
      setLoading(false);
    });
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Admin Dashboard</h1>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <StatCard
          label="Total Users"
          value={overview?.total_users || 0}
          link="/admin/users"
        />
        <StatCard
          label="Active Users"
          value={overview?.active_users || 0}
          link="/admin/users"
        />
        <StatCard
          label="Total Items"
          value={overview?.total_items || 0}
          link="/admin/items"
        />
        <StatCard
          label="Available Qty"
          value={overview?.available_quantity || 0}
          link="/admin/items"
          color="green"
        />
        <StatCard
          label="Issued Qty"
          value={overview?.issued_quantity || 0}
          link="/admin/issues"
          color="yellow"
        />
        <StatCard
          label="Suppliers"
          value={overview?.total_suppliers || 0}
          link="/admin/suppliers"
        />
        <StatCard
          label="Categories"
          value={overview?.total_categories || 0}
        />
        <StatCard
          label="Low Stock"
          value={overview?.low_stock_items || 0}
          link="/admin/items?low_stock=true"
          color="red"
        />
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <StatCard
          label="Pending Requests"
          value={overview?.pending_requests || 0}
          link="/admin/requests?status=pending"
          color="yellow"
        />
        <StatCard
          label="Pending Issues"
          value={overview?.pending_issues || 0}
          link="/admin/issues?status=pending"
          color="yellow"
        />
        <StatCard
          label="Overdue"
          value={overview?.overdue_items || 0}
          link="/admin/issues?overdue=true"
          color="red"
        />
        <StatCard
          label="Active Notices"
          value={overview?.active_notices || 0}
          link="/admin/notices"
          color="blue"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <QuickLinkCard
          title="Manage Users"
          description="Add, edit, deactivate users"
          link="/admin/users"
          icon="users"
        />
        <QuickLinkCard
          title="Manage Items"
          description="Add, edit inventory"
          link="/admin/items"
          icon="items"
        />
        <QuickLinkCard
          title="Requests"
          description="Review pending requests"
          link="/admin/requests"
          icon="requests"
        />
        <QuickLinkCard
          title="Issues"
          description="Track issued items"
          link="/admin/issues"
          icon="issues"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-lg font-semibold">Requests by Status</h2>
            <Link to="/admin/requests" className="text-sm text-blue-600 hover:underline">
              View all
            </Link>
          </div>
          <div className="space-y-2">
            {analytics?.requests_by_status?.map((item) => (
              <div key={item.status} className="flex justify-between items-center">
                <span className="text-gray-600 capitalize">{item.status}</span>
                <span className="font-medium">{item.count}</span>
              </div>
            ))}
            {(!analytics?.requests_by_status || analytics.requests_by_status.length === 0) && (
              <p className="text-gray-500 text-sm">No requests</p>
            )}
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-lg font-semibold">Issues by Status</h2>
            <Link to="/admin/issues" className="text-sm text-blue-600 hover:underline">
              View all
            </Link>
          </div>
          <div className="space-y-2">
            {analytics?.issues_by_status?.map((item) => (
              <div key={item.status} className="flex justify-between items-center">
                <span className="text-gray-600 capitalize">{item.status}</span>
                <span className="font-medium">{item.count}</span>
              </div>
            ))}
            {(!analytics?.issues_by_status || analytics.issues_by_status.length === 0) && (
              <p className="text-gray-500 text-sm">No issues</p>
            )}
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-lg font-semibold">Items by Status</h2>
            <Link to="/admin/items" className="text-sm text-blue-600 hover:underline">
              View all
            </Link>
          </div>
          <div className="space-y-2">
            {analytics?.items_by_status?.map((item) => (
              <div key={item.status} className="flex justify-between items-center">
                <span className="text-gray-600 capitalize">{item.status}</span>
                <span className="font-medium">{item.count}</span>
              </div>
            ))}
            {(!analytics?.items_by_status || analytics.items_by_status.length === 0) && (
              <p className="text-gray-500 text-sm">No items</p>
            )}
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-lg font-semibold">Top Categories</h2>
          </div>
          <div className="space-y-2">
            {analytics?.top_categories?.slice(0, 5).map((cat) => (
              <div key={cat.category_id} className="flex justify-between items-center">
                <span className="text-gray-600 truncate max-w-[150px]">
                  {cat.category_name}
                </span>
                <span className="font-medium">{cat.total_items} items</span>
              </div>
            ))}
            {(!analytics?.top_categories || analytics.top_categories.length === 0) && (
              <p className="text-gray-500 text-sm">No categories</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

interface StatCardProps {
  label: string;
  value: number;
  link?: string;
  color?: 'blue' | 'green' | 'yellow' | 'red';
}

function StatCard({ label, value, link, color = 'blue' }: StatCardProps) {
  const colors: Record<string, string> = {
    blue: 'bg-blue-50 border-blue-200 text-blue-700',
    green: 'bg-green-50 border-green-200 text-green-700',
    yellow: 'bg-yellow-50 border-yellow-200 text-yellow-700',
    red: 'bg-red-50 border-red-200 text-red-700',
  };

  const content = (
    <div className={`p-4 rounded-lg border ${colors[color]}`}>
      <div className="text-2xl font-bold">{value}</div>
      <div className="text-sm opacity-80">{label}</div>
    </div>
  );

  if (link) {
    return <Link to={link}>{content}</Link>;
  }
  return content;
}

interface QuickLinkCardProps {
  title: string;
  description: string;
  link: string;
  icon: string;
}

function QuickLinkCard({ title, description, link, icon }: QuickLinkCardProps) {
  const icons: Record<string, JSX.Element> = {
    users: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
      </svg>
    ),
    items: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
      </svg>
    ),
    requests: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
      </svg>
    ),
    issues: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0v6m0-6L4 7m16 10l-4-4m4 4l4-4" />
      </svg>
    ),
  };

  return (
    <Link
      to={link}
      className="block bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div className="text-blue-600 mb-3">{icons[icon]}</div>
      <h3 className="font-semibold text-gray-900">{title}</h3>
      <p className="text-sm text-gray-500 mt-1">{description}</p>
    </Link>
  );
}