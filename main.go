package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ZILINYE/GoMovie/Process"
	"os"
	"strconv"
	"strings"
	"time"
)

var spider bool
var home_url string
var min_score float64
var thread_num int

var Download_url string
var h bool
var g bool

func init() {
	flag.BoolVar(&h, "h", false, "This help")
	// Spider on 80s Website and send url to synology
	flag.BoolVar(&spider, "a", false, "Open Spider?")
	flag.StringVar(&home_url, "u", "http://www.y80s.com/movie/list", "send home url")
	flag.Float64Var(&min_score, "s", 7.5, "set Minimum Douban Score")
	flag.IntVar(&thread_num, "t", 5, "Set Thread number`")

	// Send Download Url directly to Synology
	flag.StringVar(&Download_url, "l", "none", "Input Download URL (http://, https://, ftp://, ftps://, sftp://, magnet://, thunder://, flashget://, qqdl://)`")

	// Search Movie
	flag.BoolVar(&g, "g", false, "Search Movie")

	flag.Usage = usage
}
func main() {
	config, uname, upass := Process.ReadConf() // Read Conf.json Get Data Saving Method and NAS User Name And Password

	for {
		flag.Parse()

		if h {
			flag.Usage()
			return
		}

		// Download Movie By providing the Download URL
		if Download_url != "none" {
			cookie := Process.Api_cookie(uname, upass)
			err := Process.Api(Download_url, cookie)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Send URL to Synology Successfully")
			}
			return

		} else if g {
			// Download Movie By providing Movie name
			y := 0
			for y == 0 {
				fmt.Print("Search Movie : ")
				name := bufio.NewScanner(os.Stdin)
				name.Scan()
				title_list, url_list := Process.Search(name.Text())
				fmt.Printf("|%-6s|%-12s|%-6s\n", "序号", "电影名", "详情链接")
				for i := 0; i < len(title_list); i++ {
					fmt.Printf("|%-6s|%-12s|%-6s\n", strconv.Itoa(i+1), strings.TrimSpace(title_list[i]), url_list[i])
				}
				fmt.Printf("%v : Cancel\n", y)
				fmt.Print("Choose Number : ")
				num := bufio.NewScanner(os.Stdin)
				num.Scan()
				x := num.Text()
				y, _ := strconv.Atoi(x)
				if y != 0 {
					dl_url := Process.Download_search(url_list[y-1])
					cookie := Process.Api_cookie(uname, upass)
					err := Process.Api(dl_url, cookie)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("Send URL to Synology Successfully")
					}
				}

			}
			return

		} else if spider {
			// Daily Schedule spider
			now := time.Now()
			initialize := Process.Initialize(home_url, min_score, thread_num, uname, upass)
			firstmatch := Process.Get_urls(&initialize)
			secondmatch := config.CheckRecord(firstmatch)
			Process.Downloader(secondmatch)
			end := time.Now()
			spend := end.Sub(now)
			fmt.Printf("total spend %v s\n", spend)
			return
		} else {
			// Daily Schedule spider
			fmt.Print("Main Menu\n\nSelect Option : \n1.Open Spider\n2.Send URL\n3.Search Movie\n4.Cancel\n\n\nYou Select :")
			option := bufio.NewScanner(os.Stdin)
			option.Scan()
			cnum := option.Text()
			num, _ := strconv.Atoi(cnum)
			if num == 1 {
				spider = true
				fmt.Print("Open Spider\n\nInput Home URL :")
				hurl := bufio.NewScanner(os.Stdin)
				hurl.Scan()
				if hurl.Text() != "" {
					home_url = hurl.Text()
				}

				fmt.Print("Set The Minimum Douban Score :")
				mscore := bufio.NewScanner(os.Stdin)
				mscore.Scan()
				if mscore.Text() != "" {
					min_score, _ = strconv.ParseFloat(mscore.Text(), 64)
				}

				fmt.Print("Input Thread Number :")
				tnum := bufio.NewScanner(os.Stdin)
				tnum.Scan()
				if tnum.Text() != "" {
					thread_num, _ = strconv.Atoi(tnum.Text())
				}

			} else if num == 2 {
				fmt.Print("Send URL\n\nInput Download URL :")
				durl := bufio.NewScanner(os.Stdin)
				durl.Scan()
				if durl.Text() != "" {
					Download_url = durl.Text()
				}

			} else if num == 3 {
				g = true
			} else {
				return
			}
		}
	}

}
func usage() {
	fmt.Fprintf(os.Stderr, `Auto Download Movie Version 1.0
Usage (Spider): ./Movie_search [-u Home URL] [-s Minimum Douban Score] [-t Thread number] [-h help]
Usage (Send Url): ./Movie_search [-l 'Download url(http://, https://, ftp://, ftps://, sftp://, magnet://, thunder://, flashget://, qqdl://)']
Usage (Search Movie): ./Movie_search -g
Options:
`)
	flag.PrintDefaults()
}
