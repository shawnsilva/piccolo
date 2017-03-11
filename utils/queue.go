package utils

import (
	"sync"
)

type (
	queueItem struct {
		data interface{}
		next *queueItem
	}

	// Queue contains pointers to start and end of a queue, the length, and a lock
	Queue struct {
		start  *queueItem
		end    *queueItem
		length int
		lock   *sync.Mutex
	}
)

// NewQueue will create a new queue for use.
func NewQueue() *Queue {
	q := &Queue{}
	q.lock = &sync.Mutex{}
	return q
}

// Length returns the length of a queue
func (q *Queue) Length() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.length
}

// Push adds an item to the end of the queue
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

// Pop will return the first item in the queue, and remove it.
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

// Look will return the first item in the queue, but not remove it.
func (q *Queue) Look() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := q.start
	if n == nil {
		return nil
	}
	return n.data
}
