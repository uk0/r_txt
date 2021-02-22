package main

import (
	"bufio"
	"bytes"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"r_txt/lib"
	"r_txt/xreader"
	"regexp"
	"strings"
	"time"
)

var (
	pageHost = "http://www.boquku.com"
)

var err error

type PageDataAndNextChapter struct {
	Data        []string
	NextChapter string
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
	http.ListenAndServe(":18311", nil)
}

// nextLine
func nextLine(w http.ResponseWriter, rx *http.Request) {
	enableCors(&w)
	fmt.Fprintf(w, " %s!", r.Next())
}


func nextOnePage(w http.ResponseWriter, rx *http.Request) {
	enableCors(&w)
	pdc := getTxt(r.GetNextChapter())
	savePos(pdc.CurrChapter)
	r.Load(pdc.Data, pdc.NextChapter,pdc.CurrChapter)
	fmt.Fprintf(w, " %s!", "下一章")
}


func prveLine(w http.ResponseWriter, rx *http.Request) {
	enableCors(&w)
	fmt.Fprintf(w, " %s!", r.Prev())
}



func convrtToUTF8(str string, origEncoding string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
	strBytes, _ = ioutil.ReadAll(reader)
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

func nextPage(jumpPage string) string {
	proxyUrl, err := url.Parse("http://127.0.0.1:2087")
	http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	resp, err := http.Get(fmt.Sprintf("%s%s", pageHost, jumpPage))
	if err != nil {
		fmt.Println(err)
		return "nil"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		return convrtToUTF8(string(body), "gbk")
	}
	return "nil"

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

func savePos(url string) {
	//os.Mkdir("/tmp/r_txt/", 0777)
	err := ioutil.WriteFile("/tmp/r_txt/pos", []byte(url), 0644)
	if err != nil {
		panic(err)
	}
}

func getTxt(jumpPage string) PageDataAndNextChapter {
	txt := nextPage(jumpPage)
	start := false
	end := false
	startIndex := 0
	endIndex := 0
	nextPageIndex := 0
	textDataArray := SplitLines(txt)
	for index, line := range textDataArray {
		if strings.Contains(line, "<div id=\"txtContent\">") {
			startIndex = index
			start = true
		}
		if start && strings.Contains(line, "</div>") {
			endIndex = index
			end = true
		}

		if end && start {
			if start && strings.Contains(line, "下一章") {
				nextPageIndex = index
				break
			}
		}
	}

	txtBody := strings.Join(textDataArray[startIndex:endIndex], "")

	textALLPageData := strings.Replace(txtBody, "<br/> ", "", -1)
	// 去除空格
	textALLPageData = strings.Replace(textALLPageData, " ", "", -1)
	// 去除换行符
	textALLPageData = strings.Replace(textALLPageData, "\n", "", -1)
	//TerTxt:=SplitLines(strings.TrimSpace(textALLPageData))
	TerTxt := SplitSubWarpWidth(strings.TrimSpace(textALLPageData), 25)

	re := regexp.MustCompile(`<a.+?href=\"(.+?)\".*>(.+)</a>`)
	matches := re.FindStringSubmatch(textDataArray[nextPageIndex])
	nextPageUrl := matches[1]
	//fmt.Println(nextPageUrl)

	//fmt.Println(len(strings.Split(txt,"<br/>")))
	//fmt.Println(strings.Index(txt, "下一章"))
	return PageDataAndNextChapter{TerTxt, nextPageUrl,jumpPage}
}

func handleEvents() {
	// http://www.boquku.com/book/4622/3084968.html

	//<div id="txtContent"> count = 1

	//  </div> count = 1

	// 下一章 <a href="/book/4622/3084969.html">下一章</a><span class="red">（快捷键:→）</span>
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
			r.Load(pdc.Data, pdc.NextChapter,pdc.CurrChapter)
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
		tmp, _ := ioutil.ReadFile("/tmp/r_txt/pos")
		urlPath = string(tmp)
	} else {
		urlPath = os.Args[1]

	}

	r = xreader.XReader(xreader.NewTxtReader())

	///book/4622/3084968.html

	pdc := getTxt(urlPath)
	if err := r.Load(pdc.Data, pdc.NextChapter,urlPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go start_server();
	Init(r)

}
