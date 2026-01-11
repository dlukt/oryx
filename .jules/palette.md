## 2024-05-22 - Password Input Refactoring
**Learning:** Swapping components (like toggling between two `Form.Control`s for password visibility) can cause focus management issues and disrupt the accessibility tree.
**Action:** Always use a single component and toggle its `type` attribute (e.g., `type={plaintext ? "text" : "password"}`) instead of conditional rendering of different components. This preserves focus, selection state, and event listeners naturally.
