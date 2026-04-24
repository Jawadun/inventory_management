export interface Notice {
  id: string;
  title: string;
  content: string;
  posted_by?: string;
  is_pinned: boolean;
  is_active: boolean;
  priority: number;
  created_at: string;
  updated_at: string;
}

export interface CreateNoticeRequest {
  title: string;
  content: string;
  is_pinned?: boolean;
}

export interface UpdateNoticeRequest {
  title?: string;
  content?: string;
  is_pinned?: boolean;
  is_active?: boolean;
  priority?: number;
}