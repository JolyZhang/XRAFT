package main

import (
	"flag"
	"fmt"
	"strings"
	_"os"
	//"os"
	"net"
	"strconv"
	"log"
	"bufio"
	_"bytes"
	"encoding/gob"
	"bytes"
	_"io/ioutil"
	_"time"
	"os"
)


//发送信息
func send_msg(conn net.Conn, msg string) error{
	writer := bufio.NewWriter(conn)
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	err := enc.Encode(msg)
	if err != nil {
		log.Println("Can not connect to the server." )
		os.Exit(0)
		return err
	}
	writer.Write(buff.Bytes())
	writer.Flush()
	return nil
	//time.Sleep(time.Second)
}


//接收信息
func receive_msg(conn net.Conn) {

	//reader :=bufio.NewReader(conn)
	//fmt.Println("HHHHHHHHHHHHHHHHHHHHHHHHHHH")
	raw := make([]byte, 2048)
	len, err := conn.Read(raw)
	if err != nil {
		log.Fatal(err)
	}
	buff := bytes.NewBuffer(raw[:len])
	var param string
	dec := gob.NewDecoder(buff)
	err = dec.Decode(&param)

	if err == nil {
		fmt.Println(param)
	} else {
		log.Println("Can not connect to the server." )
		os.Exit(0)
		//fmt.Println("KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK")
	}

	//time.Sleep(time.Second)
}

func on_operation(ip string, port int ,signal chan int){
	//defer func (){x<-0}()
	flag := false




	conn, err := net.Dial("tcp",ip + ":"+ strconv.Itoa(port))
	if err != nil{
		log.Fatal(err)
	}

	for !flag{
		fmt.Print(ip,port,":")
		var op, key, val string
		fmt.Scanln(&op, &key, &val)

		switch strings.ToUpper(op) {
		case "GET":
			//fmt.Println(op, key, val)
			if key != "" && val == ""{
				msg := op + "," + key
				send_msg(conn, msg)
				receive_msg(conn)
			}else{
				log.Println("COMMAND ERROR.")
				on_menu()
			}

		case "PUT":
			if key != "" && val != ""{
				msg := op + "," + key +","+val
				send_msg(conn, msg)
				receive_msg(conn)
			}else{
				log.Println("COMMAND ERROR.")
				on_menu()
			}
		case "APPEND":
			if key != "" && val != ""{
				msg := op + "," + key +"," + val
				send_msg(conn, msg)
				receive_msg(conn)
			}else{
				log.Println("COMMAND ERROR.")
				on_menu()
			}
		case "DELETE":
			if key != "" && val == ""{
				msg := op + "," + key
				send_msg(conn, msg)
				receive_msg(conn)
			}else{
				log.Println("COMMAND ERROR.")
				on_menu()
			}
		case "QUIT":
			//flg <-0
			//os.Exit(0)
			flag = true
			//break
		default:
			log.Println("COMMAND ERROR.")
			 on_menu()

		}


	}
	signal<-0


}

func on_menu(){
	prompt := "┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓\n" +
		"┃                WELCOME TO USE THIS DISTRIBUTED KV-STORE                ┃\n" +
		"┃  The Commands below are available:                                     ┃\n" +
		"┃     -GET Key: Query the value specified by key.                        ┃\n" +
		"┃     -PUT Key value: Store a new key-value pair.                        ┃\n" +
		"┃     -Append Key value:                                                 ┃\n" +
		"┃          Append new value on an exsist key-value pair. If the key do   ┃\n" +
		"┃      not exsist, store the new key-value pair as PUT does.             ┃\n" +
		"┃     -DELETE key: Delete a key-value pair specified by key. If the key  ┃\n" +
		"┃      DO NOT exsist, there is NO effect.                                ┃\n" +
		"┃     -QUIT: Quit the system.                                            ┃\n" +
			"┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛"
			fmt.Println(prompt)

}
func main(){
	ip:= flag.String("ip","localhost","The IP Address of Server.")
	port := flag.Int("port", 2000, "The port of Server.")

	flag.Parse()
	var x chan int
	x = make(chan int )

	on_menu()
	go on_operation(*ip, *port,x)

	select {

		case <-x:
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>Good bye.<<<<<<<<<<<<<<<<<<<<<<<")
	}
	//fmt.Println("Good Bye.",<-x)

	//fmt.Println(<-x)
}