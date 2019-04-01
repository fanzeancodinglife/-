package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)
var wg sync.WaitGroup
const (
	SEVER_NETWORK = "tcp"
	SEVER_ADDRESS = "127.0.0.1:8000"
	DELIMITER = '\t'
)

func serverGo (){
	defer  wg.Done()
	var listener net.Listener
	listener, err := net.Listen(SEVER_NETWORK,SEVER_ADDRESS)
	if err != nil {
		log.Printf("Listen Error: %s",err)
		return
	}
	defer listener.Close()
	log.Printf("Got listener for the server.(local address :%s)",listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept Error: %s",err)
		}
		log.Printf("建立连接，远程地址为：%s",conn.RemoteAddr())
		go handleConn(conn)
	}
}

func clientGo(id int){
	defer  wg.Done()
	conn, err := net.DialTimeout(SEVER_NETWORK,SEVER_ADDRESS,2*time.Second)
	if err != nil {
		log.Printf("%d客户端,连接失败,错误信息：%s",id,err)
		return
	}
	defer conn.Close()
	log.Printf("%d号机器已经建立连接，远程地址为：%s,本地地址为：%s",id,conn.RemoteAddr(),conn.LocalAddr())
	time.Sleep(200*time.Millisecond)

	requestNumber := 5
	conn.SetDeadline(time.Now().Add(5*time.Millisecond))
	for i := 0; i < requestNumber; i++{
		req := rand.Int31()
		n, err := write(conn,fmt.Sprintf("%d",req))
		if err != nil{
			log.Printf("%d号机器报错(已经写入%dbyte)，错误：%s",id,n,err)
			continue
		}
		log.Printf("%d号机器写入成功",id)
	}
//	接收端程序
	for j := 0; j < requestNumber;j++{
		strResp, err := read(conn)
		if err != nil{
			if err != io.EOF{
				log.Printf("%d号连接关闭！",id)
			} else {
				log.Printf("%d号读错误！",id)
			}
			break
		}
		log.Printf("%d号接到服务器响应，结果为%s",id,strResp)
	}

}
func handleConn(conn net.Conn){
	defer conn.Close()
	for{
		conn.SetReadDeadline(time.Now().Add(10*time.Second))
		strReq, err := read(conn)
		if err != nil{
			if err == io.EOF{
				log.Printf("连接关闭！")
			} else{
				log.Printf("")
			}
			break
		}
		log.Printf("接到请求：%s",strReq)
		intReq, err := strconv.Atoi(strReq)
		if err != nil{
			n, err := write(conn, err.Error())
			log.Printf("写入了%d错误信息：%s",n,err)
			continue
		}
		SumResp := intReq + intReq
		respMsg := fmt.Sprintf("原参数：%d，计算结果：%d",intReq,SumResp)
		_, err = write(conn,respMsg)
		if err != nil {
			log.Printf("写入错误！%s",err)
		}
		log.Printf("发送服务端响应！%s",respMsg)
	}
}


func read(conn net.Conn)(string,error){
	readBytes := make([]byte,1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes)
		if err != nil{
			return "",err
		}
		// 一个字节一个字节读
		readByte := readBytes[0]
		if readByte == DELIMITER{
			break
		}
		buffer.WriteByte(readByte)
	}
	return  buffer.String(), nil
}

func write(conn net.Conn,content string)(int, error){
	var buffer bytes.Buffer
	buffer.WriteString(content)
	buffer.WriteByte(DELIMITER)
	return conn.Write(buffer.Bytes())
}


func main(){
	wg.Add(2)
	go serverGo()
	time.Sleep(5*time.Millisecond)
	go clientGo(1)
	wg.Wait()
}