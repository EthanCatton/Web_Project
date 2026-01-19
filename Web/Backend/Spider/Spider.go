package main

import (
	"encoding/json"
	"fmt" //basic functions
	"net"
	"time"

	"github.com/gocolly/colly" //scraper library
	"github.com/gocolly/colly/queue"
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

var site_counter int = 1


// Covers actually crawling
func main() {
	fmt.Println("Running")

	//Important vars

	var emptystruct Webpage
	crawl_queue, _ := queue.New(
		2, 
		&queue.InMemoryQueueStorage{MaxSize: 10000}, 
	)


	//Setup spider
	spider := colly.NewCollector() 
	
	spider.SetRequestTimeout(15 * time.Second)  
	spider.IgnoreRobotsTxt = false

	//gpt debug code - spider has had some weird behaviour around rate limiting of some kind
	//spider.Limit(&colly.LimitRule{
	//	DomainGlob:  "*",             // Apply to all subdomains of example.com
	//	RandomDelay: 1*time.Second, // Add a random delay of up to 1 second
	//	Parallelism: 2,               // Max concurrent requests
	})

	//new page
	spider.OnRequest(func(r *colly.Request) {
		//fmt.Println("new url")
		page := &Webpage{
			URL: r.URL.String(),
			Id:  site_counter,
		}
		r.Ctx.Put("page_data", page)
		site_counter++
	})

	//when hyperlink found
	spider.OnHTML("a[href]", func(ele *colly.HTMLElement) {
		//fmt.Println("link found")
		newlink := ele.Request.AbsoluteURL(ele.Attr("href"))
		//ele.Request.Visit(newlink) 
		crawl_queue.AddURL(newlink)  
	})

	//when title found
	spider.OnHTML("title", func(ele *colly.HTMLElement) {
		page := ele.Request.Ctx.GetAny("page_data").(*Webpage)
		page.Name = ele.Text
	})

	//h1 found ---- fallback to title
	spider.OnHTML("h1", func(ele *colly.HTMLElement) {
		page := ele.Request.Ctx.GetAny("page_data").(*Webpage)

		if page.Name == "" {
			page.Name = ele.Text
		}
	})

	spider.OnHTML("p", func(ele *colly.HTMLElement) {
		page := ele.Request.Ctx.GetAny("page_data").(*Webpage)
		page.Content += " " + ele.Text
	})

	//on error
	spider.OnError(func(r *colly.Response, err error) {
		fmt.Println(red + "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		fmt.Println(red+"error 1:", err)
		fmt.Println(red + "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		if r.StatusCode == 429 {
			return
		}
	})

	//when finished with a page
	spider.OnScraped(func(r *colly.Response) {
		fmt.Println("page scraped")

		//fmt.Println(green+"page data:", page_data)
		//fmt.Println(reflect.TypeOf(page_data))
		page := r.Ctx.GetAny("page_data").(*Webpage)
		datahandling(*page, emptystruct)
	})

	//Actual Main from here
  
	//OLD Crawling
	//spider.Visit("https://www.scrapethissite.com/pages/simple/")
	//debug start
	fmt.Println("starting crawl")
	

	//NEW crawling  
	crawl_queue.AddURL("https://www.scrapethissite.com/")  
	crawl_queue.Run(spider)

	//wait to stop any weird queue behaviour
	fmt.Println("spider.wait triggered")
	spider.Wait()

}

//manages sending data to next step
func datahandling(page_data Webpage, emptystruct Webpage) {
	fmt.Println(blue + "spider attempting to send data")
	send_data, err := json.Marshal(page_data)
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
	buffer := make([]byte, 200000)
	d, err := conn.Read(buffer)

	//recieved commented out because its not needed right now
	recieved := string(buffer[:d])
	if recieved == "full" {
		fmt.Println(blue + "Data sent and continue recieved")
	}
	_ = recieved

	page_data = emptystruct
}
