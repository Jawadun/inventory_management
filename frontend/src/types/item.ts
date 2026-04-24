export interface Item {
  id: string;
  name: string;
  category_id?: string;
  category?: Category;
  supplier_id?: string;
  supplier?: Supplier;
  sku?: string;
  barcode?: string;
  description?: string;
  quantity: number;
  min_quantity: number;
  unit: string;
  location?: string;
  storage_location?: string;
  purchase_date?: string;
  purchase_price?: number;
  warranty_months?: number;
  status: ItemStatus;
  condition?: string;
  image_url?: string;
  notes?: string;
  created_by?: string;
  created_at: string;
  updated_at: string;
}

export type ItemStatus = 'available' | 'issued' | 'reserved' | 'damaged' | 'retired';

export interface Category {
  id: string;
  name: string;
  description?: string;
  parent_id?: string;
  created_by?: string;
  created_at: string;
  updated_at: string;
}

export interface Supplier {
  id: string;
  name: string;
  contact_person?: string;
  phone?: string;
  email?: string;
  address?: string;
  notes?: string;
  is_active: boolean;
  created_by?: string;
  created_at: string;
  updated_at: string;
}

export interface ItemListResponse {
  items: Item[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ItemFilter {
  category_id?: string;
  supplier_id?: string;
  status?: string;
  search?: string;
  low_stock?: boolean;
}

export interface CreateItemRequest {
  name: string;
  category_id?: string;
  supplier_id?: string;
  sku?: string;
  barcode?: string;
  description?: string;
  quantity: number;
  min_quantity?: number;
  unit?: string;
  location?: string;
  storage_location?: string;
  purchase_date?: string;
  purchase_price?: number;
  warranty_months?: number;
  condition?: string;
  image_url?: string;
  notes?: string;
}

export interface UpdateItemRequest {
  name?: string;
  category_id?: string;
  supplier_id?: string;
  sku?: string;
  barcode?: string;
  description?: string;
  quantity?: number;
  min_quantity?: number;
  unit?: string;
  location?: string;
  storage_location?: string;
  purchase_date?: string;
  purchase_price?: number;
  warranty_months?: number;
  status?: string;
  condition?: string;
  image_url?: string;
  notes?: string;
}

export interface ItemHistory {
  id: string;
  item_id: string;
  quantity_change: number;
  previous_quantity: number;
  new_quantity: number;
  change_type: string;
  reason?: string;
  changed_by?: string;
  created_at: string;
}