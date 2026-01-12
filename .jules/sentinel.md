# Sentinel's Journal

## 2026-01-12 - Arbitrary File Write via FFmpeg Protocol Abuse

**Vulnerability:** The previous fix for command injection (blocking `-`) was insufficient. Attackers could still use FFmpeg protocols like `file://` in the `Server` configuration to overwrite arbitrary files on the system when combined with the `Secret` field to form the output URL.

**Learning:** Blacklisting dangerous characters (like `-`) is often insufficient. When constructing URLs or paths passed to powerful tools like FFmpeg, allowlisting protocols (e.g., only allowing `rtmp://`, `srt://`) is mandatory to prevent protocol smuggling or abuse (like `file://`, `http://` for SSRF, etc.).

**Prevention:** Updated `ValidateServerURL` to strictly allowlist only `rtmp://`, `rtmps://`, `srt://`, and `rtsp://` protocols.

## 2024-05-23 - Command Argument Injection in FFmpeg Configurations

**Vulnerability:** Found an argument injection vulnerability in the `forward` and `transcode` modules. The `Server` configuration field was used to construct the output URL for `ffmpeg`. Since `exec.Command` passes arguments directly to the process, a malicious `Server` value starting with `-` (e.g., `-version` or `-f`) would be interpreted by `ffmpeg` as a flag rather than a URL.

**Learning:** When using `exec.Command` (or `exec.CommandContext`), arguments are safe from shell injection (globbing, pipes, etc.), but NOT safe from argument injection if the called program parses them as flags. Always validate that user-controlled inputs passed as arguments do not start with `-` unless intended.

**Prevention:** Added input validation to ensure the `Server` field does not start with `-` in both `platform/forward.go` and `platform/trancode.go` (via a helper in `platform/utils.go`). In the future, prefer to validate against a strict schema (e.g., `^rtmp://`).

## 2024-03-25 - [JWT Algorithm Confusion Defense]
**Vulnerability:** The `Authenticate` function in `platform/utils.go` used `jwt.Parse` without verifying the signing method (`token.Method`).
**Learning:** This is a classic JWT vulnerability (CWE-327). Even if we use a symmetric key (`apiSecret`) and expect `HS256`, failure to explicitly check `token.Method` in the keyfunc allows attackers to potentially use the `none` algorithm (if supported/enabled by the library or configuration) or perform algorithm confusion attacks (e.g. changing RS256 to HS256 if we were using RSA keys). Although this project uses `HS256`, the lack of check is a violation of secure coding practices for JWTs and could be exploited if the library behavior changes or if future refactoring introduces asymmetric keys.
**Prevention:** Always verify `token.Method` inside the `jwt.Parse` callback function. Ensure it matches the expected signing method (e.g., `*jwt.SigningMethodHMAC`).
