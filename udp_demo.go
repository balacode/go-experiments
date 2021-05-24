// -----------------------------------------------------------------------------
// Go Language Experiments                          go-experiments/[udp_demo.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

var _ = udpServer
var _ = udpClient
var _ = udpDemo

//- MAX_BUFFER_SIZE specifies the size of the buffers that are used
//- to temporarily hold data from the UDP packets that we receive.
const MAX_BUFFER_SIZE = 1024

// -----------------------------------------------------------------------------

//- udpServer wraps all the UDP echo server functionality.
//- ps.: the server is capable of answering to a single client at a time.
func udpServer(address string) (err error) {
	//- ListenPacket provides us a wrapper around ListenUDP so
	//- that we don't need to call `net.ResolveUDPAddr` and then
	//- subsequentially perform a `ListenUDP` with the UDP address.
	//-
	//- The returned value (PacketConn) is pretty much the same as the one
	//- from ListenUDP (UDPConn) - the only difference is that `Packet*`
	//- methods and interfaces are more broad, also covering `ip`.
	fmt.Println("SV: started")
	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}
	fmt.Println("SV: called ListenPacket")
	//
	//- `Close`ing the packet "connection" means cleaning the data structures
	//- allocated for holding information about the listening socket.
	defer pc.Close()
	//
	buffer := make([]byte, MAX_BUFFER_SIZE)
	//
	//- Given that waiting for packets to arrive is blocking by nature and we want
	//- to be able of canceling such action if desired, we do that in a separate
	//- go routine.
	go func() {
		for {
			//- By reading from the connection into the buffer, we block until there's
			//- new content in the socket that we're listening for new packets.
			//-
			//- Whenever new packets arrive, `buffer` gets filled and we can continue
			//- the execution.
			//-
			//- note.: `buffer` is not being reset between runs.
			//-    It's expected that only `n` reads are read from it whenever
			//-    inspecting its contents.
			//
			// the contents of 'buffer' or overwritten after every call to ReadFrom.
			fmt.Println("S2: before ReadFrom. buffer:", string(buffer))
			nRead, addr, err := pc.ReadFrom(buffer)
			fmt.Println("S2: after ReadFrom. Buffer:", string(buffer))
			if err != nil {
				return
			}
			msg := string(buffer[:nRead])
			fmt.Printf("S2: received %q [%d] from %s\n",
				msg, nRead, addr.String())

			msg = strings.ToLower(msg)
			//
			//- Setting a deadline for the `write` operation allows us to not block
			//- for longer than a specific timeout.
			//-
			//- In the case of a write operation, that'd mean waiting for the
			//- send queue to be freed enough so that we are able to proceed.
			deadline := time.Now().Add(time.Second)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				return
			}
			//- Write the packet's contents back to the client.
			nWrit, err := pc.WriteTo([]byte(msg), addr)
			if err != nil {
				return
			}
			fmt.Printf("S2: replied %q [%d] bytes to %s\n",
				msg, nWrit, addr.String())
		}
	}()
	time.Sleep(2 * time.Hour)
	return nil
} // udpServer

// -----------------------------------------------------------------------------

//- udpClient wraps the whole functionality of a UDP client that sends
//- a message and waits for a response coming back from the server
//- that it initially targetted.
func udpClient(address string, message string) (err error) {
	fmt.Println("CL: started")
	//- Resolve the UDP address so that we can make use of DialUDP
	//- with an actual IP and port instead of a name (in case a
	//- hostname is specified).
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	fmt.Println("CL: resolved", raddr.IP.String())
	//
	//- Although we're not in a connection-oriented transport,
	//- the act of `dialing` is analogous to the act of performing
	//- a `connect(2)` syscall for a socket of type SOCK_DGRAM:
	//- - it forces the underlying socket to only read and write
	//-   to and from a specific remote address.
	fmt.Println("CL: dial", raddr.IP.String())
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}
	fmt.Println("CL: dial done")
	//
	//- Closes the underlying file descriptor associated with the,
	//- socket so that it no longer refers to any file.
	defer conn.Close()
	go func() {
		rd := strings.NewReader(message)
		//- It is possible that this action blocks, although this
		//- should only occur in very resource-intensive situations:
		//- - when you've filled up the socket buffer and the OS
		//-   can't dequeue the queue fast enough.
		nw, err := io.Copy(conn, rd)
		if err != nil {
			return
		}
		fmt.Printf("C2: sent %q [%d] to %s\n",
			message, nw, conn.RemoteAddr().String())
		//
		buffer := make([]byte, MAX_BUFFER_SIZE)
		//
		//- Set a deadline for the ReadOperation so that we don't
		//- wait forever for a server that might not respond on
		//- a resonable amount of time.
		deadline := time.Now().Add(time.Second)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			return
		}
		nr, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			return
		}
		fmt.Printf(
			"C2: received %q [%d] from %s\n",
			string(buffer)[:nr], nr, addr.String(),
		)
	}()
	time.Sleep(1 * time.Hour)
	return nil
} //                                                                   udpClient

// -----------------------------------------------------------------------------

func udpDemo() {
	addr := "127.0.0.1:9009"
	//
	// start the server
	go udpServer(addr)
	//
	go udpClient(addr, "AAAAAAAA")
	go udpClient(addr, "BBBBBB")
	go udpClient(addr, "CCCC")
	go udpClient(addr, "DD")
	//
	time.Sleep(1 * time.Hour)
} //                                                                     udpDemo

// end
