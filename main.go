package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type LogRecord struct {
	Level    string `json:"level"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Raw      string `json:"-"`
}

type Report struct {
	Total            int
	Kept             int
	DiscardedErrors  int
	UniqueSignatures int
	LostSignatures   int
	ReductionPct     float64
	SignatureLossPct float64
}

func ReadLogs(r io.Reader) ([]LogRecord, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	var logs []LogRecord
	for line := 1; scanner.Scan(); line++ {
		raw := scanner.Text()
		var record LogRecord
		if err := json.Unmarshal([]byte(raw), &record); err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		if record.Level == "" {
			record.Level = record.Severity
		}
		record.Raw = raw
		logs = append(logs, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return logs, nil
}

func Analyze(logs []LogRecord, strategy string, every int) (Report, error) {
	if every < 1 {
		return Report{}, errors.New("every must be at least 1")
	}
	if strategy != "every" && strategy != "severity" {
		return Report{}, fmt.Errorf("unknown strategy %q", strategy)
	}
	all, keptSignatures := map[string]bool{}, map[string]bool{}
	report := Report{Total: len(logs)}
	for i, record := range logs {
		signature := strings.ToLower(record.Level) + "|" + strings.TrimSpace(record.Message)
		all[signature] = true
		keep := i%every == 0
		if strategy == "severity" && isError(record.Level) {
			keep = true
		}
		if keep {
			report.Kept++
			keptSignatures[signature] = true
		} else if isError(record.Level) {
			report.DiscardedErrors++
		}
	}
	report.UniqueSignatures = len(all)
	for signature := range all {
		if !keptSignatures[signature] {
			report.LostSignatures++
		}
	}
	if report.Total > 0 {
		report.ReductionPct = 100 * float64(report.Total-report.Kept) / float64(report.Total)
	}
	if report.UniqueSignatures > 0 {
		report.SignatureLossPct = 100 * float64(report.LostSignatures) / float64(report.UniqueSignatures)
	}
	return report, nil
}

func isError(level string) bool {
	switch strings.ToLower(level) {
	case "error", "fatal", "critical", "panic", "alert", "emergency":
		return true
	default:
		return false
	}
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("log-sampling-advisor", flag.ContinueOnError)
	flags.SetOutput(stderr)
	input := flags.String("input", "", "JSON-lines logs file; stdin when omitted")
	strategy := flags.String("strategy", "every", "every or severity (retain all error-level records)")
	every := flags.Int("every", 10, "retain one out of each N records")
	if err := flags.Parse(args); err != nil {
		return err
	}
	var reader io.Reader = stdin
	if *input != "" {
		file, err := os.Open(*input)
		if err != nil {
			return err
		}
		defer file.Close()
		reader = file
	}
	logs, err := ReadLogs(reader)
	if err != nil {
		return err
	}
	report, err := Analyze(logs, *strategy, *every)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(stdout, "Strategy: %s (1 in %d)\nRecords: %d -> %d (%.1f%% volume reduction)\nUnique level/message signatures: %d; absent after sampling: %d (%.1f%%)\nError records discarded: %d\n\nEstimate is based on this input sample and assumes it represents the production mix.\n", *strategy, *every, report.Total, report.Kept, report.ReductionPct, report.UniqueSignatures, report.LostSignatures, report.SignatureLossPct, report.DiscardedErrors)
	return err
}
