import { useState, useEffect, useCallback } from 'react';
import { noticeService } from '../../services';
import { Notice } from '../../types';
import { LoadingSpinner, EmptyState } from '../common/Pagination';
import { NoticeCard, NoticeDetail } from './NoticeCard';

export function NoticeBoard() {
  const [notices, setNotices] = useState<Notice[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedNotice, setSelectedNotice] = useState<Notice | null>(null);
  const [filter, setFilter] = useState<'all' | 'active' | 'pinned'>('active');

  const fetchNotices = useCallback(async () => {
    setLoading(true);
    try {
      const data = await noticeService.list(filter === 'all');
      setNotices(data);
    } catch (err) {
      console.error('Failed to fetch notices:', err);
    } finally {
      setLoading(false);
    }
  }, [filter]);

  useEffect(() => {
    fetchNotices();
  }, [fetchNotices]);

  const pinnedNotices = notices.filter((n) => n.is_pinned);
  const regularNotices = notices.filter((n) => !n.is_pinned);

  return (
    <div className="max-w-3xl mx-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Notice Board</h1>
        <p className="text-sm text-gray-500 mt-1">
          Stay updated with latest announcements
        </p>
      </div>

      <div className="mb-4 flex gap-2">
        <button
          onClick={() => setFilter('active')}
          className={`px-3 py-1.5 rounded text-sm font-medium ${
            filter === 'active'
              ? 'bg-blue-100 text-blue-700'
              : 'text-gray-600 hover:bg-gray-100'
          }`}
        >
          Active
        </button>
        <button
          onClick={() => setFilter('all')}
          className={`px-3 py-1.5 rounded text-sm font-medium ${
            filter === 'all'
              ? 'bg-blue-100 text-blue-700'
              : 'text-gray-600 hover:bg-gray-100'
          }`}
        >
          All Notices
        </button>
      </div>

      {loading ? (
        <div className="py-12">
          <LoadingSpinner size="lg" />
        </div>
      ) : notices.length === 0 ? (
        <EmptyState
          title="No notices"
          description={filter === 'active' ? 'No active notices at the moment' : 'No notices available'}
        />
      ) : (
        <div className="space-y-4">
          {pinnedNotices.length > 0 && (
            <div>
              <div className="flex items-center gap-2 mb-3">
                <svg className="w-4 h-4 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M5 5a2 2 0 012-2h6a2 2 0 012 2v2H5V5zM3 7a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 13a1 1 0 011-1h8a1 1 0 110 2H4a1 1 0 01-1-1z" />
                </svg>
                <span className="text-sm font-medium text-gray-500">Pinned</span>
              </div>
              <div className="space-y-4">
                {pinnedNotices.map((notice) => (
                  <NoticeCard
                    key={notice.id}
                    notice={notice}
                    onClick={() => setSelectedNotice(notice)}
                  />
                ))}
              </div>
            </div>
          )}

          {regularNotices.length > 0 && (
            <div>
              {pinnedNotices.length > 0 && (
                <div className="flex items-center gap-2 mb-3 mt-6">
                  <span className="text-sm font-medium text-gray-500">Announcements</span>
                </div>
              )}
              <div className="space-y-4">
                {regularNotices.map((notice) => (
                  <NoticeCard
                    key={notice.id}
                    notice={notice}
                    onClick={() => setSelectedNotice(notice)}
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {selectedNotice && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div
            className="fixed inset-0 bg-black bg-opacity-50"
            onClick={() => setSelectedNotice(null)}
          />
          <div className="relative max-w-2xl w-full max-h-[90vh] overflow-y-auto">
            <NoticeDetail
              notice={selectedNotice}
              onClose={() => setSelectedNotice(null)}
            />
          </div>
        </div>
      )}
    </div>
  );
}