package collect

type Unit struct {
	ID      string
	URL     string
	Payload []byte
}

func NewUnit(id, url string) Unit {
	return Unit{
		ID:  id,
		URL: url,
	}
}