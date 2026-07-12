## Summary

What does this change do, and why?

## Wire protocol changes

If this touches `client.go`, `events.go`, `socketio.go`, `rawdata.go`, or
`rest.go`: what part of `~/nrepos/nerimity.js` did you verify this against?
(file + line, or a short quote). Leave this section out entirely if the change
doesn't touch the wire protocol.

## Checklist

- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] `gofmt -l .` prints nothing
- [ ] `go test ./...` passes
- [ ] New parsing/logic code has unit tests
- [ ] `examples/*/main.go` still compile, if the public API changed
- [ ] Docs in `docs/` updated, if behavior or the public API changed
