package main

import (
	"encoding/json"
	"fmt" //basic functions
	"net"
	"time"

	"github.com/gocolly/colly" //scraper library
	"github.com/gocolly/colly/queue"
)

// colours
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

var site_counter int = 1

// covers actually crawling
func main() {
	fmt.Println("Running")

	var emptystruct Webpage

	//colly specific queue
	crawl_queue, _ := queue.New(
		//parallelism, cant handle different numbers
		2,
		//prevents throughput issues later
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	//Setup spider
	spider := colly.NewCollector()
	spider.SetRequestTimeout(15 * time.Second)
	spider.IgnoreRobotsTxt = false

	//new page
	spider.OnRequest(func(response *colly.Request) {
		//fmt.Println("new url")
		page := &Webpage{
			URL: response.URL.String(),
			Id:  site_counter,
		}
		response.Ctx.Put("page_data", page)
		site_counter++
	})

	//when hyperlink found
	spider.OnHTML("a[href]", func(element *colly.HTMLElement) {
		//fmt.Println("link found")
		newlink := element.Request.AbsoluteURL(element.Attr("href"))
		crawl_queue.AddURL(newlink)
	})

	//when title found
	spider.OnHTML("title", func(element *colly.HTMLElement) {
		page := element.Request.Ctx.GetAny("page_data").(*Webpage)
		page.Name = element.Text
	})

	//h1 found ---- fallback to title
	spider.OnHTML("h1", func(element *colly.HTMLElement) {
		page := element.Request.Ctx.GetAny("page_data").(*Webpage)

		if page.Name == "" {
			page.Name = element.Text
		}
	})

	spider.OnHTML("p", func(element *colly.HTMLElement) {
		page := element.Request.Ctx.GetAny("page_data").(*Webpage)
		page.Content += " " + element.Text
	})

	//on error
	spider.OnError(func(response *colly.Response, err error) {
		fmt.Println(red+"error 1:", err)
		fmt.Println(red+"code:", response.StatusCode)
		//too many requests catch
		if response.StatusCode == 429 {
			return
		}
	})

	//when finished with a page
	spider.OnScraped(func(response *colly.Response) {
		fmt.Println("page scraped")

		//fmt.Println(green+"page data:", page_data)
		//fmt.Println(reflect.TypeOf(page_data))

		page := response.Ctx.GetAny("page_data").(*Webpage)
		data_handling(*page, emptystruct)
	})

	//Actual Main from here
	fmt.Println("starting crawl")
	crawl_queue.AddURL("https://www.scrapethissite.com/")
	crawl_queue.Run(spider)

	//wait to stop any weird queue behaviour
	fmt.Println("spider.wait triggered")
	spider.Wait()
}

// manages sending data to next step
func data_handling(page_data Webpage, emptystruct Webpage) {
	fmt.Println(blue + "spider attempting to send data")
	send_data, err := json.Marshal(page_data)
	//splits with newline
	send_data = append(send_data, '\n') //doesnt work with ""?
	conn, err := net.Dial("tcp", "192.168.57.5:5757")
	if err != nil {
		fmt.Println(red+"error 2:", err)
	}
	defer conn.Close()

	_, err = conn.Write(send_data)
	//fmt.Println(pink + "data sent")
	if err != nil {
		fmt.Println(red+"error 3:", err)
	}
	//arbitrary data size to prevent issues
	buffer := make([]byte, 200000)
	data, err := conn.Read(buffer)
	recieved := string(buffer[:data])
	//waits to recieve before moving on
	if recieved == "full" {
		fmt.Println(blue + "Data sent and continue recieved")
	}
	//wipes
	_ = recieved
	page_data = emptystruct
}
