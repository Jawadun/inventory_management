import { useState, useEffect, useCallback } from 'react';
import { adminService } from '../../services';
import { ItemRequest, AdminFilter, PaginationParams } from '../../types';
import { Pagination, LoadingSpinner, EmptyState, ConfirmDialog } from '../common/Pagination';
import { SelectInput } from '../common/SelectInput';

const STATUS_OPTIONS = [
  { value: '', label: 'All Status' },
  { value: 'pending', label: 'Pending' },
  { value: 'approved', label: 'Approved' },
  { value: 'rejected', label: 'Rejected' },
  { value: 'fulfilled', label: 'Fulfilled' },
];

export function AdminRequestManagement() {
  const [requests, setRequests] = useState<ItemRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [status, setStatus] = useState('');

  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [action, setAction] = useState<'approve' | 'reject' | 'fulfill' | null>(null);
  const [rejectReason, setRejectReason] = useState('');
  const [processing, setProcessing] = useState(false);

  const fetchRequests = useCallback(async () => {
    setLoading(true);
    try {
      const filter: AdminFilter = {};
      if (status) filter.status = status;
      const pageParams: PaginationParams = { page, page_size: 20 };
      const res = await adminService.listRequests(filter, pageParams);
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

  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      setSelectedIds(
        requests
          .filter((r) => r.status === 'pending')
          .map((r) => r.id)
      );
    } else {
      setSelectedIds([]);
    }
  };

  const handleSelect = (id: string, checked: boolean) => {
    if (checked) {
      setSelectedIds([...selectedIds, id]);
    } else {
      setSelectedIds(selectedIds.filter((i) => i !== id));
    }
  };

  const handleAction = async () => {
    if (!action || selectedIds.length === 0) return;
    if (action === 'reject' && !rejectReason) {
      alert('Please provide a rejection reason');
      return;
    }

    setProcessing(true);
    try {
      for (const id of selectedIds) {
        if (action === 'approve') {
          await adminService.manageRequest(id, 'approved');
        } else if (action === 'reject') {
          await adminService.manageRequest(id, 'rejected', rejectReason);
        } else if (action === 'fulfill') {
          await adminService.manageRequest(id, 'fulfilled');
        }
      }
      setSelectedIds([]);
      setAction(null);
      setRejectReason('');
      fetchRequests();
    } catch (err: any) {
      alert(err.response?.data?.message || 'Action failed');
    } finally {
      setProcessing(false);
    }
  };

  const pendingCount = requests.filter((r) => r.status === 'pending').length;

  return (
    <div>
      <div className="mb-6">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Requests</h1>
            <p className="text-sm text-gray-500 mt-1">
              {totalCount} total requests
              {pendingCount > 0 && (
                <span className="ml-2 text-yellow-600">
                  ({pendingCount} pending)
                </span>
              )}
            </p>
          </div>
          {selectedIds.length > 0 && (
            <div className="flex gap-2">
              <button
                onClick={() => setAction('approve')}
                disabled={pendingCount === 0}
                className="px-3 py-1 bg-green-600 text-white rounded text-sm disabled:opacity-50"
              >
                Approve ({selectedIds.length})
              </button>
              <button
                onClick={() => setAction('reject')}
                disabled={pendingCount === 0}
                className="px-3 py-1 bg-red-600 text-white rounded text-sm disabled:opacity-50"
              >
                Reject ({selectedIds.length})
              </button>
              <button
                onClick={() => setAction('fulfill')}
                disabled={pendingCount === 0}
                className="px-3 py-1 bg-blue-600 text-white rounded text-sm disabled:opacity-50"
              >
                Fulfill ({selectedIds.length})
              </button>
            </div>
          )}
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
          <EmptyState title="No requests" description="All caught up!" />
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3">
                    <input
                      type="checkbox"
                      checked={
                        selectedIds.length > 0 &&
                        selectedIds.length ===
                          requests.filter((r) => r.status === 'pending').length
                      }
                      onChange={(e) => handleSelectAll(e.target.checked)}
                      className="rounded border-gray-300"
                    />
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Item
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Requested By
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Qty
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Type
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Reason
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Date
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {requests.map((req) => (
                  <tr key={req.id} className="hover:bg-gray-50">
                    <td className="px-4 py-4">
                      {req.status === 'pending' && (
                        <input
                          type="checkbox"
                          checked={selectedIds.includes(req.id)}
                          onChange={(e) =>
                            handleSelect(req.id, e.target.checked)
                          }
                          className="rounded border-gray-300"
                        />
                      )}
                    </td>
                    <td className="px-6 py-4">
                      <div className="font-medium text-gray-900">
                        {req.item?.name || 'Unknown'}
                      </div>
                    </td>
                    <td className="px-6 py-4 text-gray-600">
                      {req.user?.full_name || req.user?.username || '-'}
                    </td>
                    <td className="px-6 py-4 text-gray-900">{req.quantity}</td>
                    <td className="px-6 py-4 text-gray-600 capitalize">
                      {req.request_type.replace('_', ' ')}
                    </td>
                    <td className="px-6 py-4 text-gray-500 text-sm max-w-xs truncate">
                      {req.reason || '-'}
                    </td>
                    <td className="px-6 py-4">
                      <StatusBadge status={req.status} />
                    </td>
                    <td className="px-6 py-4 text-gray-500 text-sm">
                      {new Date(req.requested_at).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4">
                      {req.status === 'pending' && (
                        <div className="flex gap-2">
                          <button
                            onClick={async () => {
                              try {
                                await adminService.manageRequest(
                                  req.id,
                                  'approved'
                                );
                                fetchRequests();
                              } catch (err: any) {
                                alert(
                                  err.response?.data?.message ||
                                    'Failed to approve'
                                );
                              }
                            }}
                            className="text-green-600 hover:text-green-700 text-sm"
                          >
                            Approve
                          </button>
                          <button
                            onClick={async () => {
                              const reason = prompt(
                                'Enter rejection reason:'
                              );
                              if (!reason) return;
                              try {
                                await adminService.manageRequest(
                                  req.id,
                                  'rejected',
                                  reason
                                );
                                fetchRequests();
                              } catch (err: any) {
                                alert(
                                  err.response?.data?.message ||
                                    'Failed to reject'
                                );
                              }
                            }}
                            className="text-red-600 hover:text-red-700 text-sm"
                          >
                            Reject
                          </button>
                        </div>
                      )}
                      {req.status === 'approved' && (
                        <button
                          onClick={async () => {
                            try {
                              await adminService.manageRequest(
                                req.id,
                                'fulfilled'
                              );
                              fetchRequests();
                            } catch (err: any) {
                              alert(
                                err.response?.data?.message ||
                                  'Failed to fulfill'
                              );
                            }
                          }}
                          className="text-blue-600 hover:text-blue-700 text-sm"
                        >
                          Fulfill
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
        open={action === 'reject'}
        title="Reject Requests"
        message="Please provide a reason for rejection:"
        confirmLabel="Reject"
        variant="danger"
        onConfirm={handleAction}
        onCancel={() => {
          setAction(null);
          setRejectReason('');
        }}
      >
        <input
          type="text"
          value={rejectReason}
          onChange={(e) => setRejectReason(e.target.value)}
          placeholder="Rejection reason..."
          className="w-full px-3 py-2 border border-gray-300 rounded mt-3"
        />
      </ConfirmDialog>
    </div>
  );
}

function StatusBadge({ status }: { status: string }) {
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