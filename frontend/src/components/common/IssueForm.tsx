import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { itemService, issueService } from '../../services';
import { Item, User, CreateIssueRequest } from '../../types';
import { LoadingSpinner } from '../common/Pagination';

export function IssueForm() {
  const { requestId } = useParams();
  const navigate = useNavigate();

  const [formData, setFormData] = useState<CreateIssueRequest>({
    item_id: '',
    recipient_id: '',
    quantity: 1,
    issue_type: 'personal',
    due_date: undefined,
    notes: '',
    auto_approve: false,
  });

  const [items, setItems] = useState<Item[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [itemsRes, usersRes] = await Promise.all([
          itemService.list({ status: 'available' }, 1, 100),
          // Load users - in a real app, this would be an admin service
          Promise.resolve([]),
        ]);
        setItems(itemsRes.items);
        setUsers(usersRes);
      } catch (err) {
        setError('Failed to load data');
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const selectedItem = items.find((item) => item.id === formData.item_id);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.item_id || !formData.recipient_id) {
      setError('Please fill all required fields');
      return;
    }
    if (formData.quantity > (selectedItem?.quantity || 0)) {
      setError('Quantity exceeds available stock');
      return;
    }

    setError('');
    setSubmitting(true);
    try {
      await issueService.create(formData);
      navigate(-1);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create issue');
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

  return (
    <div className="max-w-2xl">
      {error && (
        <div className="mb-4 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-lg font-medium">Issue Item</h2>
          <p className="text-sm text-gray-500 mt-1">
            Record the issuance of an item to a user
          </p>
        </div>

        <div className="p-6 space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Item *
            </label>
            <select
              value={formData.item_id}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  item_id: e.target.value,
                  quantity: Math.min(formData.quantity, selectedItem?.quantity || 1),
                })
              }
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Select an item</option>
              {items.map((item) => (
                <option key={item.id} value={item.id} disabled={item.quantity <= 0}>
                  {item.name} ({item.quantity} {item.unit} available)
                </option>
              ))}
            </select>
          </div>

          {selectedItem && (
            <div className="bg-gray-50 rounded-lg p-4 text-sm">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <span className="text-gray-500">Available:</span>
                  <span className="ml-2 font-medium">
                    {selectedItem.quantity} {selectedItem.unit}
                  </span>
                </div>
                <div>
                  <span className="text-gray-500">Location:</span>
                  <span className="ml-2 font-medium">
                    {selectedItem.location || '-'}
                  </span>
                </div>
              </div>
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Recipient *
            </label>
            <select
              value={formData.recipient_id}
              onChange={(e) =>
                setFormData({ ...formData, recipient_id: e.target.value })
              }
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Select recipient</option>
              {users.map((user) => (
                <option key={user.id} value={user.id}>
                  {user.full_name} ({user.username})
                </option>
              ))}
              {/* Fallback: using a text input for demo */}
            </select>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Quantity *
              </label>
              <input
                type="number"
                min="1"
                max={selectedItem?.quantity || 1}
                value={formData.quantity}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    quantity: parseInt(e.target.value) || 1,
                  })
                }
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Issue Type *
              </label>
              <select
                value={formData.issue_type}
                onChange={(e) =>
                  setFormData({ ...formData, issue_type: e.target.value as any })
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="personal">Personal Use</option>
                <option value="classroom">Classroom</option>
                <option value="lab">Lab</option>
                <option value="teachers_room">Teacher's Room</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Due Date
            </label>
            <input
              type="date"
              value={formData.due_date?.split('T')[0] || ''}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  due_date: e.target.value
                    ? `${e.target.value}T23:59:59Z`
                    : undefined,
                })
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Notes
            </label>
            <textarea
              value={formData.notes || ''}
              onChange={(e) =>
                setFormData({ ...formData, notes: e.target.value })
              }
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>

        <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 flex justify-end gap-3">
          <button
            type="button"
            onClick={() => navigate(-1)}
            className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={submitting}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {submitting ? 'Issuing...' : 'Issue Item'}
          </button>
        </div>
      </form>
    </div>
  );
}