import { api } from "./client";
import type { AuthResult, UserResponse } from "../types/api";

export const authApi = {
  registerAnonymous: () =>
    api.post<AuthResult>("/auth/register", {}),

  getMe: () => api.get<UserResponse>("/auth/me"),
};
