import { api, ApiError } from "../client";
import { useAuthStore } from "../../stores/auth";

const mockFetch = jest.fn();
global.fetch = mockFetch;

function jsonResponse(body: unknown, status = 200): Response {
  const text = JSON.stringify(body);
  return {
    ok: status >= 200 && status < 300,
    status,
    text: () => Promise.resolve(text),
    json: () => Promise.resolve(body),
  } as Response;
}

beforeEach(() => {
  mockFetch.mockReset();
  useAuthStore.setState({
    user: null,
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false,
    isLoading: false,
  });
});

describe("api.get", () => {
  it("makes a GET request and returns data", async () => {
    mockFetch.mockResolvedValueOnce(
      jsonResponse({ success: true, data: { id: "1", name: "Ethiopia" } })
    );

    const result = await api.get<{ id: string; name: string }>("/beans/1");

    expect(result).toEqual({ id: "1", name: "Ethiopia" });
    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/v1/beans/1",
      expect.objectContaining({
        headers: expect.objectContaining({ "Content-Type": "application/json" }),
      })
    );
  });

  it("attaches Authorization header when accessToken exists", async () => {
    useAuthStore.setState({ accessToken: "test-token" });
    mockFetch.mockResolvedValueOnce(
      jsonResponse({ success: true, data: {} })
    );

    await api.get("/beans");

    expect(mockFetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: "Bearer test-token",
        }),
      })
    );
  });
});

describe("api.post", () => {
  it("sends JSON body with POST method", async () => {
    mockFetch.mockResolvedValueOnce(
      jsonResponse({ success: true, data: { id: "2" } })
    );

    await api.post("/beans", { name: "Colombia" });

    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/v1/beans",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ name: "Colombia" }),
      })
    );
  });
});

describe("token refresh on 401", () => {
  it("refreshes token and retries on 401", async () => {
    useAuthStore.setState({
      accessToken: "expired-token",
      refreshToken: "valid-refresh",
    });

    // First call: 401
    mockFetch.mockResolvedValueOnce(
      jsonResponse({ success: false }, 401)
    );
    // Refresh call: success
    mockFetch.mockResolvedValueOnce(
      jsonResponse({
        success: true,
        data: { access_token: "new-access", refresh_token: "new-refresh" },
      })
    );
    // Retry call: success
    mockFetch.mockResolvedValueOnce(
      jsonResponse({ success: true, data: { id: "1" } })
    );

    const result = await api.get<{ id: string }>("/beans/1");

    expect(result).toEqual({ id: "1" });
    expect(mockFetch).toHaveBeenCalledTimes(3);

    // Verify retry used the new token
    const retryCall = mockFetch.mock.calls[2];
    expect(retryCall[1].headers.Authorization).toBe("Bearer new-access");
  });

  it("リフレッシュ失敗時にログアウトしてエラーをスローする", async () => {
    useAuthStore.setState({
      accessToken: "expired-token",
      refreshToken: "invalid-refresh",
      isAuthenticated: true,
    });

    // First call: 401
    mockFetch.mockResolvedValueOnce(
      jsonResponse({ success: false }, 401)
    );
    // Refresh call: failure
    mockFetch.mockResolvedValueOnce(
      jsonResponse({ success: false }, 401)
    );

    await expect(api.get("/beans")).rejects.toThrow("セッションが切れました");

    expect(useAuthStore.getState().isAuthenticated).toBe(false);
    expect(useAuthStore.getState().accessToken).toBeNull();
    expect(useAuthStore.getState().refreshToken).toBeNull();
  });

  it("リフレッシュトークンなしでサーバーがUNAUTHORIZEDを返した場合ログアウトする", async () => {
    useAuthStore.setState({
      accessToken: "some-token",
      refreshToken: null,
      isAuthenticated: true,
    });

    mockFetch.mockResolvedValueOnce(
      jsonResponse(
        {
          success: false,
          error: { code: "UNAUTHORIZED", message: "認証トークンが必要です" },
        },
        401
      )
    );

    await expect(api.get("/beans")).rejects.toThrow("認証トークンが必要です");

    expect(useAuthStore.getState().isAuthenticated).toBe(false);
    expect(useAuthStore.getState().accessToken).toBeNull();
    expect(useAuthStore.getState().refreshToken).toBeNull();
  });
});

describe("error handling", () => {
  it("throws ApiError on non-success response", async () => {
    mockFetch.mockResolvedValueOnce(
      jsonResponse({
        success: false,
        error: { code: "NOT_FOUND", message: "豆が見つかりません" },
      })
    );

    try {
      await api.get("/beans/999");
      fail("should have thrown");
    } catch (e) {
      expect(e).toBeInstanceOf(ApiError);
      expect((e as ApiError).code).toBe("NOT_FOUND");
      expect((e as ApiError).message).toBe("豆が見つかりません");
    }
  });

  it("throws ApiError with PARSE_ERROR on invalid JSON", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      status: 200,
      text: () => Promise.resolve("not json"),
    } as Response);

    try {
      await api.get("/beans");
      fail("should have thrown");
    } catch (e) {
      expect(e).toBeInstanceOf(ApiError);
      expect((e as ApiError).code).toBe("PARSE_ERROR");
    }
  });
});

describe("ApiError", () => {
  it("stores code, message, and details", () => {
    const details = [{ field: "name", message: "必須です" }];
    const err = new ApiError("VALIDATION", "バリデーションエラー", details);

    expect(err.code).toBe("VALIDATION");
    expect(err.message).toBe("バリデーションエラー");
    expect(err.details).toEqual(details);
    expect(err).toBeInstanceOf(Error);
  });
});
