export type RequestStatus = 'pending' | 'approved' | 'rejected' | 'cancelled' | 'fulfilled';
export type RequestType = 'classroom' | 'lab' | 'teachers_room' | 'personal';

export interface ItemRequest {
  id: string;
  user_id: string;
  user?: User;
  item_id: string;
  item?: Item;
  request_type: RequestType;
  quantity: number;
  status: RequestStatus;
  reason?: string;
  requested_at: string;
  reviewed_by?: string;
  reviewed_at?: string;
  rejection_reason?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateRequestRequest {
  item_id: string;
  quantity: number;
  request_type: string;
  reason?: string;
}

export interface ReviewRequestRequest {
  approved: boolean;
  notes?: string;
  rejection_reason?: string;
}

export interface RequestListResponse {
  requests: ItemRequest[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}