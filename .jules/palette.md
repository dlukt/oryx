## 2024-05-22 - Default Password Visibility
**Learning:** Users (and security best practices) expect password fields to be masked by default. Defaulting to plaintext, even with a toggle, creates immediate friction and potential security exposure.
**Action:** Always default password inputs to `type="password"`. Ensure tests explicitly check for this default state.

## 2024-05-22 - Semantic Buttons for Icons
**Learning:** The application frequently used `div role='button'` for icon-only actions (like copy to clipboard). This pattern lacks keyboard accessibility (tab focus, Enter/Space activation) and semantic meaning.
**Action:** Replaced these instances with a reusable `IconButton` component that wraps a semantic `<button>` with appropriate `aria-label` and `onClick` handlers, preserving the visual "ghost" style using Bootstrap classes.

## 2025-02-12 - Accessible File Inputs
**Learning:** `display: none` on file inputs removes them from the accessibility tree, making them inaccessible to keyboard users.
**Action:** Use a visually hidden pattern (e.g., `clip: rect(0 0 0 0)`) and track the input's focus state to apply a visible outline to the custom label.
