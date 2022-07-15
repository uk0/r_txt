package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	ui "github.com/gizak/termui/v3"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"r_txt/lib"
	"r_txt/xreader"
	"strings"
	"time"
)

var (
	//pageHost = "https://wap.kanshu5.la"
	pageHost = "https://m.51kanshu.cc"
)

var err error

type PageDataAndNextChapter struct {
	Data        []string
	NextChapter string
	PrevChapter string
	CurrChapter string
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

// start server
func start_server() {
	http.HandleFunc("/nl", nextLine)
	http.HandleFunc("/pl", prveLine)
	http.HandleFunc("/np", nextOnePage)
	http.HandleFunc("/up", prevOnePage)
	http.ListenAndServe(":18311", nil)
}

// nextLine
func nextLine(w http.ResponseWriter, rx *http.Request) {

	enableCors(&w)
	if rx.Method == "GET" {
		fmt.Fprintf(w, "[ %s ]", r.Next())
	}
}

func prevOnePage(w http.ResponseWriter, rx *http.Request) {
	enableCors(&w)
	if rx.Method == "GET" {
		pdc := getTxt(r.GetPrevChapter())
		savePos(pdc.CurrChapter)
		r.Load(pdc.Data, pdc.PrevChapter, pdc.NextChapter, pdc.CurrChapter)
		fmt.Fprintf(w, "[ %s ]", "上一章")
	}
}

func nextOnePage(w http.ResponseWriter, rx *http.Request) {
	enableCors(&w)
	if rx.Method == "GET" {
		pdc := getTxt(r.GetNextChapter())
		savePos(pdc.CurrChapter)
		r.Load(pdc.Data, pdc.PrevChapter, pdc.NextChapter, pdc.CurrChapter)
		fmt.Fprintf(w, "[ %s ]", "下一章")
	}
}

func prveLine(w http.ResponseWriter, rx *http.Request) {
	enableCors(&w)
	if rx.Method == "GET" {
		fmt.Fprintf(w, "[ %s ]", r.Prev())
	}
}

func convrtToUTF8(str string, origEncoding string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
	strBytes, _ = ioutil.ReadAll(reader)
	fmt.Println(string(strBytes))
	return string(strBytes)
}

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func nextPage(jumpPage string) io.Reader {
	if getUseProxy(pageHost) {
		proxyUrl, _ := url.Parse("http://127.0.0.1:9998")
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}

	resp, err := http.Get(fmt.Sprintf("%s%s", pageHost, jumpPage))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if resp.StatusCode == 200 {
		return resp.Body
	}
	return nil

}

func SplitSubWarpWidth(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}
	return subs
}

func getContent(host string, doc *goquery.Document) string {
	if strings.Contains(host, "biqugeu") || strings.Contains(host, "kanshu5") || strings.Contains(host, "51kanshu") {
		return doc.Find("div#chaptercontent").Text()
	}
	return ""
}

func getUseProxy(host string) bool {
	if strings.Contains(host, "biqugeu") {
		return true
	}
	if strings.Contains(host, "kanshu5") || strings.Contains(host, "51kanshu") {
		return false
	}
	return false
}

/**
<a href="/73/73133/79540004.html" id="pt_prev" class="Readpage_up">上一页</a>
<a href="/73/73133/79540004_3.html" id="pt_next" class="Readpage_down js_page_down">下一页</a>
**/
func getPageUp(host string, doc *goquery.Document) string {
	if strings.Contains(host, "biqugeu") {
		prevPageUrl, _ := doc.Find("a#pb_prev").Attr("href")
		return prevPageUrl
	}
	if strings.Contains(host, "kanshu5") {
		prevPageUrl, _ := doc.Find("a#pt_prev").Attr("href")
		return prevPageUrl
	}

	if strings.Contains(host, "51kanshu") {
		prevPageUrl, _ := doc.Find("a#pb_prev").Attr("href")
		return prevPageUrl
	}

	return ""
}
func getPageDown(host string, doc *goquery.Document) string {
	if strings.Contains(host, "biqugeu") {
		nextPageUrl, _ := doc.Find("a#pb_next").Attr("href")
		return nextPageUrl
	}
	if strings.Contains(host, "kanshu5") {
		nextPageUrl, _ := doc.Find("a#pt_next").Attr("href")
		return nextPageUrl
	}

	if strings.Contains(host, "51kanshu") {
		nextPageUrl, _ := doc.Find("a#pb_next").Attr("href")
		return nextPageUrl
	}

	return ""
}

