import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { itemService, requestService } from '../../services';
import { Item, CreateRequestRequest } from '../../types';
import { LoadingSpinner } from '../common/Pagination';
import { SelectInput, SearchInput } from '../common/SelectInput';

const REQUEST_TYPE_OPTIONS = [
  { value: 'personal', label: 'Personal Use' },
  { value: 'classroom', label: 'Classroom' },
  { value: 'lab', label: 'Lab' },
  { value: 'teachers_room', label: "Teacher's Room" },
];

export function RequestForm() {
  const navigate = useNavigate();

  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const [formData, setFormData] = useState<CreateRequestRequest>({
    item_id: '',
    quantity: 1,
    request_type: 'personal',
    reason: '',
  });

  useEffect(() => {
    loadItems();
  }, []);

  const loadItems = async () => {
    setLoading(true);
    try {
      const res = await itemService.list({ status: 'available' }, 1, 100);
      setItems(res.items);
    } catch (err) {
      setError('Failed to load items');
    } finally {
      setLoading(false);
    }
  };

  const filteredItems = search
    ? items.filter(
        (item) =>
          item.name.toLowerCase().includes(search.toLowerCase()) ||
          item.sku?.toLowerCase().includes(search.toLowerCase())
      )
    : items;

  const selectedItem = items.find((item) => item.id === formData.item_id);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.item_id) {
      setError('Please select an item');
      return;
    }
    if (formData.quantity < 1) {
      setError('Quantity must be at least 1');
      return;
    }
    if (selectedItem && formData.quantity > selectedItem.quantity) {
      setError(`Maximum available quantity is ${selectedItem.quantity}`);
      return;
    }

    setError('');
    setSubmitting(true);
    try {
      await requestService.create(formData);
      navigate('/requests');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to submit request');
    } finally {
      setSubmitting(false);
    }
  };

  const handleItemSelect = (itemId: string) => {
    setFormData((prev) => ({
      ...prev,
      item_id: itemId,
      quantity: Math.min(prev.quantity, selectedItem?.quantity || 1),
    }));
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
          <h2 className="text-lg font-medium">Request Item</h2>
          <p className="text-sm text-gray-500 mt-1">
            Submit a request to borrow or obtain an item from the inventory
          </p>
        </div>

        <div className="p-6 space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Select Item *
            </label>
            <SearchInput
              value={search}
              onChange={setSearch}
              placeholder="Search by name or SKU..."
              className="mb-3"
            />
            <div className="max-h-60 overflow-y-auto border border-gray-200 rounded-md">
              {filteredItems.length === 0 ? (
                <div className="p-4 text-center text-gray-500">
                  No items available
                </div>
              ) : (
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">
                        Item
                      </th>
                      <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">
                        Available
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {filteredItems.map((item) => (
                      <tr
                        key={item.id}
                        onClick={() => handleItemSelect(item.id)}
                        className={`cursor-pointer hover:bg-gray-50 ${
                          formData.item_id === item.id
                            ? 'bg-blue-50'
                            : ''
                        }`}
                      >
                        <td className="px-4 py-3">
                          <div className="font-medium text-gray-900">
                            {item.name}
                          </div>
                          {item.sku && (
                            <div className="text-sm text-gray-500">
                              SKU: {item.sku}
                            </div>
                          )}
                        </td>
                        <td className="px-4 py-3">
                          <span
                            className={
                              item.quantity <= item.min_quantity
                                ? 'text-red-600 font-medium'
                                : ''
                            }
                          >
                            {item.quantity} {item.unit}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          </div>

          {selectedItem && (
            <div className="bg-gray-50 rounded-lg p-4">
              <h3 className="font-medium text-gray-900 mb-2">
                Selected Item
              </h3>
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <span className="text-gray-500">Name:</span>
                  <div className="font-medium">{selectedItem.name}</div>
                </div>
                <div>
                  <span className="text-gray-500">Available:</span>
                  <div
                    className={
                      selectedItem.quantity <= selectedItem.min_quantity
                        ? 'text-red-600 font-medium'
                        : ''
                    }
                  >
                    {selectedItem.quantity} {selectedItem.unit}
                  </div>
                </div>
                {selectedItem.location && (
                  <div>
                    <span className="text-gray-500">Location:</span>
                    <div className="font-medium">{selectedItem.location}</div>
                  </div>
                )}
                {selectedItem.category && (
                  <div>
                    <span className="text-gray-500">Category:</span>
                    <div className="font-medium">
                      {selectedItem.category.name}
                    </div>
                  </div>
                )}
              </div>
            </div>
          )}

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
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Request Type *
              </label>
              <select
                value={formData.request_type}
                onChange={(e) =>
                  setFormData({ ...formData, request_type: e.target.value as any })
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              >
                {REQUEST_TYPE_OPTIONS.map((opt) => (
                  <option key={opt.value} value={opt.value}>
                    {opt.label}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Reason (Optional)
            </label>
            <textarea
              value={formData.reason || ''}
              onChange={(e) =>
                setFormData({ ...formData, reason: e.target.value })
              }
              rows={3}
              placeholder="Explain why you need this item..."
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>

        <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 flex justify-end gap-3">
          <button
            type="button"
            onClick={() => navigate('/inventory')}
            className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={submitting || !formData.item_id}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {submitting ? 'Submitting...' : 'Submit Request'}
          </button>
        </div>
      </form>
    </div>
  );
}