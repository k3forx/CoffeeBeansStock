import { useAuthStore } from "../stores/auth";
import type { ApiResponse } from "../types/api";

// TODO: 環境変数に切り出す
const API_BASE = "http://localhost:8080/api/v1";

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const { accessToken, refreshToken, setTokens } = useAuthStore.getState();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  if (accessToken) {
    headers["Authorization"] = `Bearer ${accessToken}`;
  }

  let res = await fetch(`${API_BASE}${path}`, { ...options, headers });

  // Token expired → try refresh
  if (res.status === 401 && refreshToken) {
    const refreshRes = await fetch(`${API_BASE}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (refreshRes.ok) {
      const refreshData: ApiResponse<{ access_token: string; refresh_token: string }> =
        await refreshRes.json();
      if (refreshData.success && refreshData.data) {
        await setTokens(refreshData.data.access_token, refreshData.data.refresh_token);
        headers["Authorization"] = `Bearer ${refreshData.data.access_token}`;
        res = await fetch(`${API_BASE}${path}`, { ...options, headers });
      }
    } else {
      throw new ApiError("UNAUTHORIZED", "セッションが切れました。再ログインしてください。");
    }
  }

  const text = await res.text();
  let json: ApiResponse<T>;
  try {
    json = JSON.parse(text);
  } catch {
    throw new ApiError("PARSE_ERROR", `サーバーエラー (${res.status}): ${text.slice(0, 200)}`);
  }

  if (!json.success) {
    throw new ApiError(
      json.error?.code ?? "UNKNOWN",
      json.error?.message ?? "エラーが発生しました",
      json.error?.details
    );
  }

  return json.data as T;
}

export class ApiError extends Error {
  code: string;
  details?: { field: string; message: string }[];

  constructor(code: string, message: string, details?: { field: string; message: string }[]) {
    super(message);
    this.code = code;
    this.details = details;
  }
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "POST", body: JSON.stringify(body) }),
  put: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PUT", body: JSON.stringify(body) }),
  delete: <T>(path: string) => request<T>(path, { method: "DELETE" }),
};
