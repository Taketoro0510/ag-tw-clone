import { test, expect } from '@playwright/test';

test('login flow and verify timeline', async ({ page }) => {
  // Capture console logs to debug
  page.on('console', msg => console.log('BROWSER CONSOLE:', msg.text()));
  page.on('pageerror', err => console.log('BROWSER ERROR:', err.message));
  page.on('requestfailed', request =>
    console.log('REQUEST FAILED:', request.url(), request.failure()?.errorText)
  );

  // Navigate to login
  await page.goto('/login');
  
  // Wait a bit to ensure it settles
  await page.waitForTimeout(1000);

  // Click Sign in
  await page.click('button:has-text("Sign in with Google")');

  // Firebase Emulator uses a redirect or popup
  // Wait for the redirect to emulator URL
  await expect(page).toHaveURL(/.*localhost:9099\/emulator\/auth\/handler.*/);

  // In the emulator UI, click "Add new account" if it exists, otherwise just click the auto-generate button
  try {
    const addNewAccount = page.locator('text="Add new account"');
    await addNewAccount.waitFor({ state: 'visible', timeout: 5000 });
    await addNewAccount.click();
  } catch (e) {
    console.log("Add new account button not found or not needed");
  }

  // "Auto-generate user information"
  const autoGen = page.locator('text="Auto-generate user information"');
  await autoGen.waitFor({ state: 'visible', timeout: 5000 });
  await autoGen.click();

  // "Sign-in" button
  const signInButton = page.locator('button', { hasText: /Sign in/i }).last();
  await signInButton.waitFor({ state: 'visible', timeout: 5000 });
  await signInButton.click();

  // Wait for redirect back
  await page.waitForURL(/.*localhost:5173.*/, { timeout: 10000 });

  // Wait for timeline
  await expect(page.locator('text="Timeline"')).toBeVisible({ timeout: 10000 });
});
