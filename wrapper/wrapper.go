package wrapper

type wrapper struct {
	url string
}

func New(url string) *wrapper {
	w := wrapper{
		url: url,
	}
	return &w
}
