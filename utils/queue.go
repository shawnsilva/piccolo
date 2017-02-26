package utils

import (
	"sync"
)

type (
	queueItem struct {
		data interface{}
		next *queueItem
	}

	Queue struct {
		start  *queueItem
		end    *queueItem
		length int
		lock   *sync.Mutex
	}
)

func NewQueue() *Queue {
	q := &Queue{}
	q.lock = &sync.Mutex{}
	return q
}

func (q *Queue) Length() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.length
}

func (q *Queue) Push(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := &queueItem{data: item}

	if q.end == nil {
		q.end = n
		q.start = n
	} else {
		q.end.next = n
		q.end = n
	}
	q.length++
}

func (q *Queue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.start == nil {
		return nil
	}

	n := q.start
	q.start = n.next

	if q.start == nil {
		q.end = nil
	}
	q.length--

	return n.data
}

func (q *Queue) Look() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := q.start
	if n == nil {
		return nil
	}
	return n.data
}
