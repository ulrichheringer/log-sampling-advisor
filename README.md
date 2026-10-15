# Log Sampling Advisor

[![Go](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/) [![MIT License](https://img.shields.io/badge/License-MIT-007ec6.svg)](LICENSE)

See how much log volume a sampling rule would remove and which level/message combinations could disappear.

## Install

With Go 1.23 or newer, install it with `go install github.com/ulrichheringer/log-sampling-advisor@latest`.

## Use

Pass a JSON Lines file with `-input logs.jsonl`. Choose `-strategy every` or `-strategy severity`, then set the interval with `-every N`.

## Example

The example data is simulated. Run `go run . -input examples/simulated-events.jsonl -strategy severity -every 10`. Results are estimates based on the supplied events.
