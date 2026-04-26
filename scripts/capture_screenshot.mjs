import { chromium } from 'playwright';

const [url = 'http://127.0.0.1:8080', outPath = 'artifacts/users-screen.png'] = process.argv.slice(2);
const executablePath = process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH || undefined;

const browser = await chromium.launch({ headless: true, executablePath });
const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });

await page.goto(url, { waitUntil: 'networkidle', timeout: 30000 });
await page.screenshot({ path: outPath, fullPage: true });

await browser.close();
console.log(`saved screenshot: ${outPath}`);
