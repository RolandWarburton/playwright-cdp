import { chmodSync } from 'fs';
import net from 'net';
import { Emitter } from 'mitt';

type ISocketAction = 'connect' | 'eval' | 'action-example';
type ISocketEvents = {
  data: { action?: ISocketAction; [key: string]: any };
  end: string;
};

function createSocketServer(socketPath: string, emitter: Emitter<ISocketEvents>) {
  const server = net.createServer((socket) => {
    socket.on('data', (socketData) => {
      try {
        const data: object = JSON.parse(socketData.toString());
        if (!data) {
          throw new Error('failed to marshal data');
        }
        emitter.emit('data', data);
      } catch (err) {
        console.error(err);
      }
    });

    socket.on('end', () => {
      emitter.emit('end', 'closed');
    });
  });

  server.listen(socketPath, () => {
    console.log('Unix socket server is listening');
  });

  // Set the permissions of the socket file
  chmodSync(socketPath, '777');

  // Trap Ctrl+C (SIGINT) to shut down the server gracefully
  process.on('SIGINT', () => {
    console.log('Shutting down server...');
    server.close(() => {
      console.log('Server has been shut down');
      process.exit(0);
    });
  });
}

export { createSocketServer, ISocketEvents, ISocketAction };
