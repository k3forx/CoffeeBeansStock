package apperrors

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var errBase = errors.New("base error")

func TestWrap(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   error
		wantNil bool
	}{
		"nilを渡すとnilを返す": {
			input:   nil,
			wantNil: true,
		},
		"エラーをラップするとwithStackErrorが返される": {
			input:   errBase,
			wantNil: false,
		},
		"すでにwithStackErrorを含むエラーは二重ラップしない": {
			input:   Wrap(errBase),
			wantNil: false,
		},
		"fmt.Errorfでラップされた中にwithStackErrorがある場合も二重ラップしない": {
			input:   fmt.Errorf("outer: %w", Wrap(errBase)),
			wantNil: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := Wrap(tt.input)
			if tt.wantNil {
				if got != nil {
					t.Errorf("Wrap() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Errorf("Wrap() = nil, want non-nil")
				return
			}
			_, ok := GetStackTrace(got)
			if !ok {
				t.Errorf("Wrap() returned error without stack trace")
			}
		})
	}
}

func TestWrap_ErrorsIs(t *testing.T) {
	t.Parallel()
	wrapped := Wrap(errBase)
	if !errors.Is(wrapped, errBase) {
		t.Errorf("errors.Is(Wrap(errBase), errBase) = false, want true")
	}
}

func TestWrap_DoubleWrapReturnsSameError(t *testing.T) {
	t.Parallel()
	first := Wrap(errBase)
	second := Wrap(first)
	if first != second {
		t.Errorf("二重ラップが防止されていない: Wrap(Wrap(err)) should return same error")
	}
}

func TestWrapf(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   error
		format  string
		args    []any
		wantNil bool
	}{
		"nilを渡すとnilを返す": {
			input:   nil,
			format:  "operation failed",
			wantNil: true,
		},
		"フォーマットメッセージ付きでラップできる": {
			input:  errBase,
			format: "operation %s failed",
			args:   []any{"save"},
		},
		"すでにwithStackErrorを含むエラーは二重ラップしない": {
			input:  Wrap(errBase),
			format: "outer",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := Wrapf(tt.input, tt.format, tt.args...)
			if tt.wantNil {
				if got != nil {
					t.Errorf("Wrapf() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Errorf("Wrapf() = nil, want non-nil")
				return
			}
			_, ok := GetStackTrace(got)
			if !ok {
				t.Errorf("Wrapf() returned error without stack trace")
			}
		})
	}
}

func TestWrapf_ErrorsIs(t *testing.T) {
	t.Parallel()
	wrapped := Wrapf(errBase, "operation failed")
	if !errors.Is(wrapped, errBase) {
		t.Errorf("errors.Is(Wrapf(errBase, ...), errBase) = false, want true")
	}
}

func TestWrapf_MessageIncluded(t *testing.T) {
	t.Parallel()
	wrapped := Wrapf(errBase, "operation %s failed", "save")
	want := "operation save failed: base error"
	if got := wrapped.Error(); got != want {
		t.Errorf("Wrapf().Error() = %q, want %q", got, want)
	}
}

func TestGetStackTrace(t *testing.T) {
	t.Parallel()

	t.Run("スタックトレースが存在する場合はtrueを返す", func(t *testing.T) {
		t.Parallel()
		err := Wrap(errBase)
		stack, ok := GetStackTrace(err)
		if !ok {
			t.Errorf("GetStackTrace() ok = false, want true")
		}
		if len(stack) == 0 {
			t.Errorf("GetStackTrace() returned empty stack trace")
		}
	})

	t.Run("スタックトレースが存在しない場合はfalseを返す", func(t *testing.T) {
		t.Parallel()
		_, ok := GetStackTrace(errBase)
		if ok {
			t.Errorf("GetStackTrace() ok = true, want false")
		}
	})

	t.Run("スタックトレースに呼び出し元のフレームが含まれる", func(t *testing.T) {
		t.Parallel()
		err := Wrap(errBase)
		stack, _ := GetStackTrace(err)
		if len(stack) == 0 {
			t.Fatal("stack trace is empty")
		}
		found := false
		for _, f := range stack {
			if strings.Contains(f.Function, "TestGetStackTrace") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("stack trace does not contain TestGetStackTrace frame: %+v", stack)
		}
	})

	t.Run("フレーム情報が正しく設定される", func(t *testing.T) {
		t.Parallel()
		err := Wrap(errBase)
		stack, _ := GetStackTrace(err)
		if len(stack) == 0 {
			t.Fatal("stack trace is empty")
		}
		first := stack[0]
		if diff := cmp.Diff(false, first.File == ""); diff != "" {
			t.Errorf("Frame.File should not be empty: %s", diff)
		}
		if first.Line == 0 {
			t.Errorf("Frame.Line should not be 0")
		}
		if first.Function == "" {
			t.Errorf("Frame.Function should not be empty")
		}
	})
}
