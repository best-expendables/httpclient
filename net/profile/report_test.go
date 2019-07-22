package profile

import (
	"testing"
	"time"
)

var timestamp = time.Unix(0, 0)

func TestReport_Duration(t *testing.T) {
	timestamp := timestamp

	report := Report{
		ConnectStart: timestamp,
		ConnectDone:  timestamp.Add(time.Second),
	}

	if report.ConnectionTime() != time.Second {
		t.Error("Duration not equals")
	}

	if report.ConnectionTimeMs() != 1000 {
		t.Error("Duration in milliseconds not equals")
	}
}

func TestReport_DNSLookupDuration(t *testing.T) {
	report := Report{
		DNSLookupStart: timestamp,
		DNSLookupDone:  timestamp.Add(time.Second),
	}

	if report.DNSLookupTime() != time.Second {
		t.Error("Duration not equals")
	}

	if report.DNSLookupTimeMs() != 1000 {
		t.Error("Duration in milliseconds not equals")
	}
}

func TestReport_TLSHandshakeDuration(t *testing.T) {
	report := Report{
		TLSHandshakeStart: timestamp,
		TLSHandshakeDone:  timestamp.Add(time.Second),
	}

	if report.TLSHandshakeTime() != time.Second {
		t.Error("Duration not equals")
	}

	if report.TLSHandshakeTimeMs() != 1000 {
		t.Error("Duration in milliseconds not equals")
	}
}
