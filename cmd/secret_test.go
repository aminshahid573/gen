package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestGenerateSecret_Hex(t *testing.T) {
	out, err := generateSecret(32, "hex")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 64 { // 32 bytes → 64 hex chars
		t.Errorf("hex len = %d, want 64", len(out))
	}
	if _, err := hex.DecodeString(out); err != nil {
		t.Errorf("output is not valid hex: %v", err)
	}
}

func TestGenerateSecret_Base64(t *testing.T) {
	out, err := generateSecret(32, "base64")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := base64.StdEncoding.DecodeString(out); err != nil {
		t.Errorf("output is not valid base64: %v", err)
	}
}

func TestGenerateSecret_Base64URL(t *testing.T) {
	out, err := generateSecret(32, "base64url")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := base64.RawURLEncoding.DecodeString(out); err != nil {
		t.Errorf("output is not valid base64url: %v", err)
	}
}

func TestGenerateSecret_InvalidEncoding(t *testing.T) {
	_, err := generateSecret(32, "invalid")
	if err == nil {
		t.Fatal("expected error for invalid encoding")
	}
}

func TestGenerateSecret_Uniqueness(t *testing.T) {
	a, _ := generateSecret(32, "hex")
	b, _ := generateSecret(32, "hex")
	if a == b {
		t.Error("two secrets should not be identical")
	}
}

func TestResolveSecretDefaults(t *testing.T) {
	tests := []struct {
		secretType      string
		wantBytes       int
		wantEncoding    string
	}{
		{"generic", 32, "hex"},
		{"jwt", 64, "base64url"},
		{"hmac", 32, "hex"},
		{"api", 32, "base64url"},
	}

	for _, tt := range tests {
		t.Run(tt.secretType, func(t *testing.T) {
			gotBytes, gotEnc := resolveSecretDefaults(tt.secretType, 0, "", false, false)
			if gotBytes != tt.wantBytes {
				t.Errorf("bytes = %d, want %d", gotBytes, tt.wantBytes)
			}
			if gotEnc != tt.wantEncoding {
				t.Errorf("encoding = %q, want %q", gotEnc, tt.wantEncoding)
			}
		})
	}
}

func TestResolveSecretDefaults_ExplicitFlagsOverride(t *testing.T) {
	// explicit --bytes should override type default
	gotBytes, _ := resolveSecretDefaults("jwt", 16, "", true, false)
	if gotBytes != 16 {
		t.Errorf("explicit bytes not respected: got %d", gotBytes)
	}

	// explicit --encoding should override type default
	_, gotEnc := resolveSecretDefaults("jwt", 0, "hex", false, true)
	if gotEnc != "hex" {
		t.Errorf("explicit encoding not respected: got %q", gotEnc)
	}
}
