## 2024-05-23 - React Polling and State Updates
**Learning:** Polling intervals that blindly update state with new object references (even if content is identical) cause unnecessary re-renders. `useState` only bails out if the *reference* is the same.
**Action:** When polling for data that changes infrequently, use functional state updates and compare deep equality (e.g., `JSON.stringify` for simple objects) before returning the new state. If identical, return `prevState` to skip the re-render.

## 2025-01-25 - Efficient Global State Management
**Learning:** Using `setInterval` to poll for changes in global state (like language/locale) that rarely changes is inefficient, especially when used in many components.
**Action:** Replace polling with a Pub/Sub mechanism (Observer pattern) for global state that is read frequently but updated rarely. Add a `subscribe` method to the state manager and use `useEffect` to listen for changes in consuming components.
