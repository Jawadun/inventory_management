export interface PublicStats {
  total_items: number;
  total_categories: number;
  total_suppliers: number;
  available_items: number;
  issued_items: number;
  pending_requests: number;
  active_notices: number;
  low_stock_items: number;
  updated_at: string;
}

export interface CategoryStat {
  category_id: string;
  category_name: string;
  item_count: number;
  total_quantity: number;
}

export interface DashboardStats {
  public: PublicStats;
  categories: CategoryStat[];
  issued_by_type: Record<string, number>;
  updated_at: string;
}