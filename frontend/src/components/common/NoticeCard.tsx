import { Notice } from '../../types';

interface NoticeCardProps {
  notice: Notice;
  onClick?: () => void;
  isAdmin?: boolean;
  onEdit?: () => void;
  onDelete?: () => void;
}

export function NoticeCard({ notice, onClick, isAdmin, onEdit, onDelete }: NoticeCardProps) {
  const formattedDate = new Date(notice.created_at).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <div
      className={`bg-white rounded-lg shadow hover:shadow-md transition-shadow ${
        notice.is_pinned ? 'border-l-4 border-blue-500' : ''
      } ${onClick ? 'cursor-pointer' : ''}`}
      onClick={onClick}
    >
      <div className="p-6">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              {notice.is_pinned && (
                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-700">
                  <svg className="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M5 5a2 2 0 012-2h6a2 2 0 012 2v2H5V5zM3 7a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 13a1 1 0 011-1h8a1 1 0 110 2H4a1 1 0 01-1-1z" />
                  </svg>
                  Pinned
                </span>
              )}
              {!notice.is_active && (
                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-500">
                  Inactive
                </span>
              )}
              {notice.priority > 0 && notice.priority < 3 && (
                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-700">
                  Priority {notice.priority}
                </span>
              )}
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              {notice.title}
            </h3>
            <div className="text-gray-600 whitespace-pre-wrap">
              {notice.content.length > 200
                ? `${notice.content.substring(0, 200)}...`
                : notice.content}
            </div>
            <div className="mt-3 text-sm text-gray-400">
              Posted {formattedDate}
            </div>
          </div>

          {isAdmin && (onEdit || onDelete) && (
            <div className="flex flex-col gap-2 ml-4">
              {onEdit && (
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onEdit();
                  }}
                  className="p-2 text-gray-400 hover:text-gray-600 rounded hover:bg-gray-100"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
              )}
              {onDelete && (
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onDelete();
                  }}
                  className="p-2 text-gray-400 hover:text-red-600 rounded hover:bg-gray-100"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

interface NoticeDetailProps {
  notice: Notice;
  onClose: () => void;
  isAdmin?: boolean;
  onEdit?: () => void;
  onDelete?: () => void;
}

export function NoticeDetail({ notice, onClose, isAdmin, onEdit, onDelete }: NoticeDetailProps) {
  const formattedDate = new Date(notice.created_at).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-6 border-b border-gray-200">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-2">
            {notice.is_pinned && (
              <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-700">
                <svg className="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M5 5a2 2 0 012-2h6a2 2 0 012 2v2H5V5zM3 7a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 13a1 1 0 011-1h8a1 1 0 110 2H4a1 1 0 01-1-1z" />
                </svg>
                Pinned
              </span>
            )}
            {!notice.is_active && (
              <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-500">
                Inactive
              </span>
            )}
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 rounded hover:bg-gray-100"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <h2 className="text-2xl font-bold text-gray-900 mt-2">
          {notice.title}
        </h2>
        <div className="mt-2 text-sm text-gray-500">
          Posted {formattedDate}
        </div>
      </div>

      <div className="p-6">
        <div className="prose prose-sm max-w-none text-gray-600 whitespace-pre-wrap">
          {notice.content}
        </div>
      </div>

      {isAdmin && (onEdit || onDelete) && (
        <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 flex justify-end gap-3 rounded-b-lg">
          {onDelete && (
            <button
              onClick={onDelete}
              className="px-4 py-2 text-red-600 hover:bg-red-50 rounded"
            >
              Delete
            </button>
          )}
          {onEdit && (
            <button
              onClick={onEdit}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Edit
            </button>
          )}
        </div>
      )}
    </div>
  );
}