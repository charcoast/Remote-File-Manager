package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
)

const DiscoveryPort = 7070

var port string
var discoveryIP string
var prefixes = []string{"list", "li", "ls"}

func main() {
	port = getPortFromOSArgs(os.Args)
	discoveryIP = getDiscoveryIPFromOSArgs(os.Args)
	listener, _ := openTCPListener(port)
	defer closeTCPListener(listener)
	sysFormat("TCP Server initialized at port", port)
	go communicateDiscovery()
	for {
		conn, _ := openConnection(listener)
		go parallelReader(conn)
	}
}

func communicateDiscovery() {
	for {
		conn, err := net.Dial("tcp", discoveryIP+":"+strconv.Itoa(DiscoveryPort))
		if err != nil {
			continue
		}

		selfPort, _ := strconv.Atoi(port)
		featureRegister := FeatureRegister{Port: selfPort, Prefixes: prefixes}
		data, _ := json.Marshal(featureRegister)
		sendData(conn, data)
		closeConnection(conn)
		return
	}
}

func parallelReader(conn net.Conn) {
	defer closeConnection(conn)
	for {
		addr := conn.RemoteAddr()
		data, _ := readConnectionData(conn)
		upperTrimSpaceData := strings.ToUpper(strings.TrimSpace(data))
		sysOut(upperTrimSpaceData)
		if upperTrimSpaceData == "PING" {
			sendData(conn, []byte("PONG"))
		}
		if slices.Contains(prefixes, strings.TrimSpace(data)) {
			sysOut("lendo a list")
			dirs := readDir()
			var directories = make([]string, 0)
			for _, dir := range dirs {
				dirName := strings.TrimSpace(dir.Name())
				if dirName == "" || len(dirName) == 0 {
					continue
				}
				directories = append(directories, dirName)
			}
			sendData(conn, []byte(strings.Join(directories, ",")))
		}
		if upperTrimSpaceData == "EXIT" || upperTrimSpaceData == "QUIT" || upperTrimSpaceData == "EOF" {
			sysFormat("Closing connection with %s", addr)
			return
		}
	}
}

func readDir() []os.DirEntry {
	dir, err := os.ReadDir("./")
	if err != nil {
		return nil
	}
	return dir
}

func getPortFromOSArgs(args []string) string {
	if len(args) == 1 {
		fmt.Println("Enter with port number in argument")
		panic("The port number must be entered")
	}
	return args[1]
}

func getDiscoveryIPFromOSArgs(args []string) string {
	if len(args) == 2 {
		fmt.Println("Enter with port number and discovery IP in argument")
		panic("The discovery IP must be entered")
	}
	return args[2]
}

func openTCPListener(port string) (net.Listener, error) {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return l, nil
}

func closeTCPListener(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func openConnection(listener net.Listener) (net.Conn, error) {
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return conn, nil
}

func closeConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func readConnectionData(conn net.Conn) (string, error) {
	data, err := bufio.NewReader(conn).ReadString('\n')
	if data == "" && err != nil {
		sysOut("data is empty")
		data = err.Error()
	}
	return data, nil
}

func sendData(conn net.Conn, b []byte) {
	_, err := fmt.Fprintf(conn, string(b)+"\n")
	if err != nil {
		sysOut(err.Error())
		return
	}
}

func sysOut(value interface{}) {
	fmt.Println(value)
}

func sysFormat(s string, a ...interface{}) {
	fmt.Printf(s, a...)
}
