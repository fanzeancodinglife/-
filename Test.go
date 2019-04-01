package main

import (
	"fmt"
	"io"
	"net"
)

func checkerr(err error){
	if err != nil{
		println(err)
		panic(err)
	}
}
func server() net.Listener{
	listener,err := net.Listen("tcp","10.108.113.229:8009")
	checkerr(err)
	return listener
}

func main() {
	listener := server()
	lisConn,err := listener.Accept()
	checkerr(err)
	cliConn, err := net.Dial("tcp", "10.108.113.229:8009")
	checkerr(err)
	s := "hello"
	s1 := []byte(s)
	cliConn.Write(s1)
	b := make([]byte ,100)
	for {
		_ ,err := cliConn.Read(b)
		if err != nil{
			if err == io.EOF {
				fmt.Println("Cloesd")
			} else {
				fmt.Printf("Read error: %s\n", err)
			}
		} else {
			break
		}
	}
	println(string(b))
	defer lisConn.Close()
	defer cliConn.Close()
}
