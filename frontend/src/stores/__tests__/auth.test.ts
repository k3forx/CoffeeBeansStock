import * as SecureStore from "expo-secure-store";
import { useAuthStore } from "../auth";
import type { UserResponse } from "../../types/api";

const mockedSecureStore = jest.mocked(SecureStore);

const mockUser: UserResponse = {
  id: "user-1",
  name: "testuser",
  created_at: "2026-01-01T00:00:00Z",
};

describe("useAuthStore", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    useAuthStore.setState({
      user: null,
      accessToken: null,
      refreshToken: null,
      isLoading: true,
      isAuthenticated: false,
    });
  });

  it("setAuthでトークンを保存し認証状態をセットする", async () => {
    await useAuthStore.getState().setAuth(mockUser, "access-123", "refresh-456");

    expect(mockedSecureStore.setItemAsync).toHaveBeenCalledWith("access_token", "access-123");
    expect(mockedSecureStore.setItemAsync).toHaveBeenCalledWith("refresh_token", "refresh-456");

    const state = useAuthStore.getState();
    expect(state.user).toEqual(mockUser);
    expect(state.accessToken).toBe("access-123");
    expect(state.refreshToken).toBe("refresh-456");
    expect(state.isAuthenticated).toBe(true);
  });

  it("setTokensでユーザーを変更せずトークンのみ更新する", async () => {
    await useAuthStore.getState().setTokens("new-access", "new-refresh");

    expect(mockedSecureStore.setItemAsync).toHaveBeenCalledWith("access_token", "new-access");
    expect(mockedSecureStore.setItemAsync).toHaveBeenCalledWith("refresh_token", "new-refresh");

    const state = useAuthStore.getState();
    expect(state.accessToken).toBe("new-access");
    expect(state.refreshToken).toBe("new-refresh");
    expect(state.user).toBeNull();
  });

  it("setUserでユーザーのみ更新する", () => {
    useAuthStore.getState().setUser(mockUser);

    const state = useAuthStore.getState();
    expect(state.user).toEqual(mockUser);
    expect(state.accessToken).toBeNull();
    expect(state.refreshToken).toBeNull();
  });

  it("logoutでSecureStoreからトークンを削除し状態をリセットする", async () => {
    useAuthStore.setState({
      user: mockUser,
      accessToken: "access-123",
      refreshToken: "refresh-456",
      isAuthenticated: true,
    });

    await useAuthStore.getState().logout();

    expect(mockedSecureStore.deleteItemAsync).toHaveBeenCalledWith("access_token");
    expect(mockedSecureStore.deleteItemAsync).toHaveBeenCalledWith("refresh_token");

    const state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.accessToken).toBeNull();
    expect(state.refreshToken).toBeNull();
    expect(state.isAuthenticated).toBe(false);
  });

  it("loadStoredTokensでトークンが存在する場合に認証状態を復元する", async () => {
    mockedSecureStore.getItemAsync.mockImplementation(async (key) => {
      if (key === "access_token") return "stored-access";
      if (key === "refresh_token") return "stored-refresh";
      return null;
    });

    await useAuthStore.getState().loadStoredTokens();

    const state = useAuthStore.getState();
    expect(state.accessToken).toBe("stored-access");
    expect(state.refreshToken).toBe("stored-refresh");
    expect(state.isAuthenticated).toBe(true);
    expect(state.isLoading).toBe(false);
  });

  it("loadStoredTokensでトークンがない場合にisLoadingをfalseにする", async () => {
    mockedSecureStore.getItemAsync.mockResolvedValue(null);

    await useAuthStore.getState().loadStoredTokens();

    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(false);
    expect(state.isLoading).toBe(false);
  });

  it("loadStoredTokensでエラー時もisLoadingをfalseにする", async () => {
    mockedSecureStore.getItemAsync.mockRejectedValue(new Error("SecureStore error"));

    await useAuthStore.getState().loadStoredTokens().catch(() => {});

    const state = useAuthStore.getState();
    expect(state.isLoading).toBe(false);
  });
});
