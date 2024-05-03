/*
A very simple TCP server written in Go.

This is a toy project that I used to learn the fundamentals of writing
Go code and doing some really basic network stuff.

Maybe it will be fun for you to read. It's not meant to be
particularly idiomatic, or well-written for that matter.
*/
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type server struct {
	Host string
	Port int
	Type string
}

var SERVERS = []server{
	{Host: "localhost", Port: 8000, Type: "tcp"},
	{Host: "localhost", Port: 8001, Type: "tcp"},
	{Host: "localhost", Port: 8002, Type: "tcp"},
	{Host: "localhost", Port: 8003, Type: "tcp"},
	// Add more servers here if needed
}

type KhoGiay struct {
	Width  int `json:"width"`
	Length int `json:"length"`
}

type KhoGiay_Kq struct {
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
	X      float64 `json:"x"`
	V      float64 `json:"v"`
}

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 9999, "The port to listen on; default is 9999.")

func main() {
	flag.Parse()

	fmt.Println("Starting load balance...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	var queue = make(chan KhoGiay, 100)

	var queue_kq = make(chan KhoGiay_Kq, 100)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		// Bắt đầu một goroutine để đọc phản hồi từ server
		var wg sync.WaitGroup
		wg.Add(1)
		go handleConnection(conn, queue)

		go readConnection(conn, queue, queue_kq)

		go func() {

			for v := range queue_kq {

				sendKg_kq(v, conn, 0)

			}

		}()

		wg.Wait()

	}

}

func handleConnection(conn net.Conn, queue chan<- KhoGiay) {
	defer conn.Close()
	// Tạo một scanner để đọc dữ liệu từ kết nối
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// Lấy dữ liệu từ scanner
		data := scanner.Bytes()

		// Phân tích dữ liệu JSON
		var kg KhoGiay
		err := json.Unmarshal(data, &kg)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		queue <- kg

	}
}

func readConnection(connClient net.Conn, queue <-chan KhoGiay, queue_kq chan<- KhoGiay_Kq) {
	var wg sync.WaitGroup
	wg.Add(len(SERVERS))

	// Biến đếm để lựa chọn server
	counter := 0

	for v := range queue {
		// Tăng giá trị của biến đếm và lấy phần dư để lựa chọn server
		counter = (counter + 1) % len(SERVERS)
		serverIndex := counter

		// Lấy server dựa trên giá trị của biến đếm
		server := SERVERS[serverIndex]

		// Kết nối đến server và gửi dữ liệu
		connSer, err := net.Dial(server.Type, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			fmt.Printf("Error connecting to server (thread %d): %v\n", serverIndex, err)
			return
		}

		go func() {
			defer wg.Done()
			sendKg(v, connSer, serverIndex)
			readResponse(connSer, &wg, queue_kq)

		}()
	}

	wg.Wait()
}

func readResponse(conn net.Conn, wg *sync.WaitGroup, queue_kq chan<- KhoGiay_Kq) {
	defer wg.Done()
	defer conn.Close()

	// Tạo một scanner để đọc dữ liệu từ kết nối
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// Lấy dữ liệu từ scanner
		data := scanner.Bytes()

		// Phân tích dữ liệu JSON
		var khoGiay_Kq KhoGiay_Kq
		err := json.Unmarshal(data, &khoGiay_Kq)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// In ra dữ liệu đã phân tích được

		fmt.Println("Received data in load balance server:", khoGiay_Kq)

		queue_kq <- khoGiay_Kq

	}

}

func sendKg(kg KhoGiay, conn net.Conn, numThread int) {
	// Encode data to JSON
	jsonData, err := json.Marshal(kg)
	if err != nil {
		fmt.Println("Error encoding JSON from client:", err)
		return
	}

	// Append delimiter to JSON data
	jsonData = append(jsonData, []byte("\n")...)

	// Send data to the server
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Printf("Error sending data to server (thread %d): %v\n", numThread, err)
		return
	}

	fmt.Printf("Send data success %+v (thread %d)\n", kg, numThread)

}

func sendKg_kq(kg_kq KhoGiay_Kq, conn net.Conn, numThread int) {
	// Encode data to JSON
	jsonData, err := json.Marshal(kg_kq)
	if err != nil {
		fmt.Println("Error encoding JSON from client:", err)
		return
	}

	// Append delimiter to JSON data
	jsonData = append(jsonData, []byte("\n")...)

	// Send data to the server
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Printf("Error sending data to server (thread %d): %v\n", numThread, err)
		return
	}

	fmt.Printf("Send data success %+v (thread %d)\n", kg_kq, numThread)
}
