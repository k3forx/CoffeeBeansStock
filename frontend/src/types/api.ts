// API response envelope
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: { field: string; message: string }[];
  };
}

// Auth
export interface UserResponse {
  id: string;
  name?: string;
  low_stock_threshold?: number;
  notification_enabled?: boolean;
  created_at: string;
  updated_at?: string;
}

export interface AuthResult {
  user: UserResponse;
  access_token: string;
  refresh_token: string;
}

export interface RefreshResult {
  access_token: string;
  refresh_token: string;
}

// Coffee Beans
export interface CoffeeBean {
  id: string;
  name: string;
  origin?: string;
  roast_level?: string;
  current_stock: number;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface ListBeansResult {
  beans: CoffeeBean[];
  pagination: {
    total: number;
    limit: number;
    offset: number;
    has_more: boolean;
  };
}

export interface CreateBeanInput {
  name: string;
  origin?: string;
  roast_level?: string;
  current_stock: number;
  notes?: string;
}

export interface UpdateBeanInput {
  name?: string;
  origin?: string;
  roast_level?: string;
  current_stock?: number;
  notes?: string;
}
