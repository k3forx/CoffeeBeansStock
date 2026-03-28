import { create } from "zustand";
import * as SecureStore from "expo-secure-store";
import type { UserResponse } from "../types/api";

interface AuthState {
  user: UserResponse | null;
  accessToken: string | null;
  refreshToken: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;

  setAuth: (user: UserResponse, accessToken: string, refreshToken: string) => Promise<void>;
  setTokens: (accessToken: string, refreshToken: string) => Promise<void>;
  logout: () => Promise<void>;
  loadStoredTokens: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  accessToken: null,
  refreshToken: null,
  isLoading: true,
  isAuthenticated: false,

  setAuth: async (user, accessToken, refreshToken) => {
    await SecureStore.setItemAsync("access_token", accessToken);
    await SecureStore.setItemAsync("refresh_token", refreshToken);
    set({ user, accessToken, refreshToken, isAuthenticated: true });
  },

  setTokens: async (accessToken, refreshToken) => {
    await SecureStore.setItemAsync("access_token", accessToken);
    await SecureStore.setItemAsync("refresh_token", refreshToken);
    set({ accessToken, refreshToken });
  },

  logout: async () => {
    await SecureStore.deleteItemAsync("access_token");
    await SecureStore.deleteItemAsync("refresh_token");
    set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false });
  },

  loadStoredTokens: async () => {
    try {
      const accessToken = await SecureStore.getItemAsync("access_token");
      const refreshToken = await SecureStore.getItemAsync("refresh_token");
      if (accessToken && refreshToken) {
        set({ accessToken, refreshToken, isAuthenticated: true });
      }
    } finally {
      set({ isLoading: false });
    }
  },
}));
