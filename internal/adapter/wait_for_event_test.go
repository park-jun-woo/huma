package adapter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWaitForEvent_DeadlineFires(t *testing.T) {
	deadline := make(chan time.Time, 1)
	deadline <- time.Now()
	tick := make(chan time.Time)

	err := waitForEvent(http.DefaultClient, "http://localhost:99999", deadline, tick)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestWaitForEvent_TickSucceeds(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	deadline := make(chan time.Time)
	tick := make(chan time.Time, 1)
	tick <- time.Now()

	err := waitForEvent(http.DefaultClient, ts.URL, deadline, tick)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWaitForEvent_TickRetry(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	deadline := make(chan time.Time)
	tick := make(chan time.Time, 1)
	tick <- time.Now()

	err := waitForEvent(http.DefaultClient, ts.URL, deadline, tick)
	if err != errRetry {
		t.Fatalf("expected errRetry, got %v", err)
	}
}
