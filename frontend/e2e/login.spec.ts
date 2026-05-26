import { test, expect } from '@playwright/test';

test('login, post, comment, bookmark, and profile interactions', async ({ page }) => {
  // Capture console logs to debug
  page.on('console', msg => console.log('BROWSER CONSOLE:', msg.text()));
  page.on('pageerror', err => console.log('BROWSER ERROR:', err.message));
  page.on('requestfailed', request =>
    console.log('REQUEST FAILED:', request.url(), request.failure()?.errorText)
  );

  // Navigate to login
  await page.goto('/login');
  await page.waitForTimeout(1000);
  await page.click('button:has-text("Sign in with Google")');

  // Wait for redirect to Firebase emulator UI
  await expect(page).toHaveURL(/.*localhost:9099\/emulator\/auth\/handler.*/, { timeout: 15000 });
  await page.waitForLoadState('networkidle');

  // In the emulator UI, click "Add new account" if visible
  const addNewAccount = page.locator('text="Add new account"');
  const autoGen = page.locator('text="Auto-generate user information"');

  try {
    await Promise.race([
      addNewAccount.waitFor({ state: 'visible', timeout: 8000 }),
      autoGen.waitFor({ state: 'visible', timeout: 8000 })
    ]);
  } catch (e) {
    console.log("Timed out waiting for initial emulator screen, proceeding anyway");
  }

  if (await addNewAccount.isVisible()) {
    console.log("Clicking Add new account");
    await addNewAccount.click();
  }

  // Click "Auto-generate user information"
  await autoGen.waitFor({ state: 'visible', timeout: 10000 });
  await autoGen.click();

  // Click "Sign in"
  const signInButton = page.locator('button', { hasText: /Sign in/i }).last();
  await signInButton.waitFor({ state: 'visible', timeout: 10000 });
  await signInButton.click();

  // Wait for redirect back to the app
  await page.waitForURL(/.*localhost:5173.*/, { timeout: 10000 });

  // Wait for Timeline header to be visible
  await expect(page.locator('h5:has-text("Timeline")')).toBeVisible({ timeout: 10000 });

  // 1. Create a Post
  await page.click('button[aria-label="add"]');
  await expect(page.locator('text="Create Post"')).toBeVisible();

  const postBody = "E2E testing new premium features: follows and details!";
  await page.locator('textarea').first().fill(postBody);
  await page.click('button:has-text("Post")');

  // Wait for modal to close and post to appear
  await page.waitForTimeout(2000);

  // Verify post appeared
  const postCard = page.locator(`text="${postBody}"`).first();
  await expect(postCard).toBeVisible({ timeout: 5000 });

  // 2. Create a Comment
  const commentIcon = page.locator('button[aria-label="comment"]').first();
  await commentIcon.click();

  const commentText = "Premium comment works perfectly!";
  await page.fill('input[placeholder="Write a comment..."]', commentText);
  await page.click('button[type="submit"]');

  // Verify comment appeared
  await expect(page.locator(`text="${commentText}"`)).toBeVisible({ timeout: 5000 });

  // 3. Bookmark the Post
  const bookmarkIcon = page.locator('button[aria-label="bookmark"]').first();
  await bookmarkIcon.click();

  // Wait for bookmark mutation to complete
  await page.waitForTimeout(1000);

  // Verify Bookmarks page displays it
  await page.click('text="Bookmarks"');
  await expect(page.locator(`text="${postBody}"`).first()).toBeVisible({ timeout: 5000 });

  // 4. View Profile
  await page.click('text="Profile"');
  await expect(page.locator('span:has-text("Followers")').first()).toBeVisible({ timeout: 5000 });
  await expect(page.locator(`text="${postBody}"`).first()).toBeVisible({ timeout: 5000 });
});
