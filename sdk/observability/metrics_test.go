package observability

import "testing"

func TestObserveWSMessageIgnoreNonPositive(t *testing.T) {
	ObserveWSMessage("test", 0)
	ObserveWSMessage("test", -1)
}

func TestWSConnectionCountersDoNotPanic(t *testing.T) {
	IncWSConnection()
	DecWSConnection()
	IncWSRejected()
	IncHTTPPanic()
	IncHTTPInflight()
	DecHTTPInflight()
	ObserveHTTPRequest("GET", "/metrics", 200, 0.001)
}
