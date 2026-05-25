package adapter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAwaitReady_ImmediateSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	deadline := make(chan time.Time, 1)
	tick := make(chan time.Time, 1)
	tick <- time.Now() // fire tick immediately

	err := awaitReady(http.DefaultClient, ts.URL, deadline, tick)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAwaitReady_RetryThenSuccess(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	deadline := make(chan time.Time) // never fires
	tick := make(chan time.Time, 10)
	// Pre-fill 3 ticks: first two will return errRetry, third will succeed
	tick <- time.Now()
	tick <- time.Now()
	tick <- time.Now()

	err := awaitReady(http.DefaultClient, ts.URL, deadline, tick)
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestAwaitReady_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	deadline := make(chan time.Time, 1)
	deadline <- time.Now() // fire deadline immediately

	tick := make(chan time.Time)

	err := awaitReady(http.DefaultClient, ts.URL, deadline, tick)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}
