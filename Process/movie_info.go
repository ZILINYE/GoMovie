package Process

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var cookie string
var urls []string
var wg sync.WaitGroup
var detail []Movie_info

type Movie_info struct {
	D_url    string // Movie Download URL
	Title    string // Movie Title
	Category string // Movie Category
	Area     string // Movie Area
	Douban   string // Douban Score
}
type Outside_pattern struct {
	Home_url   string
	Score      float64
	Thread_num int // Thread Number
}

//Initialize Movie_info.Outside_pattern and NAS API
func Initialize(homeurl string, score float64, Thread_num int) Outside_pattern {

	fmt.Println("Initializing The Spider")
	result := &Outside_pattern{Home_url: homeurl, Score: score, Thread_num: Thread_num}
	cookie = Api_cookie()

	return *result
}

// Initialize API Cookie
func Api_cookie() string {
	// Login Synology Nas Via API
	Session := "http://192.168.2.20:5000/webapi/auth.cgi?api=SYNO.API.Auth&version=2&method=login&account=***&passwd=***&session=DownloadStation&format=cookie"
	result1, _ := http.Get(Session)
	//Initial Login Cookie
	coo := result1.Cookies()[0]
	cookie = strings.Split(coo.String(), ";")[0][3:]

	return cookie
}

// Search Movie By Providing Movie Name
func Search(Movie_name string) ([]string, []string) {
	var title_list []string
	var url_list []string
	hc := http.Client{}
	form := url.Values{}
	form.Add("skey", Movie_name)
	req, _ := http.NewRequest("POST", "http://www.y80s.com/movie/search/", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)

	if err != nil {
		fmt.Println(err)
	}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	doc.Find(".h3 a").Each(func(i int, selection *goquery.Selection) {

		url, _ := selection.Attr("href")
		title := selection.Text()
		url = "http://www.y80s.com" + url
		url_list = append(url_list, url)
		title_list = append(title_list, title)

	})
	return title_list, url_list
}

// Download Specific Movie
func Download_search(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)

	}
	if resp.StatusCode != 200 {
		fmt.Println("err")
	}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	dlurl, _ := doc.Find("#myform > ul > li:nth-child(2) > span.dlname.nm > span:nth-child(2) > a").First().Attr("href")
	return dlurl
}

// Crawl the home page and get all movie page url
func Get_urls(o *Outside_pattern) []Movie_info {
	resp, err := http.Get(o.Home_url)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		fmt.Println("err")
	}

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	// Collect all movie page url
	doc.Find(".h3 a").Each(func(i int, selection *goquery.Selection) {

		url, _ := selection.Attr("href")
		url = "http://www.y80s.com" + url
		urls = append(urls, url)
	})

	//Pass urls (array) to the Get_detail function
	Get_detail(urls, o.Thread_num)
	firstmatch := Filter(detail, o.Score)

	return firstmatch

}
func Filter(Download_list []Movie_info, min_score float64) []Movie_info {
	var matched []Movie_info
	for _, value := range Download_list {
		score, _ := strconv.ParseFloat(value.Douban, 64)
		if score >= min_score {
			matched = append(matched, value)
		}
	}
	return matched
}

func Get_detail(u []string, t int) {
	// Calculate number of the tasks for each thread
	each_num := int(math.Ceil(float64(len(u)) / float64(t)))

	for i := 0; i < t; i++ {
		wg.Add(1)
		part := u[i*each_num : (i+1)*each_num]
		// Pass Url array and thread code  to multithread spider
		go Spider(part, i)

	}
	fmt.Printf("each thread can dealing with %v movie pages\n", each_num)
	fmt.Printf("%v movie found\n", len(u))
	wg.Wait()

}

func Spider(s []string, i int) {

	for _, url := range s {
		fmt.Printf("Thread %v is Processing URL : %v\n", i+1, url)
		// Http request each movie page site
		resp, err := http.Get(url)
		if err != nil {
			panic(err)

		}
		if resp.StatusCode != 200 {
			fmt.Println("err")
		}

		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		title := doc.Find("#minfo > div.info > h1").First().Text()
		cate := doc.Find("#minfo > div.info > div:nth-child(6) > span:nth-child(3) > a").First().Text()
		area := doc.Find("#minfo > div.info > div:nth-child(6) > span:nth-child(4) > a").First().Text()
		douban := doc.Find("#minfo > div.info > div:nth-child(7) > div").First().Text()

		re := regexp.MustCompile("[+-]?([0-9]*[.])?[0-9]+")
		douban = re.FindString(douban)
		if len(douban) == 0 {
			douban = "0"
		}
		dlurl, _ := doc.Find("#myform > ul > li:nth-child(2) > span.dlname.nm > span:nth-child(2) > a").First().Attr("href")
		// save  each movie page info to data set
		whole := &Movie_info{Title: title, Category: cate, Area: area, Douban: douban, D_url: dlurl}
		// append all movie page info by each thread to list
		detail = append(detail, *whole)
	}
	wg.Done()

}

func Downloader(Download_list []Movie_info) {
	if len(Download_list) > 0 {
		for _, value := range Download_list {
			err := Api(value.D_url, cookie)

			if err != nil {
				fmt.Println(err)
				err_msg := "Error : " + err.Error()
				fmt.Printf(err_msg)

			}
			if err == nil {
				m := "Downloaded Movie : " + value.Douban + strings.TrimSpace(value.Title) + "to the Synology Nas  "
				fmt.Printf("Send URL: %v to the Synology Nas \n", strings.TrimSpace(value.Title))
				msg := "Success : " + m
				fmt.Printf(msg)
			}
			fmt.Printf("%v", value.Title)
		}
	} else {
		fmt.Print("No Available Movie Found\n")
	}

}

// synology API (Send Download URL to Synology)
func Api(url, co string) error {
	result, _ := http.NewRequest("GET", "http://192.168.2.20:5000/webapi/DownloadStation/task.cgi?api=SYNO.DownloadStation.Task&version=2&method=create&uri="+url, nil)
	fmt.Printf(url)
	//Add Cookie to Synology Method API
	result.AddCookie(&http.Cookie{Name: "id", Value: cookie})

	client := &http.Client{}
	resp, _ := client.Do(result)
	_, err := ioutil.ReadAll(resp.Body)

	return err
}
