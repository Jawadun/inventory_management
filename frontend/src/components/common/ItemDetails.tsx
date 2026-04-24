import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { itemService } from '../../services';
import { Item, ItemHistory, CreateIssueRequest } from '../../types';
import { LoadingSpinner } from '../common/Pagination';

interface ItemDetailsProps {
  isAdmin?: boolean;
}

export function ItemDetails({ isAdmin = false }: ItemDetailsProps) {
  const { id } = useParams();
  const navigate = useNavigate();

  const [item, setItem] = useState<Item | null>(null);
  const [history, setHistory] = useState<ItemHistory[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const [showIssueForm, setShowIssueForm] = useState(false);
  const [issueForm, setIssueForm] = useState<CreateIssueRequest>({
    item_id: id || '',
    recipient_id: '',
    quantity: 1,
    issue_type: 'personal',
  });
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!id) return;

    const fetchData = async () => {
      try {
        const [itemData, historyData] = await Promise.all([
          itemService.get(id),
          isAdmin ? itemService.getHistory(id) : Promise.resolve([]),
        ]);
        setItem(itemData);
        setHistory(historyData);
      } catch (err) {
        setError('Failed to load item');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id, isAdmin]);

  const handleIssue = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id) return;

    setSubmitting(true);
    try {
      // Use issue service to create issue
      const { issueService } = await import('../../services');
      await issueService.create(issueForm);
      setShowIssueForm(false);
      // Refresh data
      const [itemData, historyData] = await Promise.all([
        itemService.get(id),
        itemService.getHistory(id),
      ]);
      setItem(itemData);
      setHistory(historyData);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to issue item');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (error || !item) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600">{error || 'Item not found'}</p>
        <button
          onClick={() => navigate(-1)}
          className="mt-4 text-blue-600 hover:underline"
        >
          Go back
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-4xl">
      <button
        onClick={() => navigate(-1)}
        className="mb-4 text-sm text-gray-600 hover:text-gray-900 flex items-center gap-1"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
        </svg>
        Back
      </button>

      {error && (
        <div className="mb-4 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-200">
          <div className="flex items-start justify-between">
            <div className="flex items-start gap-4">
              {item.image_url && (
                <img
                  src={item.image_url}
                  alt={item.name}
                  className="w-24 h-24 rounded-lg object-cover"
                />
              )}
              <div>
                <h1 className="text-2xl font-bold text-gray-900">{item.name}</h1>
                {item.description && (
                  <p className="mt-1 text-gray-600">{item.description}</p>
                )}
                <div className="mt-2 flex items-center gap-3">
                  <StatusBadge status={item.status} />
                  {item.sku && (
                    <span className="text-sm text-gray-500">SKU: {item.sku}</span>
                  )}
                  {item.barcode && (
                    <span className="text-sm text-gray-500">Barcode: {item.barcode}</span>
                  )}
                </div>
              </div>
            </div>
            {isAdmin && (
              <button
                onClick={() => navigate(`/admin/items/${item.id}/edit`)}
                className="px-3 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
              >
                Edit
              </button>
            )}
          </div>
        </div>

        <div className="px-6 py-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div>
              <div className="text-sm text-gray-500">Quantity</div>
              <div className="text-lg font-semibold">
                <span className={item.quantity <= item.min_quantity ? 'text-red-600' : ''}>
                  {item.quantity}
                </span>{' '}
                <span className="text-sm font-normal text-gray-500">{item.unit}</span>
              </div>
              {item.quantity <= item.min_quantity && (
                <div className="text-sm text-red-600">
                  Minimum: {item.min_quantity}
                </div>
              )}
            </div>

            <div>
              <div className="text-sm text-gray-500">Category</div>
              <div className="font-medium">{item.category?.name || '-'}</div>
            </div>

            <div>
              <div className="text-sm text-gray-500">Supplier</div>
              <div className="font-medium">{item.supplier?.name || '-'}</div>
            </div>

            <div>
              <div className="text-sm text-gray-500">Location</div>
              <div className="font-medium">{item.location || '-'}</div>
            </div>

            <div>
              <div className="text-sm text-gray-500">Storage Location</div>
              <div className="font-medium">{item.storage_location || '-'}</div>
            </div>

            <div>
              <div className="text-sm text-gray-500">Condition</div>
              <div className="font-medium">{item.condition || '-'}</div>
            </div>

            <div>
              <div className="text-sm text-gray-500">Purchase Price</div>
              <div className="font-medium">
                {item.purchase_price
                  ? `$${item.purchase_price.toFixed(2)}`
                  : '-'}
              </div>
            </div>

            <div>
              <div className="text-sm text-gray-500">Warranty</div>
              <div className="font-medium">
                {item.warranty_months
                  ? `${item.warranty_months} months`
                  : '-'}
              </div>
            </div>
          </div>

          {item.notes && (
            <div className="mt-6 pt-6 border-t border-gray-200">
              <div className="text-sm text-gray-500">Notes</div>
              <div className="mt-1">{item.notes}</div>
            </div>
          )}

          <div className="mt-6 pt-6 border-t border-gray-200">
            <div className="text-sm text-gray-500">
              Created: {new Date(item.created_at).toLocaleString()}
            </div>
            <div className="text-sm text-gray-500">
              Updated: {new Date(item.updated_at).toLocaleString()}
            </div>
          </div>
        </div>

        {isAdmin && history.length > 0 && (
          <div className="border-t border-gray-200">
            <div className="px-6 py-4">
              <h2 className="text-lg font-medium mb-4">Item History</h2>
              <div className="space-y-3">
                {history.map((entry) => (
                  <div
                    key={entry.id}
                    className="flex items-start justify-between text-sm"
                  >
                    <div>
                      <span className="font-medium">
                        {entry.change_type}
                      </span>
                      <span className="text-gray-500 mx-2">
                        {entry.previous_quantity} → {entry.new_quantity}
                      </span>
                      {entry.reason && (
                        <span className="text-gray-500">
                          ({entry.reason})
                        </span>
                      )}
                    </div>
                    <div className="text-gray-500">
                      {new Date(entry.created_at).toLocaleString()}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>

      {showIssueForm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="fixed inset-0 bg-black bg-opacity-50" onClick={() => setShowIssueForm(false)} />
          <form
            onSubmit={handleIssue}
            className="relative bg-white rounded-lg shadow-xl p-6 max-w-md w-full mx-4 z-10"
          >
            <h3 className="text-lg font-medium mb-4">Issue Item</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Quantity
                </label>
                <input
                  type="number"
                  min="1"
                  max={item.quantity}
                  value={issueForm.quantity}
                  onChange={(e) =>
                    setIssueForm({
                      ...issueForm,
                      quantity: parseInt(e.target.value) || 1,
                    })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Issue Type
                </label>
                <select
                  value={issueForm.issue_type}
                  onChange={(e) =>
                    setIssueForm({ ...issueForm, issue_type: e.target.value as any })
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                >
                  <option value="personal">Personal Use</option>
                  <option value="classroom">Classroom</option>
                  <option value="lab">Lab</option>
                  <option value="teachers_room">Teacher's Room</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Notes
                </label>
                <textarea
                  value={issueForm.notes || ''}
                  onChange={(e) =>
                    setIssueForm({ ...issueForm, notes: e.target.value })
                  }
                  rows={2}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                />
              </div>
            </div>
            <div className="mt-6 flex justify-end gap-3">
              <button
                type="button"
                onClick={() => setShowIssueForm(false)}
                className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={submitting}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
              >
                {submitting ? 'Processing...' : 'Issue Item'}
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}

function StatusBadge({ status }: { status: string }) {
  const statusStyles: Record<string, string> = {
    available: 'bg-green-100 text-green-800',
    issued: 'bg-blue-100 text-blue-800',
    reserved: 'bg-yellow-100 text-yellow-800',
    damaged: 'bg-red-100 text-red-800',
    retired: 'bg-gray-100 text-gray-800',
  };

  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
      statusStyles[status] || 'bg-gray-100 text-gray-800'
    }`}>
      {status}
    </span>
  );
}