import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { noticeService } from '../../services';
import { Notice, CreateNoticeRequest, UpdateNoticeRequest } from '../../types';
import { LoadingSpinner, EmptyState, ConfirmDialog } from '../common/Pagination';
import { NoticeCard, NoticeDetail } from './NoticeCard';

export function AdminNotices() {
  const [notices, setNotices] = useState<Notice[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedNotice, setSelectedNotice] = useState<Notice | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [editingNotice, setEditingNotice] = useState<Notice | null>(null);
  const [deleteNotice, setDeleteNotice] = useState<Notice | null>(null);

  const fetchNotices = useCallback(async () => {
    setLoading(true);
    try {
      const data = await noticeService.list(false);
      setNotices(data);
    } catch (err) {
      console.error('Failed to fetch notices:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchNotices();
  }, [fetchNotices]);

  const handleDelete = async () => {
    if (!deleteNotice) return;
    try {
      await noticeService.delete(deleteNotice.id);
      setDeleteNotice(null);
      fetchNotices();
    } catch (err: any) {
      alert(err.response?.data?.message || 'Failed to delete notice');
    }
  };

  const handleEdit = (notice: Notice) => {
    setEditingNotice(notice);
    setShowForm(true);
  };

  const handleNew = () => {
    setEditingNotice(null);
    setShowForm(true);
  };

  const closeForm = () => {
    setShowForm(false);
    setEditingNotice(null);
  };

  const handleSave = async (data: CreateNoticeRequest | UpdateNoticeRequest) => {
    try {
      if (editingNotice) {
        await noticeService.update(editingNotice.id, data as UpdateNoticeRequest);
      } else {
        await noticeService.create(data as CreateNoticeRequest);
      }
      closeForm();
      fetchNotices();
    } catch (err: any) {
      throw err;
    }
  };

  const pinnedCount = notices.filter((n) => n.is_pinned).length;
  const activeCount = notices.filter((n) => n.is_active).length;

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Notices</h1>
            <p className="text-sm text-gray-500 mt-1">
              {notices.length} notices ({activeCount} active, {pinnedCount} pinned)
            </p>
          </div>
          <button
            onClick={handleNew}
            className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
          >
            <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            New Notice
          </button>
        </div>
      </div>

      {loading ? (
        <div className="py-12">
          <LoadingSpinner size="lg" />
        </div>
      ) : notices.length === 0 ? (
        <EmptyState
          title="No notices"
          description="Create your first notice to get started"
          action={
            <button
              onClick={handleNew}
              className="mt-4 inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
            >
              Create Notice
            </button>
          }
        />
      ) : (
        <div className="space-y-4">
          {notices.map((notice) => (
            <NoticeCard
              key={notice.id}
              notice={notice}
              isAdmin={true}
              onEdit={() => handleEdit(notice)}
              onDelete={() => setDeleteNotice(notice)}
            />
          ))}
        </div>
      )}

      {showForm && (
        <NoticeForm
          notice={editingNotice}
          onSave={handleSave}
          onCancel={closeForm}
        />
      )}

      <ConfirmDialog
        open={!!deleteNotice}
        title="Delete Notice"
        message={`Are you sure you want to delete "${deleteNotice?.title}"? This action cannot be undone.`}
        confirmLabel="Delete"
        variant="danger"
        onConfirm={handleDelete}
        onCancel={() => setDeleteNotice(null)}
      />
    </div>
  );
}

interface NoticeFormProps {
  notice: Notice | null;
  onSave: (data: CreateNoticeRequest | UpdateNoticeRequest) => Promise<void>;
  onCancel: () => void;
}

function NoticeForm({ notice, onSave, onCancel }: NoticeFormProps) {
  const [formData, setFormData] = useState<CreateNoticeRequest>({
    title: notice?.title || '',
    content: notice?.content || '',
    is_pinned: notice?.is_pinned || false,
  });
  const [isActive, setIsActive] = useState(notice?.is_active ?? true);
  const [priority, setPriority] = useState(notice?.priority || 0);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSubmitting(true);
    try {
      await onSave({
        ...formData,
        is_pinned: formData.is_pinned,
      });
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to save notice');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="fixed inset-0 bg-black bg-opacity-50" onClick={onCancel} />
      <div className="relative bg-white rounded-lg shadow-xl max-w-xl w-full max-h-[90vh] overflow-y-auto">
        <form onSubmit={handleSubmit}>
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-lg font-medium">
              {notice ? 'Edit Notice' : 'New Notice'}
            </h2>
          </div>

          <div className="p-6 space-y-4">
            {error && (
              <div className="bg-red-50 text-red-600 px-4 py-3 rounded text-sm">
                {error}
              </div>
            )}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Title *
              </label>
              <input
                type="text"
                value={formData.title}
                onChange={(e) =>
                  setFormData({ ...formData, title: e.target.value })
                }
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter notice title"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Content *
              </label>
              <textarea
                value={formData.content}
                onChange={(e) =>
                  setFormData({ ...formData, content: e.target.value })
                }
                required
                rows={6}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter notice content"
              />
            </div>

            <div className="flex gap-4">
              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={formData.is_pinned}
                  onChange={(e) =>
                    setFormData({ ...formData, is_pinned: e.target.checked })
                  }
                  className="rounded border-gray-300"
                />
                <span className="text-sm text-gray-700">Pin notice</span>
              </label>

              {notice && (
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={isActive}
                    onChange={(e) => setIsActive(e.target.checked)}
                    className="rounded border-gray-300"
                  />
                  <span className="text-sm text-gray-700">Active</span>
                </label>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Priority (0-9, lower is higher priority)
              </label>
              <input
                type="number"
                min="0"
                max="9"
                value={priority}
                onChange={(e) => setPriority(parseInt(e.target.value) || 0)}
                className="w-24 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          </div>

          <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 flex justify-end gap-3 rounded-b-lg">
            <button
              type="button"
              onClick={onCancel}
              className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {submitting ? 'Saving...' : notice ? 'Update' : 'Publish'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}