// -----------------------------------------------------------------------------
// Go Language Experiments            go-experiments/[tls_socket_server_demo.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

// TLS (transport layer security) - Server

package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"time"
)

var _ = runSocketServerWithTLS
var _ = handleConnection
var _ = runSocketClientWithTLS
var _ = tlsSocketServerDemo

// runSocketServerWithTLS _ _
func runSocketServerWithTLS() {
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println("Server failed loading keys:", err)
		return
	}
	fmt.Println("Server loaded keys...")
	//
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		fmt.Println("Server failed listening:", err)
		return
	}
	defer ln.Close()
	//
	fmt.Println("Server listening for incoming connections...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Server failed accepting a connection:", err)
			continue
		}
		go handleConnection(conn)
	}
} //                                                      runSocketServerWithTLS

// handleConnection _ _
func handleConnection(cn net.Conn) {
	fmt.Println("Server handling connection...")
	defer cn.Close()
	rd := bufio.NewReader(cn)
	for {
		msg, err := rd.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Server received io.EOF from client. Exiting handler.")
			break
		}
		if err != nil {
			fmt.Println("Server failed reading from connection:", err)
			break
		}
		println("Server received message:", msg)
		reply := "World"
		n, err := cn.Write([]byte(reply))
		if err != nil {
			fmt.Println("Server failed writing to connection:", n, err)
			break
		}
		fmt.Printf("Server sent %q\n", reply)
	}
} //                                                            handleConnection

// -----------------------------------------------------------------------------

// runSocketClientWithTLS _ _
func runSocketClientWithTLS() {
	fmt.Println("Client starting...")
	cert, err := tls.LoadX509KeyPair(
		"demo.crt", // certFile string
		"demo.key", // keyFile string
	)
	if err != nil {
		fmt.Println("Client failed to load keys:", err)
		return
	}
	fmt.Println("Client loaded keys...")
	//
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "127.0.0.1:443", &config)
	if err != nil {
		fmt.Println("Client failed dialling:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected to:", conn.RemoteAddr())
	//
	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Println("Client connected to peer key:")
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println("Client connected to peer subject:")
		fmt.Println(v.Subject)
	}
	fmt.Println("Client handshake: ", state.HandshakeComplete)
	//
	message := "Hello\n"
	n, err := io.WriteString(conn, message)
	if err != nil {
		fmt.Println("Client failed writing:", err)
	}
	fmt.Printf("Client wrote %q (%d bytes)\n", message, n)
	//
	reply := make([]byte, 256)
	n, err = conn.Read(reply)
	if err != nil {
		fmt.Println("Client error reading connection:", err)
	}
	fmt.Printf("Client read %q (%d bytes)\n", string(reply[:n]), n)
	fmt.Println("Client exiting")
} //                                                      runSocketClientWithTLS

// -----------------------------------------------------------------------------

// tlsSocketServerDemo _ _
func tlsSocketServerDemo() {
	fmt.Println(div)
	fmt.Println("Running tlsSocketServerDemo")
	go runSocketServerWithTLS()
	time.Sleep(1 * time.Second)
	go runSocketClientWithTLS()
	//
	time.Sleep(10 * time.Second)
} //                                                         tlsSocketServerDemo

// end
