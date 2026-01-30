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

var lis net.Listener

func main() {
	//catches spider data
	fmt.Println("Running")
	for {
		temp, err := net.Listen("tcp", "0.0.0.0:5757")
		lis = temp
		if err != nil {
			fmt.Println(red+"error 1: ", err)
			time.Sleep(time.Second * 5)
		}
		if lis != nil {
			defer lis.Close()
			break
		}
	}

	go proc_loop()

	//infinite loop
	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println(red+"error 2:", err)
			continue
		}

		if conn != nil {
			go fetch_spider(conn)
		}
	}
}

func fetch_spider(conn net.Conn) {
	defer conn.Close()
	//iobuffer used for throughput
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
		//remove any broken data
		if web_data.URL != "" {
			fmt.Println(green+"Recieved data:", counter)
			counter += 1

			//trying go
			go manage_queue(web_data)
			_, err = conn.Write([]byte("full"))
			//fmt.Println(web_data.URL)
			if err != nil {
				fmt.Println(red+"error 4: ", err)
			}

		} else {
			//continue needs to occur even on blank otherwise deadlock occurs
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

func manage_queue(web_data Webpage) {
	var prev_entry Webpage
	//mutexs to prevent multiple edit attempts by different routines
	lock.Lock()
	//fmt.Println(purple+"queue len:", len(queue))
	fmt.Println("adding to queue")

	if len(queue) > 0 {
		//protects against duplicates
		prev_entry = queue[len(queue)-1]
		if prev_entry.URL != web_data.URL {
			queue = append(queue, web_data)
			//fmt.Println(queue[len(queue)-1].URL)
		}
	} else {
		//starts off queue
		queue = append(queue, web_data)
	}
	fmt.Println("queue length:", len(queue))
	lock.Unlock()
}

func send_to_proc() {

	//old lock from old function trigger
	//lock.Lock()
	//empty protection
	if len(queue) == 0 {
		fmt.Println(red + "error 5.5, queue 0")
		lock.Unlock()
		return
	}

	//gets first data and then removes from queue
	page_data := queue[0]
	queue = queue[1:]
	//lock.Unlock()

	//probably not needed with rework but harmless rn
	checklock.Lock()

	//repeats duplicate check
	if checker == page_data.URL {
		checklock.Unlock()
		return
	}
	checker = page_data.URL
	checklock.Unlock()

	//similar code to spider
	send_data, err := json.Marshal(page_data)
	_ = err
	send_data = append(send_data, '\n')
	conn, err := net.Dial("tcp", "processor:5757")

	if err != nil {
		fmt.Println(red+"error 6:", err)
		time.Sleep(3 * time.Second)
		//uncertain if this should stay or not
		//go send_to_proc()
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
	data, err := conn.Read(buffer)
	_ = err

	recieved := string(buffer[:data])
	if recieved == "received" {
		fmt.Println(pink + "Continue Recieved")
	}
	_ = recieved

}

func proc_loop() {
	//send_to_proc used to be triggered in manage_queue, I don't like it being dependent on seperate goroutines so has been split
	//doesn't really need to be own function at all but it helps visually for me and sets pace independent on main loop
	pcount := 0
	for {
		lock.Lock()
		time.Sleep(1 * time.Second)
		if len(queue) > 0 {
			fmt.Println("calling proc", pcount)
			go send_to_proc()
			pcount += 1
		}
		lock.Unlock()
	}

}
