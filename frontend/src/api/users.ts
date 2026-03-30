import { api } from "./client";
import type { UpdateUserInput, UserResponse } from "../types/api";

export const usersApi = {
  updateMe: (input: UpdateUserInput) =>
    api.put<UserResponse>("/auth/me", input),
};
