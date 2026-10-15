package main

import "testing"

func TestSeveritySamplingPreservesErrors(t *testing.T) {
	logs := []LogRecord{{Level: "info", Message: "a"}, {Level: "error", Message: "failure"}, {Level: "info", Message: "b"}, {Level: "info", Message: "c"}}
	report, err := Analyze(logs, "severity", 3)
	if err != nil {
		t.Fatal(err)
	}
	if report.Kept != 3 || report.DiscardedErrors != 0 {
		t.Fatalf("unexpected report: %#v", report)
	}
}

func TestEverySamplingCountsLostSignatures(t *testing.T) {
	logs := []LogRecord{{Level: "info", Message: "a"}, {Level: "error", Message: "b"}, {Level: "info", Message: "c"}, {Level: "info", Message: "a"}}
	report, err := Analyze(logs, "every", 3)
	if err != nil {
		t.Fatal(err)
	}
	if report.Kept != 2 || report.LostSignatures != 2 || report.DiscardedErrors != 1 {
		t.Fatalf("unexpected report: %#v", report)
	}
}
