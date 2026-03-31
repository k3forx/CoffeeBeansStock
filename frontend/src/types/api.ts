import type { components } from "./generated";

// Re-export generated types with existing names to maintain backward compatibility.
// All types originate from the OpenAPI spec (backend/api/openapi.yaml).
export type UserResponse = components["schemas"]["UserResponse"];
export type AuthResult = components["schemas"]["AuthResponse"];
export type RefreshResult = components["schemas"]["RefreshResult"];

export type ConsumptionRate = components["schemas"]["ConsumptionRate"];
export type CoffeeBean = components["schemas"]["CoffeeBeanResponse"];
export type CreateBeanInput = components["schemas"]["CreateBeanRequest"];
export type UpdateBeanInput = components["schemas"]["UpdateBeanRequest"];
export type ListBeansResult = components["schemas"]["ListBeansResponse"];
export type RoastLevel = components["schemas"]["RoastLevel"];
export type RoastDetail = components["schemas"]["RoastDetail"];
export type CreateUsageInput = components["schemas"]["CreateUsageHistoryRequest"];
export type UsageHistory = components["schemas"]["UsageHistoryResponse"];
export type ListUsageHistoryResult = components["schemas"]["ListUsageHistoryResponse"];
export type UpdateUserInput = components["schemas"]["UpdateUserRequest"];

// Envelope type used internally by client.ts
export type ApiResponse<T> = {
  success: boolean;
  data?: T;
  error?: components["schemas"]["ErrorBody"];
};
