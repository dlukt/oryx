# Sentinel's Journal

## 2024-05-23 - Command Argument Injection in FFmpeg Configurations

**Vulnerability:** Found an argument injection vulnerability in the `forward` and `transcode` modules. The `Server` configuration field was used to construct the output URL for `ffmpeg`. Since `exec.Command` passes arguments directly to the process, a malicious `Server` value starting with `-` (e.g., `-version` or `-f`) would be interpreted by `ffmpeg` as a flag rather than a URL.

**Learning:** When using `exec.Command` (or `exec.CommandContext`), arguments are safe from shell injection (globbing, pipes, etc.), but NOT safe from argument injection if the called program parses them as flags. Always validate that user-controlled inputs passed as arguments do not start with `-` unless intended.

**Prevention:** Added input validation to ensure the `Server` field does not start with `-` in both `platform/forward.go` and `platform/trancode.go` (via a helper in `platform/utils.go`). In the future, prefer to validate against a strict schema (e.g., `^rtmp://`).
