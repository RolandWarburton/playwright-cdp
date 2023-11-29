import { Emitter, Handler, WildcardHandler } from 'mitt';
import { ISocketEvents, ISocketAction } from './server.js';
import { createAgent } from './agent.js';
import { BrowserContext } from 'playwright';

function createControllerV1(emitter: Emitter<ISocketEvents>) {
  console.log('create controler');
  let page = undefined;
  let browser: BrowserContext | undefined;
  emitter.on('data', async (data) => {
    // CONNECT TO AGENT
    if (data.action === 'connect') {
      if (browser) {
        console.log('refusing to reconnect as page has already been connected');
        return;
      }
      try {
        const ctx = await createAgent();
        console.log('agent created');
        browser = ctx;
      } catch (err) {
        console.log(err);
      }
    }

    // EXAMPLE ACTION
    if (data.action === 'action-example') {
      if (!browser) {
        console.log('no browser agent created');
        return;
      }
      try {
        page = await browser.newPage();
        page.goto('https://google.com')
      } catch (err) {
        console.log(err);
      }
    }
  });

  emitter.on('end', (message) => {
    // console.log(message)
  });
}

export { createControllerV1 };
