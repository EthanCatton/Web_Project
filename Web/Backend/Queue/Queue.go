package main

//buffers incoming data before processing to prevent dataloss when throughput doesnt match
import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// Colours
var grey = "\033[30m"
var red = "\033[31m"
var green = "\033[32m"
var blue = "\033[34m"
var purple = "\033[35m"
var pink = "\033[91m"

// Data Storage
type Webpage struct {
	URL     string
	Name    string
	Content string
	Id      int
	Score   int
}

var queue []Webpage
var lock sync.Mutex
var checklock sync.Mutex
var checker string = ""
var counter int = 1

func main() {
	fmt.Println("Running")
	lis, err := net.Listen("tcp", "192.168.57.5:5757")
	if err != nil {
		fmt.Println(red+"error 1: ", err)

	}
	defer lis.Close()

	//infinite loop to catch incoming
	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println(red+"error 2:", err)
			continue
		}

		if conn != nil {
			go fetchspider(conn)
		}
	}
}

func fetchspider(conn net.Conn) {
	defer conn.Close()
	iobuffer := bufio.NewReader(conn)
	for {
		raw_data, err := iobuffer.ReadBytes('\n')
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			fmt.Println(red+"error 2.5", err)
			continue
		}
		var web_data Webpage
		if len(raw_data) != 0 {
			err = json.Unmarshal(raw_data, &web_data)
		}
		if err != nil {
			fmt.Println(red+"error 3", err)
		}
		if web_data.URL != "" {
			//remove if counter fucks up with routines
			fmt.Println(green+"Recieved data:", counter)
			counter += 1

			//trying go
			go managequeue(web_data)
			_, err = conn.Write([]byte("full"))
			//fmt.Println(web_data.URL)
			if err != nil {
				fmt.Println(red+"error 4: ", err)
			}

		} else {
			//Continue needs to occur even on blank otherwise deadlock occurs
			_, err = conn.Write([]byte("blank"))
			//fmt.Println(web_data.URL)
			if err != nil {
				//tcp error spam
				fmt.Println(red+"error 5:", err)
				return
			}
		}
	}

}

func managequeue(web_data Webpage) {
	var prev_entry Webpage
	lock.Lock()
	//
	//fmt.Println(purple+"queue len:", len(queue))
	//
	fmt.Println("adding to queue") 
	if len(queue) > 0 {
		prev_entry = queue[len(queue)-1]
		if prev_entry.URL != web_data.URL {
			queue = append(queue, web_data)
			//fmt.Println(queue[len(queue)-1].URL)

		}

	} else {
		queue = append(queue, web_data)
		go sendtoproc()
	}
	fmt.Println("queue length:", len(queue)) 
	lock.Unlock() 

}

func sendtoproc() {
	lock.Lock()
	if len(queue) == 0 {
		fmt.Println(red + "error 5.5, queue 0")
		lock.Unlock()
		return

	}
	page_data := queue[0]
	queue = queue[1:]
	lock.Unlock()
	//Checker may need a lock to not implode things
	checklock.Lock()
	if checker == page_data.URL {
		checklock.Unlock()
		return
	}
	checker = page_data.URL
	checklock.Unlock()
	send_data, err := json.Marshal(page_data)
	_ = err
	send_data = append(send_data, '\n')
	conn, err := net.Dial("tcp", "192.168.57.15:5757")
	//defer was here
	if err != nil {
		fmt.Println(red+"error 6:", err)
		time.Sleep(3 * time.Second)
		go sendtoproc()
		return
	}
	defer conn.Close()
	_, err = conn.Write(send_data)
	fmt.Println(green + "data sent to proc")
	if err != nil {
		fmt.Println(red+"error 7:", err)
		return
	}
	buffer := make([]byte, 200000)
	d, err := conn.Read(buffer)
	_ = err

	recieved := string(buffer[:d])
	if recieved == "received" {
		fmt.Println(pink + "Continue Recieved")
	}
	_ = recieved
	//conn.Close()

}
