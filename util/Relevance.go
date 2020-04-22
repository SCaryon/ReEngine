package utils

type Relevance struct {
	*Article
	Weight		float64
	TitleSegs	[]string
	ContentSegs	[]string
}

type DocSlice []Relevance

func (a DocSlice) Len() int {
	return len(a)
}

func (a DocSlice) Swap(i,j int) {
	a[i], a[j] = a[j], a[i]
}

func (a DocSlice) Less(i,j int) bool {
	return a[i].Weight > a[j].Weight
}