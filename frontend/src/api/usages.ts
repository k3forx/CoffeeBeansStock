import { api } from "./client";
import type {
  UsageHistory,
  CreateUsageInput,
  ListUsageHistoryResult,
} from "../types/api";

export const usagesApi = {
  list: (beanId: string, limit = 20, offset = 0) =>
    api.get<ListUsageHistoryResult>(
      `/coffee-beans/${beanId}/usages/?limit=${limit}&offset=${offset}`
    ),

  create: (beanId: string, input: CreateUsageInput) =>
    api.post<UsageHistory>(`/coffee-beans/${beanId}/usages/`, input),

  delete: (beanId: string, usageId: string) =>
    api.delete<{ message: string }>(`/coffee-beans/${beanId}/usages/${usageId}`),
};
