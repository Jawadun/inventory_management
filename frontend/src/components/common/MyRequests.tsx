import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { requestService } from '../../services';
import { ItemRequest, RequestStatus } from '../../types';
import { Pagination, LoadingSpinner, EmptyState, ConfirmDialog } from '../common/Pagination';
import { SelectInput } from '../common/SelectInput';

const STATUS_OPTIONS = [
  { value: '', label: 'All Status' },
  { value: 'pending', label: 'Pending' },
  { value: 'approved', label: 'Approved' },
  { value: 'rejected', label: 'Rejected' },
  { value: 'fulfilled', label: 'Fulfilled' },
  { value: 'cancelled', label: 'Cancelled' },
];

export function MyRequests() {
  const [requests, setRequests] = useState<ItemRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [status, setStatus] = useState('');

  const [cancelId, setCancelId] = useState<string | null>(null);
  const [cancelling, setCancelling] = useState(false);

  const fetchRequests = useCallback(async () => {
    setLoading(true);
    try {
      const filter: any = {};
      if (status) filter.status = status;
      const res = await requestService.getMyRequests(page, 20);
      setRequests(res.requests);
      setTotalPages(res.total_pages);
      setTotalCount(res.total_count);
    } catch (err) {
      console.error('Failed to fetch requests:', err);
    } finally {
      setLoading(false);
    }
  }, [page, status]);

  useEffect(() => {
    fetchRequests();
  }, [fetchRequests]);

  const handleCancel = async () => {
    if (!cancelId) return;
    setCancelling(true);
    try {
      await requestService.cancel(cancelId);
      setCancelId(null);
      fetchRequests();
    } catch (err: any) {
      alert(err.response?.data?.message || 'Failed to cancel request');
    } finally {
      setCancelling(false);
    }
  };

  return (
    <div>
      <div className="mb-6">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">My Requests</h1>
            <p className="text-sm text-gray-500 mt-1">
              {totalCount} {totalCount === 1 ? 'request' : 'requests'}
            </p>
          </div>
          <Link
            to="/requests/new"
            className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
          >
            New Request
          </Link>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow mb-6">
        <div className="p-4 border-b border-gray-200">
          <div className="flex gap-4">
            <div className="w-48">
              <SelectInput
                value={status}
                onChange={setStatus}
                options={STATUS_OPTIONS}
                placeholder="Filter by status"
              />
            </div>
          </div>
        </div>

        {loading ? (
          <div className="py-12">
            <LoadingSpinner size="lg" />
          </div>
        ) : requests.length === 0 ? (
          <EmptyState
            title="No requests"
            description={status ? 'Try changing your filter' : 'Submit your first request'}
            action={
              <Link
                to="/requests/new"
                className="mt-4 inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
              >
                New Request
              </Link>
            }
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Item
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Qty
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Type
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Requested
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {requests.map((req) => (
                  <tr key={req.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4">
                      <div className="font-medium text-gray-900">
                        {req.item?.name || 'Unknown Item'}
                      </div>
                      {req.reason && (
                        <div className="text-sm text-gray-500 truncate max-w-xs">
                          {req.reason}
                        </div>
                      )}
                    </td>
                    <td className="px-6 py-4 text-gray-900">
                      {req.quantity}
                    </td>
                    <td className="px-6 py-4 text-gray-600 capitalize">
                      {req.request_type.replace('_', ' ')}
                    </td>
                    <td className="px-6 py-4">
                      <StatusBadge status={req.status} />
                    </td>
                    <td className="px-6 py-4 text-gray-500 text-sm">
                      {new Date(req.requested_at).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4">
                      {req.status === 'pending' && (
                        <button
                          onClick={() => setCancelId(req.id)}
                          className="text-red-600 hover:text-red-700 text-sm"
                        >
                          Cancel
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        <Pagination
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
        />
      </div>

      <ConfirmDialog
        open={!!cancelId}
        title="Cancel Request"
        message="Are you sure you want to cancel this request?"
        confirmLabel="Cancel Request"
        variant="danger"
        onConfirm={handleCancel}
        onCancel={() => setCancelId(null)}
      />
    </div>
  );
}

function StatusBadge({ status }: { status: RequestStatus }) {
  const styles: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
    fulfilled: 'bg-blue-100 text-blue-800',
    cancelled: 'bg-gray-100 text-gray-800',
  };

  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
        styles[status] || styles.pending
      }`}
    >
      {status}
    </span>
  );
}