package network

import (
	"bufio"
	"client/internal/protocol"
	"log"
	"net"
)

// ConnectionManager maneja la comunicación del cliente
// con el servidor.
type ConnectionManager struct {
	conn     net.Conn
	reader   *bufio.Reader
	incoming chan protocol.ServerMessage
}

// Dial establece la conexión con el servidor.
func (manager *ConnectionManager) Dial(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("The connection wasn't succesfull")
	}
	manager.conn = conn
	manager.reader = bufio.NewReader(conn)
	manager.incoming = make(chan protocol.ServerMessage)
	log.Printf("The connection was succesfull with %s", conn.LocalAddr())
}

// Send manda mensajes del cliente al servidor.
func (manager *ConnectionManager) Send(msg protocol.ClientMessage) {
	line, err := protocol.Serialize(msg)

	if err != nil {
		log.Printf("The message has an invalid format")
	}

	manager.conn.Write([]byte(line + "\n"))
}

// Listen escucha los mensajes del servidor al cliente
// y los manda a través de un canal.
func (manager *ConnectionManager) Listen() {
	for {
		line, err := manager.reader.ReadString('\n')
		if err != nil {
			log.Printf("Connection closed: %v", err)
			return
		}

		msg, err := protocol.Deserialize(line)

		if err != nil {
			log.Printf("Invalid message: %s", line)
			return
		}

		manager.incoming <- msg
	}
}

// Messages otorga acceso al canal de la conexión.
func (manager *ConnectionManager) Messages() <-chan protocol.ServerMessage {
	return manager.incoming
}

// Close cierra la conexión.
func (m *ConnectionManager) Close() error {
	return m.conn.Close()
}
