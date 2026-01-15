from playwright.sync_api import sync_playwright

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    page = browser.new_page()

    # Mock the secret query
    page.route("**/terraform/v1/hooks/srs/secret/query", lambda route: route.fulfill(
        status=200,
        body='{"code": 0, "data": {"publish": "test-secret"}}',
        headers={"Content-Type": "application/json"}
    ))

    # Mock token
    page.route("**/terraform/v1/mgmt/token", lambda route: route.fulfill(
        status=200,
        body='{"code": 0, "data": {"token": "test-token"}}'
    ))

    # Mock versions and check
    page.route("**/terraform/v1/mgmt/versions", lambda route: route.fulfill(
        status=200,
        body='{"code": 0, "data": {"version": "1.0.0"}}'
    ))
    page.route("**/terraform/v1/mgmt/check", lambda route: route.fulfill(
        status=200,
        body='{"code": 0, "data": {"upgrading": false}}'
    ))

    # Mock envs
    page.route("**/terraform/v1/mgmt/envs", lambda route: route.fulfill(
        status=200,
        body='{"code": 0, "data": {"platformDocker": true, "candidate": false}}'
    ))

    # Mock init - MUST return init: true
    page.route("**/terraform/v1/mgmt/init", lambda route: route.fulfill(
        status=200,
        body='{"code": 0, "data": {"init": true}}'
    ))

    # Mock status/user info
    page.route("**/terraform/v1/mgmt/status", lambda route: route.fulfill(
        status=200,
        body='{"code": 0, "data": {}}'
    ))

    # Mock beian/query (Footer)
    page.route("**/terraform/v1/mgmt/beian/query", lambda route: route.fulfill(
        status=200,
        body='{"code": 0, "data": {"icp": "", "title": ""}}'
    ))

    page.route("**/terraform/v1/mgmt/bilibili", lambda route: route.fulfill(
        status=200,
        body='{"code": 0, "data": {"title": "Tutorial", "desc": "Desc", "stat": {"view": 100, "like": 10, "share": 5}}}'
    ))

    # Set local storage for token
    page.add_init_script("""
        localStorage.setItem('SRS_TERRAFORM_TOKEN', JSON.stringify({token: 'test-token'}));
    """)

    # We need to go to a valid URL.
    # App.js redirects to /:locale/routers-login
    # If we go to /routers-scenario, it might redirect if logic says so.
    # But let's try direct navigation.

    url = "http://localhost:3000/mgmt/en/routers-scenario?tab=live"
    print(f"Navigating to {url}")
    try:
        page.goto(url, wait_until="networkidle")
    except Exception as e:
        print(f"Error navigating: {e}")
        # Try without /mgmt
        page.goto("http://localhost:3000/en/routers-scenario?tab=live", wait_until="networkidle")

    # Wait for the copy button
    try:
        page.wait_for_selector('button[aria-label="Copy"]', timeout=10000)
        print("Found button with aria-label='Copy'")
    except:
        print("Could not find button with aria-label='Copy'")
        # Dump screenshot to check where we are
        page.screenshot(path="verification/verification_debug.png")
        print("Saved debug screenshot")

    # Take screenshot
    page.screenshot(path="verification/verification.png")
    print("Screenshot saved to verification/verification.png")

    browser.close()

if __name__ == "__main__":
    with sync_playwright() as playwright:
        run(playwright)
