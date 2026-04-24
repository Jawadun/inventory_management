import { useState, useEffect } from 'react';
import { adminService } from '../../services';

interface PendingUser {
  id: string;
  username: string;
  full_name: string;
  email?: string;
  department?: string;
  employee_id?: string;
  phone?: string;
  status: string;
  created_at: string;
}

export function AdminPendingUsersPage() {
  const [users, setUsers] = useState<PendingUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState('');

  useEffect(() => {
    loadUsers();
  }, [page]);

  const loadUsers = async () => {
    setLoading(true);
    try {
      const res = await adminService.getPendingUsers(search, { page, page_size: 20 });
      setUsers(res.users);
      setTotalPages(res.total_pages);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (userId: string) => {
    if (!confirm('Are you sure you want to approve this user registration?')) return;
    try {
      await adminService.approveUser(userId, 2);
      loadUsers();
      alert('User approved successfully');
    } catch (err: any) {
      alert(err.response?.data?.message || 'Failed to approve user');
    }
  };

  const handleReject = async (userId: string) => {
    if (!confirm('Are you sure you want to reject this user registration?')) return;
    try {
      await adminService.rejectUser(userId);
      loadUsers();
      alert('User rejected');
    } catch (err: any) {
      alert(err.response?.data?.message || 'Failed to reject user');
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    loadUsers();
  };

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Pending User Registrations</h1>

      <form onSubmit={handleSearch} className="mb-6 flex gap-4">
        <input
          type="text"
          placeholder="Search by username or name..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="flex-1 px-3 py-2 border rounded"
        />
        <button type="submit" className="px-4 py-2 bg-gray-100 rounded">Search</button>
      </form>

      {loading ? (
        <div>Loading...</div>
      ) : users.length === 0 ? (
        <div className="text-center py-8 text-gray-500">No pending registrations</div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Username</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Full Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Email</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Department</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Employee ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Registered</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {users.map((user) => (
                <tr key={user.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">{user.username}</td>
                  <td className="px-6 py-4">{user.full_name}</td>
                  <td className="px-6 py-4">{user.email || '-'}</td>
                  <td className="px-6 py-4">{user.department || '-'}</td>
                  <td className="px-6 py-4">{user.employee_id || '-'}</td>
                  <td className="px-6 py-4">{new Date(user.created_at).toLocaleDateString()}</td>
                  <td className="px-6 py-4 flex gap-2">
                    <button
                      onClick={() => handleApprove(user.id)}
                      className="px-3 py-1 bg-green-600 text-white rounded text-sm hover:bg-green-700"
                    >
                      Approve
                    </button>
                    <button
                      onClick={() => handleReject(user.id)}
                      className="px-3 py-1 bg-red-600 text-white rounded text-sm hover:bg-red-700"
                    >
                      Reject
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {totalPages > 1 && (
        <div className="flex justify-center gap-2 mt-4">
          <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1} className="px-3 py-1 border rounded disabled:opacity-50">
            Previous
          </button>
          <span className="px-3 py-1">Page {page} of {totalPages}</span>
          <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages} className="px-3 py-1 border rounded disabled:opacity-50">
            Next
          </button>
        </div>
      )}
    </div>
  );
}