package utils

import (
	"sync"
)

type (
	queueItem struct {
		data interface{}
		next *queueItem
		lock *sync.Mutex
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
	n.lock = &sync.Mutex{}

	n.lock.Lock()
	defer n.lock.Unlock()

	if q.end == nil {
		q.end = n
		q.start = n
	} else {
		q.end.lock.Lock()
		defer q.end.lock.Unlock()
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

	q.start.lock.Lock()
	defer q.start.lock.Unlock()

	n := q.start
	q.start = n.next

	if q.start == nil {
		q.end = nil
	}
	q.length--

	return n.data
}

// Look will return the first item's data in the queue, but not remove it.
func (q *Queue) Look() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := q.start
	if n == nil {
		return nil
	}
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.data
}

// First returns the first queueItem.
func (q *Queue) First() *queueItem {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.start
}

func (i *queueItem) Next() *queueItem {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.next
}

func (i *queueItem) Data() interface{} {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.data
}
