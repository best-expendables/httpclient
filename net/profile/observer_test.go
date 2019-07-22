package profile

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Millisecond * 5)
	w.WriteHeader(http.StatusOK)
}))

func TestObserve(t *testing.T) {
	c := &http.Client{}

	report := ReportFromResponse(sendRequest(c))

	if report == nil {
		t.Fatal("Report can not be nil")
	}

	if report.ConnectStart.IsZero() || report.ConnectDone.IsZero() {
		t.Errorf("ConnectionStart [%s] or ConnectionDone [%s] is zero",
			report.ConnectStart, report.ConnectDone)
	}
	t.Logf("ConnectTime connect - %s", report.ConnectionTime())

	// Reused connection
	report = ReportFromResponse(sendRequest(c))

	if report.Reused == false {
		t.Error("Connection should be reused")
	}

	if !report.ConnectStart.IsZero() && !report.ConnectDone.IsZero() {
		t.Error("ConnecttStart and ConnectDone should be zero")
	}
}

func BenchmarkObserve10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r, _ := http.NewRequest(http.MethodGet, "", nil)
		Observe(r)
	}
}

func sendRequest(c *http.Client) *http.Response {
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	res, _ := c.Do(Observe(req))
	return res
}
