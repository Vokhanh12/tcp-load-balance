package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var host = flag.String("host", "localhost", "The hostname or IP to connect to; defaults to \"localhost\".")
var port = flag.Int("port", 9999, "The port to connect to; defaults to 9999.")

func main() {
	flag.Parse()

	// Số lượng lần chạy song song
	numRuns := 1

	// Sử dụng WaitGroup để chờ tất cả các goroutine hoàn thành
	var wg sync.WaitGroup

	for i := 0; i < numRuns; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			connectAndRun(index)
		}(i)
	}

	// Chờ cho tất cả các goroutine kết thúc trước khi kết thúc chương trình
	wg.Wait()
}

func connectAndRun(index int) {
	dest := *host + ":" + strconv.Itoa(*port)
	fmt.Printf("Connecting to %s...\n", dest)

	conn, err := net.Dial("tcp", dest)
	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("Some problem connecting.")
		} else {
			fmt.Println("Unknown error: " + err.Error())
		}
		os.Exit(1)
	}
	defer conn.Close()

	// Tạo tên file output duy nhất cho mỗi lần chạy
	fileName := fmt.Sprintf("outputDanhHCN%d.txt", index)
	fileOut, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer fileOut.Close()

	// Sử dụng WaitGroup để đợi cho goroutine đọc phản hồi từ server kết thúc
	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {
		// Bắt đầu một goroutine để đọc phản hồi từ server
		wg.Add(1)
		go func() {
			defer wg.Done()
			readResponse(conn, fileOut)
		}()

		// Gửi yêu cầu đến server
		sendRequest(conn, index)
	}

	// Chờ cho tất cả các goroutine đọc phản hồi từ server kết thúc
	wg.Wait()
}

func sendRequest(conn net.Conn, index int) {
	LINE := 15507
	countLines(conn, 0, LINE, index)
}

func countLines(conn net.Conn, start int, end int, index int) {
	filePath := "DanhHCN_31012024.txt"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	lineNumber := start
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		if lineNumber == end {
			fmt.Println("read file end")
			break
		}

		arrs := strings.Fields(line)

		crWidth, err := strconv.Atoi(arrs[0])
		if err != nil {
			fmt.Println("Loi doi width string sang int")
			return
		}

		crLength, err := strconv.Atoi(arrs[1])
		if err != nil {
			fmt.Println("Loi doi length string sang int", crWidth, crLength)
			return
		}

		// Tạo một cấu trúc KhoGiay với width và length là 123
		kg := KhoGiay{Width: crWidth, Length: crLength}

		// Chuyển cấu trúc thành JSON
		jsonData, err := json.Marshal(kg)
		if err != nil {
			fmt.Println("Error encoding JSON from client:", err)
			return
		}

		// Thêm dấu xuống dòng vào cuối dữ liệu JSON
		jsonData = append(jsonData, []byte("\n")...)

		// Gửi dữ liệu đến server
		_, err = conn.Write(jsonData)
		if err != nil {
			fmt.Println("Error writing to stream:", err)
			break
		}

		lineNumber++

		fmt.Printf("Send data success %+v (thread %d)\n", kg, index)

	}
	fmt.Println("File sent to server.")
}

func readResponse(conn net.Conn, fileOut *os.File) {
	defer conn.Close()

	addTitle := fmt.Sprintf("width	length	x	Vmax \n")
	mutex := &sync.Mutex{}
	mutex.Lock() // Lock trước khi ghi vào file
	_, _ = fileOut.WriteString(addTitle)
	mutex.Unlock() // Mở khóa sau khi ghi vào file

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
		outputString := fmt.Sprintf("%f %f %f %f\n", khoGiay_Kq.Width, khoGiay_Kq.Length, khoGiay_Kq.X, khoGiay_Kq.V)

		mutex.Lock() // Lock trước khi ghi vào file
		_, err = fileOut.WriteString(outputString)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		mutex.Unlock() // Mở khóa sau khi ghi vào file

		fmt.Println("Received data:", khoGiay_Kq)
	}
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
