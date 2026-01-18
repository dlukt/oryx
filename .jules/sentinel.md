# Sentinel's Journal

## 2026-01-18 - Path Traversal in OCR Service

**Vulnerability:** The `/terraform/v1/ai/ocr/image/` endpoint allowed path traversal because it used user-supplied filenames from the URL to construct file paths without sanitizing directory separators. An attacker could use `../` to access files outside the intended `ocr` directory (e.g., `../../secret.jpg`).

**Learning:** `path.Join` cleans paths (resolving `..`) but does not enforce that the resulting path is within a specific root directory if the input contains enough `..` segments to traverse up. Always use `path.Base` to extract just the filename if the intent is to access a file in a specific flat directory, or explicitly validate that the resolved path starts with the expected root directory.

**Prevention:** Updated `platform/ocr.go` to use `path.Base` on the input filename before using it to construct the file path, ensuring it only accesses files within the `ocr` directory. Confirmed `platform/transcript.go` already correctly uses `path.Base`.

## 2026-01-13 - Incomplete Validation Fix in Dubbing Service (Bypass)

**Vulnerability:** The previous fix for SSRF in `platform/dubbing.go` applied `ValidateServerURL` only when the input contained `://`. Attackers could bypass this check by supplying a local file path (e.g., `/etc/passwd`) which does not contain `://` but is still accepted by `ffprobe` as a valid input, resulting in arbitrary file read.

**Learning:** Do not condition security validation on the presence of specific characters (like `://`) in the input. If an input type (like `stream`) implies a protocol restriction, enforce that restriction unconditionally. If the input format doesn't match the expectation (e.g., missing protocol), reject it immediately.

**Prevention:** Updated `platform/dubbing.go` to unconditionally call `ValidateServerURL` for all `FFprobeSourceTypeStream` inputs.

## 2026-01-12 - SSRF and Arbitrary File Read in Dubbing Service

**Vulnerability:** The `/terraform/v1/dubbing/source` endpoint allowed `FFprobeSourceTypeStream` inputs without protocol validation. This allowed an attacker to supply `file://` or other dangerous protocols (like `http://` for SSRF) to `ffprobe`, potentially leaking file metadata or accessing internal services.

**Learning:** When accepting "stream" URLs or any user-defined URLs that are passed to tools like `ffmpeg` or `ffprobe`, always validate the protocol against a strict allowlist. Assuming "stream" implies safe network protocols is dangerous.

**Prevention:** Applied `ValidateServerURL` to the `FFprobeSourceTypeStream` path in `platform/dubbing.go`. This enforces that only `rtmp://`, `rtmps://`, `srt://`, and `rtsp://` protocols are allowed, preventing `file://` access and SSRF.

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

## 2026-10-24 - SSRF and Arbitrary File Read in Virtual Live and Camera Live Services

**Vulnerability:** The `platform/virtual-live-stream.go` and `platform/camera-live-stream.go` modules allowed `FFprobeSourceTypeStream` inputs without protocol validation, similar to the `dubbing` service issue. This allowed attackers to supply `file://` or other dangerous protocols to `ffprobe`.

**Learning:** Vulnerabilities often repeat across similar modules. When fixing a vulnerability in one place (e.g., `dubbing.go`), always check for similar patterns in other parts of the codebase (`virtual-live-stream.go`, `camera-live-stream.go`).

**Prevention:** Applied `ValidateServerURL` to `virtual-live-stream.go` and `camera-live-stream.go` to strictly allowlist protocols, ensuring only `rtmp://`, `rtmps://`, `srt://`, and `rtsp://` are processed.

## 2026-10-25 - [Weak RNG in Authentication]
**Vulnerability:** Usage of `math/rand` for generating JWT nonces and bucket name suffixes. `math/rand` is not cryptographically secure and was likely unseeded, making nonces predictable.
**Learning:** Default unseeded `math/rand` in Go produces deterministic sequences. Security-sensitive values must always use `crypto/rand`.
**Prevention:** Use `crypto/rand` for any security-related randomness. Audit codebase for `math/rand` usage.

## 2026-01-17 - Arbitrary File Read in AI Talk Service

**Vulnerability:** The `/terraform/v1/ai-talk/stage/hello-voices/` endpoint served files from a directory using user-supplied filenames without proper validation. Although `path.Join` and `ServeMux` mitigate simple `../` traversal in some cases, the endpoint allowed accessing any file in the configuration directory (e.g., `nginx.conf`) by specifying its name, exposing sensitive configuration.

**Learning:** When serving files based on user input, never rely solely on path joining or framework cleaning. Always validate the filename against a strict allowlist of expected files, especially if the directory contains sensitive information mixed with public resources.

**Prevention:** Implemented a strict allowlist in `platform/ai-talk.go` to only serve `hello-chinese.aac` and `hello-english.aac`.
