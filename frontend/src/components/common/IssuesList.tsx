import { useState, useEffect, useCallback } from 'react';
import { issueService } from '../../services';
import { IssueRecord, IssueStatus } from '../../types';
import { Pagination, LoadingSpinner, EmptyState, ConfirmDialog } from '../common/Pagination';
import { SelectInput } from '../common/SelectInput';

const STATUS_OPTIONS = [
  { value: '', label: 'All Status' },
  { value: 'pending', label: 'Pending' },
  { value: 'approved', label: 'Approved' },
  { value: 'issued', label: 'Issued' },
  { value: 'returned', label: 'Returned' },
  { value: 'overdue', label: 'Overdue' },
];

interface IssuesListProps {
  isAdmin?: boolean;
}

export function IssuesList({ isAdmin = false }: IssuesListProps) {
  const [issues, setIssues] = useState<IssueRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [status, setStatus] = useState('');
  const [showReturn, setShowReturn] = useState<string | null>(null);
  const [returnCondition, setReturnCondition] = useState('');
  const [returnRemarks, setReturnRemarks] = useState('');
  const [processing, setProcessing] = useState(false);

  const fetchIssues = useCallback(async () => {
    setLoading(true);
    try {
      const filter: any = {};
      if (status) filter.status = status;
      const res = isAdmin
        ? await issueService.list(filter, page, 20)
        : await issueService.getMyIssues(page, 20);
      setIssues(res.issues);
      setTotalPages(res.total_pages);
      setTotalCount(res.total_count);
    } catch (err) {
      console.error('Failed to fetch issues:', err);
    } finally {
      setLoading(false);
    }
  }, [page, status, isAdmin]);

  useEffect(() => {
    fetchIssues();
  }, [fetchIssues]);

  const handleReturn = async () => {
    if (!showReturn) return;
    setProcessing(true);
    try {
      await issueService.return(showReturn, {
        return_condition: returnCondition,
        return_remarks: returnRemarks,
      });
      setShowReturn(null);
      setReturnCondition('');
      setReturnRemarks('');
      fetchIssues();
    } catch (err: any) {
      alert(err.response?.data?.message || 'Failed to process return');
    } finally {
      setProcessing(false);
    }
  };

  const isOverdue = (issue: IssueRecord) =>
    issue.status === 'issued' &&
    issue.due_date &&
    new Date(issue.due_date) < new Date();

  const activeCount = issues.filter((i) => i.status === 'issued').length;
  const overdueCount = issues.filter(isOverdue).length;

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">
          {isAdmin ? 'Issues' : 'My Issues'}
        </h1>
        <p className="text-sm text-gray-500 mt-1">
          {totalCount} {totalCount === 1 ? 'record' : 'records'}
          {activeCount > 0 && (
            <span className="ml-2 text-blue-600">
              ({activeCount} active)
            </span>
          )}
          {overdueCount > 0 && (
            <span className="ml-2 text-red-600">
              ({overdueCount} overdue)
            </span>
          )}
        </p>
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
        ) : issues.length === 0 ? (
          <EmptyState
            title="No issues"
            description={status ? 'Try changing your filter' : 'No items have been issued yet'}
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Item
                  </th>
                  {isAdmin && (
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Recipient
                    </th>
                  )}
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Qty
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Type
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Issued
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Due Date
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    Status
                  </th>
                  {isAdmin && (
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Actions
                    </th>
                  )}
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {issues.map((issue) => (
                  <tr key={issue.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4">
                      <div className="font-medium text-gray-900">
                        {issue.item?.name || 'Unknown Item'}
                      </div>
                      {issue.notes && (
                        <div className="text-sm text-gray-500 truncate max-w-xs">
                          {issue.notes}
                        </div>
                      )}
                    </td>
                    {isAdmin && (
                      <td className="px-6 py-4 text-gray-600">
                        {issue.recipient?.full_name ||
                          issue.recipient?.username ||
                          '-'}
                      </td>
                    )}
                    <td className="px-6 py-4 text-gray-900">{issue.quantity}</td>
                    <td className="px-6 py-4 text-gray-600 capitalize">
                      {issue.issue_type.replace('_', ' ')}
                    </td>
                    <td className="px-6 py-4 text-gray-500 text-sm">
                      {new Date(issue.issue_date).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 text-gray-500 text-sm">
                      {issue.due_date
                        ? new Date(issue.due_date).toLocaleDateString()
                        : '-'}
                    </td>
                    <td className="px-6 py-4">
                      <StatusBadge status={issue.status} isOverdue={isOverdue(issue)} />
                    </td>
                    {isAdmin && (
                      <td className="px-6 py-4">
                        {issue.status === 'issued' && (
                          <button
                            onClick={() => setShowReturn(issue.id)}
                            className="text-blue-600 hover:text-blue-700 text-sm"
                          >
                            Mark Returned
                          </button>
                        )}
                      </td>
                    )}
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
        open={!!showReturn}
        title="Process Return"
        message="Enter return details:"
        confirmLabel="Confirm Return"
        onConfirm={handleReturn}
        onCancel={() => {
          setShowReturn(null);
          setReturnCondition('');
          setReturnRemarks('');
        }}
      >
        <div className="space-y-4 mt-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Condition
            </label>
            <select
              value={returnCondition}
              onChange={(e) => setReturnCondition(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded"
            >
              <option value="">Select condition</option>
              <option value="good">Good</option>
              <option value="damaged">Damaged</option>
              <option value="needs_repair">Needs Repair</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Remarks
            </label>
            <textarea
              value={returnRemarks}
              onChange={(e) => setReturnRemarks(e.target.value)}
              rows={2}
              className="w-full px-3 py-2 border border-gray-300 rounded"
              placeholder="Any notes about the return..."
            />
          </div>
        </div>
      </ConfirmDialog>
    </div>
  );
}

function StatusBadge({ status, isOverdue }: { status: IssueStatus; isOverdue?: boolean }) {
  const styles: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-blue-100 text-blue-800',
    issued: 'bg-green-100 text-green-800',
    returned: 'bg-gray-100 text-gray-800',
  };

  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
        isOverdue ? 'bg-red-100 text-red-800' : styles[status] || styles.pending
      }`}
    >
      {isOverdue ? 'overdue' : status}
    </span>
  );
}