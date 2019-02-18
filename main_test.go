package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestGameScenario(t *testing.T) {
	assert := assert.New(t)
	go withTimeout(20*time.Second, interruptableMain)

	// Connect as a game client
	conn, err := net.DialTimeout("tcp", "localhost:52000", 5*time.Second)
	assert.NoError(err, "client cannot connect")
	client := NewTestClient(conn, assert)

	// Start game
	client.send("START john")

	// Verify zombie walks & take a miscalculated shot
	x, y := client.receiveWalk()
	client.shoot(x, y+10)

	// Shot is missed
	missMsg := client.receiveUntil("MISS")
	assert.Equal("MISS", missMsg)

	// Take a precision shot
	x, y = client.receiveWalk()
	client.shoot(x, y)

	// Shot is successful
	boomMsg := client.receiveUntil("BOOM")
	assert.Equal("BOOM john night-king", boomMsg)

	// Game over (victory)
	victoryMsg := client.receive()
	assert.Equal("VICTORY john", victoryMsg)
}

type testClient struct {
	conn     net.Conn
	assert   *assert.Assertions
	incoming chan string
	outgoing chan string
}

func NewTestClient(conn net.Conn, assert *assert.Assertions) *testClient {
	incoming := make(chan string, 5)
	outgoing := make(chan string, 5)

	client := &testClient{conn: conn, assert: assert, incoming: incoming, outgoing: outgoing}
	go client.sendWorker()
	go client.receiveWorker()

	return client
}

func (client *testClient) shoot(x, y int) {
	client.outgoing <- fmt.Sprintf("SHOOT %d %d", x, y)
}

func (client *testClient) sendWorker() {
	defer client.conn.Close()
	for {
		select {
		case msg, ok := <-client.outgoing:
			if !ok {
				return
			}
			_, err2 := io.WriteString(client.conn, msg+"\r\n")
			if err2 != nil {
				client.assert.FailNow("Failed to send message")
			}
		}
	}
}

func (client *testClient) send(msg string) {
	client.outgoing <- msg
}

func (client *testClient) receiveWorker() {
	defer client.conn.Close()
	for {
		buf := make([]byte, 4096)
		n, err := client.conn.Read(buf)
		if err != nil {
			client.assert.FailNow("Failed to read message")
		}

		received := n > 0
		if !received {
			continue
		}

		line := strings.TrimRight(string(buf[:n]), "\r\n")
		messages := strings.Split(line, "\r\n")
		for _, msg := range messages {
			client.incoming <- msg
		}
	}
}

func (client *testClient) receive() string {
	for {
		select {
		case msg := <-client.incoming:
			return msg
		case <-time.After(5 * time.Second):
			client.assert.FailNow("Receive timed out")
		}
	}
}

func (client *testClient) receiveUntil(command string) string {
	for i := 0; i < 5; i++ {
		msg := client.receive()
		receivedCommand := strings.Split(msg, " ")[0]
		if command == receivedCommand {
			return msg
		}
	}
	return ""
}

func (client *testClient) receiveWalk() (x, y int) {
	msg := client.receive()
	_, err := fmt.Sscanf(msg, "WALK %d %d", &x, &y)
	client.assert.NoError(err)
	return
}

func withTimeout(duration time.Duration, fn func(kill chan bool, done chan bool)) {
	ttl := time.Tick(duration)
	kill, done := make(chan bool), make(chan bool)
	go fn(kill, done)
	for {
		select {
		case _ = <-ttl:
			kill <- true
		case _ = <-done:
			return
		default:
		}
	}
}
