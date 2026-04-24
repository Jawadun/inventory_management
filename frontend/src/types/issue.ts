export type IssueType = 'classroom' | 'lab' | 'teachers_room' | 'personal';
export type IssueStatus = 'pending' | 'approved' | 'issued' | 'returned' | 'overdue' | 'rejected';

export interface IssueRecord {
  id: string;
  request_id?: string;
  item_id: string;
  item?: Item;
  recipient_id: string;
  recipient?: User;
  issued_by?: string;
  quantity: number;
  issue_type: IssueType;
  issue_date: string;
  due_date?: string;
  actual_return_date?: string;
  return_condition?: string;
  return_remarks?: string;
  status: IssueStatus;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateIssueRequest {
  item_id: string;
  recipient_id: string;
  quantity: number;
  issue_type: string;
  due_date?: string;
  notes?: string;
  auto_approve?: boolean;
}

export interface CreateReturnRequest {
  return_condition?: string;
  return_remarks?: string;
}

export interface IssueListResponse {
  issues: IssueRecord[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}