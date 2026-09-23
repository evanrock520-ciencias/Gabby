import { Socket } from "k6/x/tcp";

export const options = {
  vus: 512
}

/**
 * Corre pruebas de estrés con m usuarios. El flujo es registrarse, mandar n mensajes públicos
 * y desconectarse. Asi, manda m*n mensajes.
 */
export default async function () {
  const host = __ENV.TCP_ECHO_HOST || "localhost";
  const port = __ENV.TCP_ECHO_PORT || "9090";
  const id = __VU
  const messages = 5


  const socket = new Socket();
  const closed = new Promise((resolve) => {
    socket.on("close", () => {
      console.log(`"Closed the connection ${id}"`);
      resolve();
    });
  });

  socket.on("error", (err) => {
    console.error("Error:", err);
  });

  socket.on("data", (data) => {
  });

  try {
    await socket.connect(port, host);

    await socket.write(`{"type":"IDENTIFY", "username":"u${id}"}\n`);

    const writes = []
    for (let i = 0; i < messages; i++) {
      writes.push(socket.write(`{"type":"PUBLIC_TEXT", "text":"u${id} send a message for the ${i} time"}\n`));
    }

    await Promise.all(writes);

    await socket.write(`{"type":"DISCONNECT"}\n`);
    socket.destroy();

    await closed;
  } catch (err) {
    console.error(`User ${id} caught error:`, err);
    socket.destroy();
  }
}
