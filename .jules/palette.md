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

## 2026-01-21 - IconButton Flexibility
**Learning:** The `IconButton` component was too rigid, lacking `className` support, which prevented its use in existing layouts that relied on utility classes for positioning or cursor styles (e.g., `ai-dubbing-command`).
**Action:** Update reusable components to accept and merge `className` props to enable wider adoption without style regression.

## 2026-02-05 - Actionable Identifiers
**Learning:** Identifiers like UUIDs in list views are often used for API debugging or configuration. Making them purely navigational links frustrates users who need to copy them.
**Action:** Always provide a dedicated 'Copy' action next to long identifiers or non-selectable link text in data tables.

## 2026-02-26 - Standardizing Copy Actions
**Learning:** Inconsistent implementation of "Copy" actions (text links vs custom divs vs buttons) leads to poor UX and accessibility. Icon-only buttons in data tables save space and provide a recognizable affordance.
**Action:** Consistently use the `CopyButton` component for all copy-to-clipboard interactions, especially in data tables and status displays.
