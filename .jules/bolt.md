## 2024-05-22 - Combined Async State Updates
**Learning:** React 17 does not batch state updates in async functions (like `Promise.then`). Independent `useEffect` hooks fetching data and setting state independently cause multiple renders.
**Action:** Use `Promise.all` to fetch multiple resources and update state in a single object to trigger only one re-render. This also simplifies error handling for authentication failures (avoiding double redirects/alerts).
