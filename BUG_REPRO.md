# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	skillsassessment/cmd/skills	[no test files]
ok  	skillsassessment/domain	0.007s
ok  	skillsassessment/reporting	0.005s
ok  	skillsassessment/repository	0.040s
ok  	skillsassessment/service	0.037s
ok  	skillsassessment/storage	0.059s
--- FAIL: TestWorkflowPublication (0.03s)
    publication_test.go:28: archived project must reject publish, got <nil>
FAIL
FAIL	skillsassessment/workflow	0.074s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/skills): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/skills): exit `0`
