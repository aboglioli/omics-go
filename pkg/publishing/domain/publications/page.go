package publications

type Point struct {
	X int
	Y int
}

type Size struct {
	Width  int
	Height int
}

type Frame struct {
	Order    int
	Position Point
	Size     Size
}

type Image struct {
	URL    string
	Frames []Frame
}

type Page struct {
	Number int
	Images []Image
}
