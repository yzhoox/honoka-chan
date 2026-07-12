package honokautils

import (
	"net/http"
	"testing"
)

func TestInternalErrorContentIsSanitized(t *testing.T) {
	content := NewInternalErrorContent()

	if content.Exception != "InternalServerError" {
		t.Fatalf("unexpected exception: %q", content.Exception)
	}
	if content.Message != http.StatusText(http.StatusInternalServerError) {
		t.Fatalf("unexpected message: %q", content.Message)
	}
}

func TestPanicContentIsSanitized(t *testing.T) {
	content := NewPanicContent()

	if content != NewInternalErrorContent() {
		t.Fatalf("panic content exposes a distinct internal detail: %#v", content)
	}
}

func TestIsMainLoginEndpoint(t *testing.T) {
	for _, path := range []string{"/main.php/login/authkey", "/main.php/login/login"} {
		if !IsMainLoginEndpoint(path) {
			t.Fatalf("expected %q to be a main.php login endpoint", path)
		}
	}

	for _, path := range []string{"/login/authkey", "/main.php/user/userInfo"} {
		if IsMainLoginEndpoint(path) {
			t.Fatalf("expected %q not to be a main.php login endpoint", path)
		}
	}
}
