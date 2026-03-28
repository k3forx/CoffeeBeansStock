---
paths:
  - "backend/**/*_test.go"
---

# バックエンド テスト方針

## Repositoryテストの基盤

- `testhelper_test.go` の `TestMain` で testcontainers-go により PostgreSQL コンテナを起動し、マイグレーションを適用
- `newTestQueries(t)` ヘルパーでテストごとにトランザクションを開始し、`t.Cleanup` で自動ロールバック
- 各テストはトランザクション内で実行され、テスト間のデータ干渉なし

## テストの書き方ルール

- テーブル駆動テスト: `map[string]struct{}` 形式
- `t.Parallel()` を Test関数・サブテスト両方で使用
- アサーション: `go-cmp` を使用（testifyは使わない）
- 構造体の比較は `cmp.Diff` で構造体全体を比較する（フィールドごとの個別アサーションは行わない）
- 比較対象外のフィールド（タイムスタンプ等）は `cmp.FilterPath` + `cmp.Ignore()` で除外

## 新規repositoryテスト テンプレート

```go
func TestXxxRepository_Method(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		want    Type
		wantErr bool
	}{
		"description": {
			want: Type{/* expected fields */},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := NewXxxRepository(newTestQueries(t))

			got, err := repo.Method(t.Context())

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.want, got, xxxCmpOpts()...); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
```
