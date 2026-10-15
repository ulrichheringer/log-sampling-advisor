# log-sampling-advisor

Estimate log volume reduction and discarded level/message signatures for a deterministic sampling policy.

## Install

Requires Go 1.23 or newer. Build with `go install github.com/ulrichheringer/log-sampling-advisor@latest`.

## Use

Pass JSON-lines logs with `-input logs.jsonl`; use `-strategy every` or `-strategy severity` and set `-every N`.

## Example

Run `go run . -input examples/simulated-events.jsonl -strategy severity -every 10`. The records are simulated; estimates reflect only the supplied sample.
