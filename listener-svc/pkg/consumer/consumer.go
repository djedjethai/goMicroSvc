package consumer

type Consumer interface {
}

type consumer struct{}

func NewConsumer() Consumer {
	return &consumer{}
}
