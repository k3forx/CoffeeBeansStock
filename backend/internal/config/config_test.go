package config

import (
	"strings"
	"testing"
)

func TestValidateJWTSecret(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		secret     string
		wantErrMsg string
	}{
		"empty secret": {
			secret:     "",
			wantErrMsg: "required",
		},
		"too short": {
			secret:     "Abc123!",
			wantErrMsg: "at least 32 characters",
		},
		"exactly 31 characters": {
			secret:     "aB3!aB3!aB3!aB3!aB3!aB3!aB3!aB3",
			wantErrMsg: "at least 32 characters",
		},
		"blocked secret": {
			secret:     "secret",
			wantErrMsg: "at least 32 characters",
		},
		"blocked secret case insensitive": {
			secret:     "JWT-SECRET",
			wantErrMsg: "at least 32 characters",
		},
		"blocked secret matching blocklist": {
			secret:     "Dev-Secret-Change-In-Production!",
			wantErrMsg: "known weak secret",
		},
		"low entropy - only lowercase": {
			secret:     strings.Repeat("a", 32),
			wantErrMsg: "at least 3 character types",
		},
		"low entropy - only digits": {
			secret:     strings.Repeat("1", 32),
			wantErrMsg: "at least 3 character types",
		},
		"low entropy - two types only": {
			secret:     strings.Repeat("aA", 16),
			wantErrMsg: "at least 3 character types",
		},
		"valid secret with 3 types": {
			secret:     "abcdefghijklmnopABCDEFGHIJK12345",
			wantErrMsg: "",
		},
		"valid base64 secret": {
			secret:     "K7gNU3sdo+OL0wNhqoVWhr3g6s1xYv72ol/pe/Unols=",
			wantErrMsg: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := validateJWTSecret(tt.secret)
			if tt.wantErrMsg == "" {
				if err != nil {
					t.Errorf("validateJWTSecret(%q) returned unexpected error: %v", tt.secret, err)
				}
				return
			}
			if err == nil {
				t.Errorf("validateJWTSecret(%q) expected error containing %q, got nil", tt.secret, tt.wantErrMsg)
				return
			}
			if !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("validateJWTSecret(%q) error = %q, want containing %q", tt.secret, err.Error(), tt.wantErrMsg)
			}
		})
	}
}
