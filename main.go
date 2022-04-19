package main

// TCP chat room server

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost" // Hostname of the server, change to 0.0.0.0 for external connections
	CONN_PORT = "8080"
	CONN_TYPE = "tcp"
)

type Message struct {
	msg      string
	username string
	sender   net.Conn
	raw      bool
	internal bool
}

type Client struct {
	conn     net.Conn
	username string
}

// Channels
var broadcaster = make(chan Message)

// clients slice
var clients = make([]Client, 0)

// List of chat log
var chatLog = make([]string, 0, 256)

func main() {
	// Create a new TCP server
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	go broadcastMessages()
	for { // Listen for incoming connections.
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		// Create client
		this_client := Client{conn: conn}
		// Add connection to clients
		clients = append(clients, this_client)
		fmt.Println("New connection from " + conn.RemoteAddr().String())
		// Handle connections in a new goroutine.
		go clientListener(conn)
	}
}

// Register username to client connection
func registerUsername(conn net.Conn, username string) {
	for i, client := range clients {
		if client.conn == conn {
			fmt.Println("Registering username " + username + " to connection: " + conn.RemoteAddr().String())
			clients[i].username = username
			return
		}
	}
}

// Broadcasts messages to all connected clients.
func broadcastMessages() {
	for {
		msg := <-broadcaster
		if msg.internal { // For internal messaging
			switch msg.msg {
			case "remove":
				fmt.Println("Removing client " + msg.username)
				chatLog = append(chatLog, time.Now().Format("15:04:05")+" "+msg.username+" disconnected")
				for i, client := range clients {
					clientConn := client.conn
					if clientConn == msg.sender {
						clients = append(clients[:i], clients[i+1:]...)
						break
					}
				}
				continue
			case "new_client":
				chatLog = append(chatLog, time.Now().Format("15:04:05")+" "+msg.username+" connected.")
				registerUsername(msg.sender, msg.username)
				continue
			}
		}
		buf := make([]byte, 1024)
		if msg.raw {
			buf = []byte(msg.msg)
		} else {
			buf = []byte(constructMsg(msg.username, msg.msg))
			chatLog = append(chatLog, time.Now().Format("15:04:05")+" "+msg.username+": "+msg.msg)
		}
		for _, client := range clients {
			clientConn := client.conn
			if clientConn != msg.sender {
				clientConn.Write([]byte(buf))
			}
		}
	}
}

func handleCommand(data string, conn net.Conn) {
	dataTrimmed := strings.TrimSuffix(strings.TrimSpace(data), "\n")
	splitData := strings.Split(dataTrimmed, " ")
	baseCommand := splitData[0][1:]
	//args := splitData[1:]
	returnBuffer := ""
	switch baseCommand {
	case "help":
		returnBuffer = `
[System]
Available commands:
	/help - displays this help message
	/list - lists all connected clients
	/log - displays the chat log`
	case "list":
		returnBuffer = "[System]\nConnected clients: "
		for _, client := range clients {
			clientConn := client.conn
			returnBuffer += client.username + " (" + clientConn.RemoteAddr().String() + ")" + " "
		}
	case "log":
		returnBuffer = "[System]\nChat log: "
		for _, log := range chatLog {
			returnBuffer += log + "\n"
		}
	default:
		returnBuffer = "[System]\nUnknown command: " + baseCommand
	}
	conn.Write([]byte(returnBuffer + "\n"))
}

func constructMsg(usr string, msg string) string {
	return "\n" + usr + ": " + msg + "\n> "
}

// Handles incoming requests.
func clientListener(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	username := ""
	conn.Write([]byte("Username: "))

	for {
		// Read the incoming connection into the buffer.
		_, err := conn.Read(buf)

		if err != nil {
			if err.Error() == "EOF" {
				broadcaster <- Message{msg: "[System] " + username + " has left the chat.\n> ", sender: conn, username: username, raw: true}
				broadcaster <- Message{msg: "remove", sender: conn, username: username, raw: true, internal: true}
				fmt.Println("Client disconnected")
			} else {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}
		buf = bytes.Trim(buf, "\x00")
		data := strings.TrimSpace(string(buf))
		if data == "" {
			continue
		}
		if username == "" {
			username = data

			conn.Write([]byte("Welcome " + username + "!\n"))
			broadcaster <- Message{msg: "[System] " + username + " has joined the chat.\n> ", sender: conn, username: username, raw: true}
			broadcaster <- Message{msg: "new_client", sender: conn, username: username, raw: true, internal: true}
		} else if strings.HasPrefix(data, "/") {
			handleCommand(data, conn)
		} else {
			broadcaster <- Message{msg: data, sender: conn, username: username, raw: false}
		}
		fmt.Printf("Recv: %q\n", data)

		conn.Write([]byte("> "))
		buf = make([]byte, 1024)
	}

	// Close the connection when you're done with it.
	conn.Close()
}
