import { useNavigate } from 'react-router-dom';
import { InventoryList } from '../../components/common/InventoryList';
import { Item } from '../../types';

export function AdminItemsPage() {
  const navigate = useNavigate();

  const handleItemSelect = (item: Item) => {
    navigate(`/admin/items/${item.id}`);
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Items</h1>
        <button
          onClick={() => navigate('/admin/items/new')}
          className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
        >
          <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Add Item
        </button>
      </div>
      <InventoryList isAdmin={true} onItemSelect={handleItemSelect} />
    </div>
  );
}