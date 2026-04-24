import { useState, useEffect } from 'react';
import { adminService } from '../../services';
import { Supplier } from '../../types';

export function AdminSuppliersPage() {
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState('');

  useEffect(() => {
    loadSuppliers();
  }, [page]);

  const loadSuppliers = async () => {
    setLoading(true);
    try {
      const res = await adminService.listSuppliers(search, { page, page_size: 20 });
      setSuppliers(res.suppliers);
      setTotalPages(res.total_pages);
    } finally {
      setLoading(false);
    }
  };

  const handleToggle = async (id: string, active: boolean) => {
    try {
      await adminService.toggleSupplier(id, active);
      loadSuppliers();
    } catch (err: any) {
      alert(err.response?.data?.message || 'Failed to update');
    }
  };

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Suppliers</h1>

      <form onSubmit={(e) => { e.preventDefault(); setPage(1); loadSuppliers(); }} className="mb-6 flex gap-4">
        <input
          type="text"
          placeholder="Search suppliers..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="flex-1 px-3 py-2 border rounded"
        />
        <button type="submit" className="px-4 py-2 bg-gray-100 rounded">Search</button>
      </form>

      {loading ? (
        <div>Loading...</div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Contact</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Phone</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Email</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {suppliers.map((supplier) => (
                <tr key={supplier.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 font-medium">{supplier.name}</td>
                  <td className="px-6 py-4">{supplier.contact_person}</td>
                  <td className="px-6 py-4">{supplier.phone}</td>
                  <td className="px-6 py-4">{supplier.email}</td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-1 rounded text-xs ${supplier.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                      {supplier.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <button
                      onClick={() => handleToggle(supplier.id, !supplier.is_active)}
                      className="text-blue-600 hover:underline text-sm"
                    >
                      {supplier.is_active ? 'Deactivate' : 'Activate'}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div className="flex justify-center gap-2 mt-4">
        <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1} className="px-3 py-1 border rounded disabled:opacity-50">Previous</button>
        <span className="px-3 py-1">Page {page} of {totalPages}</span>
        <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages} className="px-3 py-1 border rounded disabled:opacity-50">Next</button>
      </div>
    </div>
  );
}