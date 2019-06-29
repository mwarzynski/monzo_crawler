package queue

import "net/url"

type FIFO interface {
	Push(v url.URL)
	Pop() (url.URL, bool)
}
