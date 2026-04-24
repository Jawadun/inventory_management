import { ReactNode } from 'react';

export interface Column<T> {
  key: string;
  header: string;
  render?: (item: T) => ReactNode;
  className?: string;
}

interface AdminTableProps<T> {
  data: T[];
  columns: Column<T>[];
  loading?: boolean;
  onRowClick?: (item: T) => void;
  emptyMessage?: string;
  keyField: keyof T;
}

export function AdminTable<T>({
  data,
  columns,
  loading,
  onRowClick,
  emptyMessage = 'No data available',
  keyField,
}: AdminTableProps<T>) {
  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="text-center py-12 text-gray-500">
        {emptyMessage}
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            {columns.map((col) => (
              <th
                key={col.key}
                className={`px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider ${
                  col.className || ''
                }`}
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {data.map((item, idx) => (
            <tr
              key={String(item[keyField])}
              className={`hover:bg-gray-50 ${
                onRowClick ? 'cursor-pointer' : ''
              }`}
              onClick={() => onRowClick?.(item)}
            >
              {columns.map((col) => (
                <td
                  key={col.key}
                  className={`px-6 py-4 whitespace-nowrap text-sm ${
                    col.className || ''
                  }`}
                >
                  {col.render
                    ? col.render(item)
                    : String((item as any)[col.key] || '')}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

interface FilterBarProps {
  children: ReactNode;
  actions?: ReactNode;
}

export function FilterBar({ children, actions }: FilterBarProps) {
  return (
    <div className="p-4 border-b border-gray-200">
      <div className="flex flex-wrap items-center gap-4">
        {children}
        {actions && <div className="ml-auto">{actions}</div>}
      </div>
    </div>
  );
}

interface PaginationBarProps {
  page: number;
  totalPages: number;
  totalCount: number;
  onPageChange: (page: number) => void;
}

export function PaginationBar({
  page,
  totalPages,
  totalCount,
  onPageChange,
}: PaginationBarProps) {
  if (totalPages <= 1) {
    return (
      <div className="px-6 py-3 border-t border-gray-200 text-sm text-gray-500">
        Showing {totalCount} {totalCount === 1 ? 'item' : 'items'}
      </div>
    );
  }

  return (
    <div className="px-6 py-3 border-t border-gray-200 flex items-center justify-between">
      <div className="text-sm text-gray-500">
        Showing {(page - 1) * 20 + 1} to {Math.min(page * 20, totalCount)} of{' '}
        {totalCount} items
      </div>
      <div className="flex gap-2">
        <button
          onClick={() => onPageChange(page - 1)}
          disabled={page <= 1}
          className="px-3 py-1 border rounded text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
        >
          Previous
        </button>
        <button
          onClick={() => onPageChange(page + 1)}
          disabled={page >= totalPages}
          className="px-3 py-1 border rounded text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
        >
          Next
        </button>
      </div>
    </div>
  );
}