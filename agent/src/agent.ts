import { chromium } from 'playwright';

async function createAgent() {
  const browser = await chromium.connectOverCDP('http://127.0.0.1:9222');
  const ctx = browser.contexts()[0];
  return ctx
  // return page;
}

export { createAgent };
