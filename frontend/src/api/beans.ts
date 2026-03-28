import { api } from "./client";
import type { CoffeeBean, CreateBeanInput, ListBeansResult, UpdateBeanInput } from "../types/api";

export const beansApi = {
  list: (limit = 20, offset = 0) =>
    api.get<ListBeansResult>(`/coffee-beans/?limit=${limit}&offset=${offset}`),

  get: (id: string) => api.get<CoffeeBean>(`/coffee-beans/${id}`),

  create: (input: CreateBeanInput) => api.post<CoffeeBean>("/coffee-beans/", input),

  update: (id: string, input: UpdateBeanInput) =>
    api.put<CoffeeBean>(`/coffee-beans/${id}`, input),

  delete: (id: string) => api.delete<{ message: string }>(`/coffee-beans/${id}`),
};
