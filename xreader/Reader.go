package xreader

type XReader interface {
	Load(path []string,nextChapter string,currentChapter string) error
	Current() string
	Next() string
	Prev() string
	First() string
	Last() string
	CurrentPos() int
	Goto(pos int) string
	GetProgress() string
	//GoNextChapter()
	GetNextChapter()string
	GetCurrentChapter()string
}