func savePos(url string) {
	os.Mkdir("tmp/", 0777)
	err := ioutil.WriteFile("tmp/pos", []byte(url), 0644)
	if err != nil {
		panic(err)
	}
}

func getTxt(jumpPage string) PageDataAndNextChapter {
	txt := nextPage(jumpPage)

	doc, err := goquery.NewDocumentFromReader(txt)
	if err != nil {
		log.Fatal(err)
	}
	var content string
	content = getContent(pageHost, doc)

	_ = ioutil.WriteFile("tmp/text_test.txt", []byte(content), 0777)

	textALLPageData := strings.Replace(content, "<br/> ", "", -1)
	// 去除空格
	textALLPageData = strings.Replace(textALLPageData, " ", "", -1)
	// 去除换行符
	textALLPageData = strings.Replace(textALLPageData, "\n", "", -1)

	decoder := mahonia.NewDecoder(CheckCharSet(pageHost))
	textALLPageData = decoder.ConvertString(textALLPageData)
	TerTxt := SplitSubWarpWidth(strings.TrimSpace(textALLPageData), 25)

	nextPageUrl := getPageDown(pageHost, doc)
	prevPageUrl := getPageUp(pageHost, doc)

	return PageDataAndNextChapter{TerTxt, nextPageUrl, prevPageUrl, jumpPage}
}

func CheckCharSet(host string) string {
	if strings.Contains(host, "biqugeu") {
		return "GBK"
	}
	return "UTF-8"
}
func handleEvents() {

	uiEvents := ui.PollEvents()
	defer ui.Close()
	for {
		e := <-uiEvents
		switch e.ID {
		case "b":
			// boss key
			displayBossKey(r.Current())
		case "?":
			// show the help menu
			displayHelp(r.Current())
		case "p":
			// show the progress
			displayProgress(r.Current(), r.GetProgress())
		case "f":
			// show the frame
			displayBorder()
		case "q", "<C-c>":
			// quit
			savePos(r.GetCurrentChapter())
			time.Sleep(1 * time.Second)
			return
		case "c":
			pdc := getTxt(r.GetNextChapter())
			savePos(pdc.CurrChapter)
			r.Load(pdc.Data, pdc.PrevChapter, pdc.NextChapter, pdc.CurrChapter)
		case "j", "<Space>", "<Enter>":
			if rowNumber == "" {
				// show the next content
				updateParagraph(r.Next())
			} else {
				// parse the row number
				if num, err := lib.ParseRowNum(rowNumber); err != nil {
					updateParagraph(err.Error())
				} else {
					updateParagraph(r.Goto(r.CurrentPos() + num))
				}
				rowNumber = ""
			}
		case "k":
			if rowNumber == "" {
				// show the previous content
				updateParagraph(r.Prev())
			} else {
				// parse the row number
				if num, err := lib.ParseRowNum(rowNumber); err != nil {
					updateParagraph(err.Error())
				} else {
					updateParagraph(r.Goto(r.CurrentPos() + 1 - num))
				}
				rowNumber = ""
			}
		case "G":
			// jump to the specified row
			if rowNumber == "" {
				// jump to the last row
				updateParagraph(r.Last())
			} else {
				// parse the row number
				if num, err := lib.ParseRowNum(rowNumber); err != nil {
					updateParagraph(err.Error())
				} else {
					updateParagraph(r.Goto(num))
				}
				rowNumber = ""
			}
		case "g":
			if rowNumber == "g" {
				// jump to the first row
				updateParagraph(r.First())
				rowNumber = ""
			} else {
				rowNumber = "g"
			}
		case "w":
			color++
			// change front color
			switchColor()
		case "t":
			// timer
			setTimer()
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// jump to rows
			rowNumber += e.ID
			updateParagraph(rowNumber)
		}
		ui.Render(p)
	}
}

func main() {
	var urlPath = ""
	if len(os.Args) == 1 || len(os.Args) > 2 {
		fmt.Println("Use Config File || Modify Please input the Url")
		//os.Exit(1)
		tmp, _ := ioutil.ReadFile("tmp/pos")
		urlPath = string(tmp)
	} else {
		urlPath = os.Args[1]

	}

	r = xreader.XReader(xreader.NewTxtReader())

	pdc := getTxt(urlPath)
	if err := r.Load(pdc.Data, pdc.PrevChapter, pdc.NextChapter, urlPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go start_server()
	Init(r)

}
