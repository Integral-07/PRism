package webhook_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Integral-07/prism/internal/handler/webhook"
)

const testSecret = "test-secret"

func sign(secret, body string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func applyMiddleware(secret string) http.Handler {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return webhook.VerifySignature(secret)(next)
}

func TestVerifySignature(t *testing.T) {
	body := `{"action":"opened"}`

	tests := []struct {
		name       string
		signature  string
		body       string
		wantStatus int
	}{
		{
			name:       "valid signature",
			signature:  sign(testSecret, body),
			body:       body,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing signature header",
			signature:  "",
			body:       body,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong signature",
			signature:  sign("wrong-secret", body),
			body:       body,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "signature without sha256= prefix",
			signature:  hex.EncodeToString([]byte("raw-hash")),
			body:       body,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "tampered body",
			signature:  sign(testSecret, body),
			body:       `{"action":"closed"}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "empty body with valid signature",
			signature:  sign(testSecret, ""),
			body:       "",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(tt.body))
			if tt.signature != "" {
				req.Header.Set("X-Hub-Signature-256", tt.signature)
			}
			rec := httptest.NewRecorder()

			applyMiddleware(testSecret).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}

func TestVerifySignature_bodyPreserved(t *testing.T) {
	body := `{"action":"opened"}`
	var capturedBody string

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		capturedBody = string(b)
		w.WriteHeader(http.StatusOK)
	})
	handler := webhook.VerifySignature(testSecret)(next)

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	req.Header.Set("X-Hub-Signature-256", sign(testSecret, body))
	handler.ServeHTTP(httptest.NewRecorder(), req)

	if capturedBody != body {
		t.Errorf("body preserved: want %q, got %q", body, capturedBody)
	}
}
