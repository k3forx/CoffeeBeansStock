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
- `t.Fatal` / `t.Fatalf` は使わない。`t.Errorf` + `return` を使う
- **比較方法の使い分け**:
  - VO・プリミティブ（bool, string, int 等）: `!=` で直接比較
  - エンティティ・複合構造体: `cmp.Diff` で構造体全体を比較する（フィールドごとの個別アサーションは行わない）
- 構造体に unexported フィールドがある場合は `cmp.AllowUnexported(T{})` を使う（`cmpopts.IgnoreUnexported` は使わない）
- 比較対象外のフィールド（タイムスタンプ等）は `cmp.FilterPath` + `cmp.Ignore()` で除外

## Domain テスト テンプレート

### VO ファクトリ関数（`NewXxx`）

```go
func TestNewXxx(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		want    Xxx
		wantErr bool
	}{
		"valid_foo": {input: "foo", want: XxxFoo},
		"invalid":   {input: "bar", wantErr: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := NewXxx(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("NewXxx() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

### VO メソッド（プリミティブ / VO 返し）

```go
func TestXxx_Method(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		receiver Xxx
		arg      Yyy
		want     bool
	}{
		"case_true":  {receiver: XxxFoo, arg: YyyBar, want: true},
		"case_false": {receiver: XxxFoo, arg: YyyBaz, want: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.receiver.Method(tt.arg)
			if got != tt.want {
				t.Errorf("Method() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

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
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if diff := cmp.Diff(tt.want, got, xxxCmpOpts()...); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
```
