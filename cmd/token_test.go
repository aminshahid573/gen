package cmd

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestGenerateToken_IsBase64URL(t *testing.T) {
	token, err := generateToken(32, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := base64.RawURLEncoding.DecodeString(token); err != nil {
		t.Errorf("token %q is not valid base64url: %v", token, err)
	}
}

func TestGenerateToken_Prefix(t *testing.T) {
	token, err := generateToken(32, "tok_", false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(token, "tok_") {
		t.Errorf("token %q missing prefix tok_", token)
	}
}

func TestGenerateToken_Bearer(t *testing.T) {
	token, err := generateToken(32, "", true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(token, "Bearer ") {
		t.Errorf("token %q missing Bearer prefix", token)
	}
}

func TestGenerateToken_PrefixAndBearer(t *testing.T) {
	token, err := generateToken(32, "sess_", true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(token, "Bearer sess_") {
		t.Errorf("token %q missing 'Bearer sess_' prefix", token)
	}
}

func TestGenerateToken_Uniqueness(t *testing.T) {
	a, _ := generateToken(32, "", false)
	b, _ := generateToken(32, "", false)
	if a == b {
		t.Error("two tokens should not be identical")
	}
}

func TestGenerateToken_LongerBytes(t *testing.T) {
	short, _ := generateToken(16, "", false)
	long, _  := generateToken(64, "", false)
	if len(long) <= len(short) {
		t.Errorf("64-byte token (%d) should be longer than 16-byte token (%d)", len(long), len(short))
	}
}
