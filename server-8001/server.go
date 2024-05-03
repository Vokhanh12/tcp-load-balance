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
var port = flag.Int("port", 8001, "The port to listen on; default is 8001.")

func main() {
	flag.Parse()

	fmt.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	var queue = make(chan KhoGiay, 100)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go handleConnection(conn, queue)

		go readConnection(conn, queue)

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

		// In ra dữ liệu đã phân tích được
		fmt.Println("Received data:", kg)

	}
}

func readConnection(conn net.Conn, queue <-chan KhoGiay) {

	var wg sync.WaitGroup

	for v := range queue {

		width := float64(v.Width) / 100

		length := float64(v.Length) / 100

		wg.Add(1)

		Vmax := 0.0
		Xmax := 0.0

		go func() {

			defer wg.Done()

			for ix := 0.0; ix <= float64(width)/2; ix += 0.0001 {
				// x range 0 <= x <= w/2

				crV := ix * (float64(width) - 2*ix) * (float64(length) - 2*ix)

				if Vmax < crV {
					Vmax = crV
					Xmax = ix
				}

			}

		}()

		wg.Wait()

		// Tạo một cấu trúc KhoGiay với width và length là 123
		khoGiay_Kq := KhoGiay_Kq{Width: width, Length: length, X: Xmax, V: Vmax}

		// Chuyển cấu trúc thành JSON
		jsonData, err := json.Marshal(khoGiay_Kq)
		if err != nil {
			fmt.Println("Error encoding JSON from client:", err)
			return
		}

		// Thêm dấu xuống dòng vào cuối dữ liệu JSON
		jsonData = append(jsonData, []byte("\n")...)

		// Sử dụng goroutine để gửi dữ liệu bất đồng bộ
		go func(data []byte) {
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println("Error writing to stream:", err)
				return
			}
		}(jsonData)

		outputString := fmt.Sprintf("%f %f %f %f numThread:%d \n", width, length, Xmax, Vmax)

		fmt.Printf(outputString)

		fmt.Println()

	}

}
