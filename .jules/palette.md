## 2024-05-22 - Default Password Visibility
**Learning:** Users (and security best practices) expect password fields to be masked by default. Defaulting to plaintext, even with a toggle, creates immediate friction and potential security exposure.
**Action:** Always default password inputs to `type="password"`. Ensure tests explicitly check for this default state.

## 2024-05-22 - Semantic Buttons for Icons
**Learning:** The application frequently used `div role='button'` for icon-only actions (like copy to clipboard). This pattern lacks keyboard accessibility (tab focus, Enter/Space activation) and semantic meaning.
**Action:** Replaced these instances with a reusable `IconButton` component that wraps a semantic `<button>` with appropriate `aria-label` and `onClick` handlers, preserving the visual "ghost" style using Bootstrap classes.

## 2025-02-12 - Accessible File Inputs
**Learning:** `display: none` on file inputs removes them from the accessibility tree, making them inaccessible to keyboard users.
**Action:** Use a visually hidden pattern (e.g., `clip: rect(0 0 0 0)`) and track the input's focus state to apply a visible outline to the custom label.

## 2025-02-12 - Semantic Link Buttons
**Learning:** The application used `<a href="#!">` for actions that look like links but behave like buttons. This is an accessibility anti-pattern as it misleads screen readers and requires extra effort for keyboard support.
**Action:** Replace `href="#!"` anchors with `<Button variant="link">`. Use utility classes like `p-0` and `text-decoration-none` to match the table cell styling if necessary.

## 2025-02-26 - Async Loading State Handling
**Learning:** Relying on side-effects inside `new Promise` wrappers (without returning them) to manage loading states can cause `finally` blocks to execute prematurely, hiding loading indicators before the operation completes.
**Action:** Use standard Promise chaining (`.then().catch().finally()`) or `async/await` with `try/catch/finally` to ensure `setLoading(false)` always runs after the async operation effectively completes.
