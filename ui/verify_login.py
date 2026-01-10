from playwright.sync_api import sync_playwright

def verify_login_page():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # Navigate to login page
        try:
            print("Navigating to login page...")
            page.goto("http://localhost:3001/mgmt/")

            # The app likely redirects to login or shows login on home if not auth
            print("Waiting for page load...")
            page.wait_for_timeout(5000)

            print("Taking debug screenshot...")
            page.screenshot(path="/home/jules/verification/debug_page.png")

            # Wait for the password input to be visible
            print("Waiting for password input...")
            page.wait_for_selector("input[type='password']", timeout=10000)

            # 1. Verify "Show password" functionality
            # Type something
            page.fill("input[type='password']", "srs123")

            # Check checkbox
            page.click("input[type='checkbox']")

            # Verify input type changes to text
            page.wait_for_selector("input[type='text']")

            # Take screenshot of show password state
            page.screenshot(path="/home/jules/verification/login_show_password.png")

            # Uncheck
            page.click("input[type='checkbox']")

            # Verify input type changes back to password
            page.wait_for_selector("input[type='password']")

            # 2. Verify spinner inside button (simulate loading)

            # Let's take a screenshot of the clean state
            page.screenshot(path="/home/jules/verification/login_clean.png")

            # Fill password and press Enter to verify form submission trigger
            page.press("input[type='password']", "Enter")

            # Wait a bit
            page.wait_for_timeout(1000)

            # Take screenshot of loading state (if caught)
            page.screenshot(path="/home/jules/verification/login_submit.png")

        except Exception as e:
            print(f"Error: {e}")
            page.screenshot(path="/home/jules/verification/error.png")
        finally:
            browser.close()

if __name__ == "__main__":
    verify_login_page()
