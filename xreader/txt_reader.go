package xreader

import (
	"fmt"
)

type TxtReader struct {
	content []string
	pos     int
	nextChapter string
	prevChapter string
	currentChapter string
}

func NewTxtReader() *TxtReader {
	return &TxtReader{}
}

func (txt *TxtReader) Load(string2 []string,prevChapter string,nextChapter string,currentChapter string) error {
	txt.content = string2
	txt.nextChapter = nextChapter
	txt.prevChapter = prevChapter
	txt.currentChapter = currentChapter
	txt.pos = 0
	return nil
}

func (txt *TxtReader)GetNextChapter()string {
	return txt.nextChapter;
}


func (txt *TxtReader)GetPrevChapter()string {
	return txt.prevChapter;
}


func (txt *TxtReader)GetCurrentChapter()string  {
	return txt.currentChapter;
}


func (txt *TxtReader) Current() string {
	return txt.content[txt.pos]
}

func (txt *TxtReader) Next() string {
	txt.pos++

	if txt.pos <= len(txt.content)-1 {
		return txt.content[txt.pos]
	} else {
		txt.pos = len(txt.content) - 1
	}

	return "END Line"
}

func (txt *TxtReader) Prev() string {
	txt.pos--

	if txt.pos < 0 {
		txt.pos = 0
	}

	return txt.content[txt.pos]
}

func (txt *TxtReader) First() string {
	txt.pos = 0
	return txt.content[0]
}

func (txt *TxtReader) Last() string {
	txt.pos = len(txt.content) - 1
	return txt.content[len(txt.content)-1]
}

func (txt *TxtReader) CurrentPos() int {
	return txt.pos
}

func (txt *TxtReader) Goto(pos int) string {
	if pos < 0 {
		pos = 0
	}

	if pos > len(txt.content)-1 {
		pos = len(txt.content) - 1
	}

	txt.pos = pos
	return txt.content[txt.pos]
}

func (txt *TxtReader) GetProgress() string {
	return fmt.Sprintf("(%d / %d)", txt.pos+1, len(txt.content))
}
