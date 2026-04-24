import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { itemService } from '../../services';
import { Item, ItemFilter, Category, Supplier } from '../../types';
import { Pagination, LoadingSpinner, EmptyState } from '../common/Pagination';
import { SearchInput, SelectInput } from '../common/SelectInput';

interface InventoryListProps {
  isAdmin?: boolean;
  onItemSelect?: (item: Item) => void;
}

const STATUS_OPTIONS = [
  { value: '', label: 'All Status' },
  { value: 'available', label: 'Available' },
  { value: 'issued', label: 'Issued' },
  { value: 'reserved', label: 'Reserved' },
  { value: 'damaged', label: 'Damaged' },
  { value: 'retired', label: 'Retired' },
];

export function InventoryList({ isAdmin = false, onItemSelect }: InventoryListProps) {
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalCount, setTotalCount] = useState(0);

  const [categories, setCategories] = useState<Category[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);

  const [filter, setFilter] = useState<ItemFilter>({
    search: '',
    category_id: '',
    supplier_id: '',
    status: '',
    low_stock: false,
  });

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const [itemsRes, categoriesRes, suppliersRes] = await Promise.all([
        itemService.list(
          {
            search: filter.search || undefined,
            category_id: filter.category_id || undefined,
            supplier_id: filter.supplier_id || undefined,
            status: filter.status || undefined,
            low_stock: filter.low_stock || undefined,
          },
          page,
          20
        ),
        isAdmin ? itemService.listCategories() : Promise.resolve([]),
        isAdmin ? itemService.listSuppliers() : Promise.resolve([]),
      ]);

      setItems(itemsRes.items);
      setTotalPages(itemsRes.total_pages);
      setTotalCount(itemsRes.total_count);

      if (isAdmin) {
        setCategories(categoriesRes);
        setSuppliers(suppliersRes);
      }
    } catch (error) {
      console.error('Failed to fetch items:', error);
    } finally {
      setLoading(false);
    }
  }, [filter, page, isAdmin]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleFilterChange = (key: keyof ItemFilter, value: string | boolean) => {
    setFilter((prev) => ({ ...prev, [key]: value }));
    setPage(1);
  };

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };

  const categoryOptions = [
    { value: '', label: 'All Categories' },
    ...categories.map((c) => ({ value: c.id, label: c.name })),
  ];

  const supplierOptions = [
    { value: '', label: 'All Suppliers' },
    ...suppliers.map((s) => ({ value: s.id, label: s.name })),
  ];

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Inventory</h1>
        <p className="text-sm text-gray-500 mt-1">
          {totalCount} {totalCount === 1 ? 'item' : 'items'} total
        </p>
      </div>

      <div className="bg-white rounded-lg shadow mb-6">
        <div className="p-4 border-b border-gray-200">
          <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
            <div className="md:col-span-2">
              <SearchInput
                value={filter.search || ''}
                onChange={(value) => handleFilterChange('search', value)}
                placeholder="Search by name, SKU, barcode..."
                onSubmit={() => handleFilterChange('search', filter.search || '')}
              />
            </div>
            {isAdmin ? (
              <>
                <SelectInput
                  value={filter.category_id || ''}
                  onChange={(value) => handleFilterChange('category_id', value)}
                  options={categoryOptions}
                />
                <SelectInput
                  value={filter.supplier_id || ''}
                  onChange={(value) => handleFilterChange('supplier_id', value)}
                  options={supplierOptions}
                />
              </>
            ) : null}
            <SelectInput
              value={filter.status || ''}
              onChange={(value) => handleFilterChange('status', value)}
              options={STATUS_OPTIONS}
            />
            {isAdmin && (
              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={filter.low_stock || false}
                  onChange={(e) => handleFilterChange('low_stock', e.target.checked)}
                  className="rounded border-gray-300"
                />
                <span className="text-sm text-gray-700">Low Stock Only</span>
              </label>
            )}
          </div>
        </div>

        {loading ? (
          <div className="py-12">
            <LoadingSpinner size="lg" />
          </div>
        ) : items.length === 0 ? (
          <EmptyState
            title="No items found"
            description={filter.search || filter.category_id || filter.supplier_id || filter.status || filter.low_stock
              ? 'Try adjusting your filters'
              : 'Get started by adding your first item'}
            action={
              isAdmin ? (
                <Link
                  to="/admin/items/new"
                  className="mt-4 inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
                >
                  Add Item
                </Link>
              ) : undefined
            }
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Item
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    SKU
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Quantity
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  {isAdmin && (
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Location
                    </th>
                  )}
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Category
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {items.map((item) => (
                  <tr
                    key={item.id}
                    className="hover:bg-gray-50 cursor-pointer"
                    onClick={() => onItemSelect?.(item)}
                  >
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        {item.image_url && (
                          <img
                            src={item.image_url}
                            alt={item.name}
                            className="h-10 w-10 rounded object-cover mr-4"
                          />
                        )}
                        <div>
                          <div className="text-sm font-medium text-gray-900">
                            {item.name}
                          </div>
                          {item.description && (
                            <div className="text-sm text-gray-500 truncate max-w-xs">
                              {item.description}
                            </div>
                          )}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {item.sku || '-'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        <span className={item.quantity <= item.min_quantity ? 'text-red-600 font-semibold' : ''}>
                          {item.quantity}
                        </span>{' '}
                        <span className="text-gray-500">{item.unit}</span>
                      </div>
                      {item.quantity <= item.min_quantity && (
                        <div className="text-xs text-red-600">
                          Min: {item.min_quantity}
                        </div>
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <StatusBadge status={item.status} />
                    </td>
                    {isAdmin && (
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {item.location || '-'}
                      </td>
                    )}
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {item.category?.name || '-'}
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
          onPageChange={handlePageChange}
        />
      </div>
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