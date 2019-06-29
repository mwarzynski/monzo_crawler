package queue

import (
	"errors"
	"net/url"
)

var ErrEmpty = errors.New("queue is empty")

type FIFOSlice struct {
	queue []url.URL
}

func NewFIFOSlice(initLength uint) *FIFOSlice {
	return &FIFOSlice{
		queue: make([]url.URL, 0, int(initLength)),
	}
}

func (q *FIFOSlice) Push(v url.URL) {
	q.queue = append(q.queue, v)
}

func (q *FIFOSlice) Pop() (url.URL, bool) {
	if len(q.queue) == 0 {
		return url.URL{}, false
	}
	v := q.queue[0]
	q.queue = q.queue[1:]
	return v, true
}
