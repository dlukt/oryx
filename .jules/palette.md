## 2024-05-22 - Default Password Visibility
**Learning:** Users (and security best practices) expect password fields to be masked by default. Defaulting to plaintext, even with a toggle, creates immediate friction and potential security exposure.
**Action:** Always default password inputs to `type="password"`. Ensure tests explicitly check for this default state.
