## 2024-05-23 - Login Form Accessibility
**Learning:** Combining password visibility toggle into a single input with dynamic `type` is more accessible than swapping components, as it preserves focus and DOM stability.
**Action:** Use `type={show ? 'text' : 'password'}` instead of conditional rendering for password inputs.
