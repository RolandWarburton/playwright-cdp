import { join } from 'path';
import { statSync } from 'fs';
import { createSocketServer, ISocketEvents } from './server.js';
import mitt from 'mitt';
import { createControllerV1 } from './controller.js';

function folderExists(path: string): boolean {
  try {
    const stats = statSync(path);
    return stats.isDirectory();
  } catch (error) {
    return false;
  }
}

function main() {
  // get the user id
  const userID = process.getuid;
  if (!userID) {
    process.exit(1);
  }

  // get the sock path
  const sockPath = `/var/run/user/${userID()}`;

  // check the required path for the socket exists
  if (!folderExists(sockPath)) {
    process.exit(1);
  }

  const emitter = mitt<ISocketEvents>();
  createSocketServer(join(sockPath, 'playwright-agent.socket'), emitter);
  createControllerV1(emitter);
}

main();
