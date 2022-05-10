# valorcheck

Linter to check that access to an optional value is guarded against the case where the value is not present.

## Installation

```bash
go install github.com/phelmkamp/valor/valorcheck@latest
```

## Usage

```bash
valorcheck [-flag] [package]

Flags:
  -V    print version and exit
  -c int
        display offending line with this many lines of context (default -1)
  -cpuprofile string
        write CPU profile to this file
  -debug string
        debug flags, any subset of "fpstv"
  -fix
        apply all suggested fixes
  -flags
        print analyzer flags in JSON
  -json
        emit JSON output
  -memprofile string
        write memory profile to this file
  -trace string
        write trace log to this file
```

## Output

```bash
/home/phelmkamp/documents/valor/valorcheck/testdata/main.go:15:2: call to MustOk not guarded by IsOk might panic
/home/phelmkamp/documents/valor/valorcheck/testdata/main.go:16:7: call to MustOk not guarded by IsOk might panic
/home/phelmkamp/documents/valor/valorcheck/testdata/main.go:17:2: result of Ok is not checked
```
