package xreader

type XReader interface {
	Load(path []string,prevChapter string,nextChapter string,currentChapter string) error
	Current() string
	Next() string
	Prev() string
	First() string
	Last() string
	CurrentPos() int
	Goto(pos int) string
	GetProgress() string
	GetNextChapter()string
	GetPrevChapter() string
	GetCurrentChapter()string
}
