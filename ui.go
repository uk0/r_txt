package main
import (
	"log"
	"r_txt/xreader"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)


const (
	fakeShell  = "(base) ➜ ~"
	menuText   = "[↑]:j [↓]:k [b]:Boss Key [c]:下一章 [q]:exit [g]:Show/Hide Grid, [q]:Quit"
	fixedWidth = 80
	interval   = 3500
)

var (
	p      *widgets.Paragraph
	r      xreader.XReader
	ticker *time.Ticker

	showBorder   = false
	showHelp     = false
	showProgress = false
	bossKey      = false
	timer        = false
	rowNumber    = ""
	color        = 0
)

func setTimer() {
	timer = !timer

	if timer {
		ticker = time.NewTicker(interval * time.Millisecond)
		go func() {
			for range ticker.C {
				p.Text = r.Next()
				ui.Render(p)
			}
		}()
	} else {
		ticker.Stop()
	}
}

func updateParagraph(key string) {
	p.Text = key
}

func switchColor() {
	p.TextStyle.Fg = ui.Color(color % 8)
}

func displayHelp(current string) {
	showHelp = !showHelp
	if showHelp {
		p.Text = menuText
	} else {
		p.Text = current
	}
}

func displayBorder() {
	showBorder = !showBorder
	p.Border = showBorder
}

func displayProgress(current, progress string) {
	showProgress = !showProgress
	if showProgress {
		p.Text = progress
	} else {
		p.Text = current
	}
}

func displayBossKey(current string) {
	bossKey = !bossKey
	if bossKey {
		p.Border = false
		p.Text = fakeShell
	} else {
		p.Text = current
	}
}

// Init ui & components
func Init(gr xreader.XReader) {
	r = gr

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize the termui: %v", err)
	}

	//_, height := ui.TerminalDimensions()
	width := fixedWidth

	p = widgets.NewParagraph()
	p.Text = r.Current()
	p.SetRect(0, 0, width, 5)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	p.Border = showBorder

	ui.Render(p)
	handleEvents()
}
