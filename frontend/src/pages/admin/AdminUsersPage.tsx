import { useState, useEffect } from 'react';
import { adminService } from '../../services';
import { User } from '../../types';

export function AdminUsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState('');

  useEffect(() => {
    loadUsers();
  }, [page, filter]);

  const loadUsers = async () => {
    setLoading(true);
    try {
      const res = await adminService.listUsers({ search, status: filter }, { page, page_size: 20 });
      setUsers(res.users);
      setTotalPages(res.total_pages);
    } finally {
      setLoading(false);
    }
  };

  const handleToggle = async (userId: string, active: boolean) => {
    if (!confirm(`Are you sure you want to ${active ? 'activate' : 'deactivate'} this user?`)) return;
    try {
      await adminService.toggleUser(userId, active);
      loadUsers();
    } catch (err: any) {
      alert(err.response?.data?.message || 'Failed to update user');
    }
  };

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Users</h1>

      <form onSubmit={(e) => { e.preventDefault(); setPage(1); loadUsers(); }} className="mb-6 flex gap-4">
        <input
          type="text"
          placeholder="Search users..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="flex-1 px-3 py-2 border rounded"
        />
        <select
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          className="px-3 py-2 border rounded"
        >
          <option value="">All</option>
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
        </select>
        <button type="submit" className="px-4 py-2 bg-gray-100 rounded">Search</button>
      </form>

      {loading ? (
        <div>Loading...</div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Username</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Department</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Role</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {users.map((user) => (
                <tr key={user.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">{user.username}</td>
                  <td className="px-6 py-4">{user.full_name}</td>
                  <td className="px-6 py-4">{user.department}</td>
                  <td className="px-6 py-4">{user.role_id === 1 ? 'Admin' : user.role_id === 2 ? 'User' : 'Viewer'}</td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-1 rounded text-xs ${user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                      {user.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <button
                      onClick={() => handleToggle(user.id, !user.is_active)}
                      className="text-blue-600 hover:underline text-sm"
                    >
                      {user.is_active ? 'Deactivate' : 'Activate'}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div className="flex justify-center gap-2 mt-4">
        <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1} className="px-3 py-1 border rounded disabled:opacity-50">
          Previous
        </button>
        <span className="px-3 py-1">Page {page} of {totalPages}</span>
        <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages} className="px-3 py-1 border rounded disabled:opacity-50">
          Next
        </button>
      </div>
    </div>
  );
}