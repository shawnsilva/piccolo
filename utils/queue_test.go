package utils

import (
	"sync"
	"testing"
)

type (
	queueTest struct {
		data   []string
		length int
	}
)

var (
	queueTestDataStrings = queueTest{
		[]string{
			"test",
			"some",
			"queues",
			"1",
			"42",
			"hello",
			"maybe",
			"blah",
			"another",
			"last",
		},
		10,
	}

	queueTestDataStrings2 = queueTest{
		[]string{
			"ztest",
			"sascome",
			"queafues",
			"1435",
			"42twser",
			"helgsdgfslo",
			"maycxvbxcvbe",
			"blaxcbvxcbh",
			"anothexcbvxr",
			"last",
			"nope",
			"more",
			"asdaufgh",
			"why",
			"another",
		},
		15,
	}
)

func TestQueueOrder(t *testing.T) {
	queue := NewQueue()
	for _, qItem := range queueTestDataStrings.data {
		queue.Push(qItem)
	}
	for _, qItem := range queueTestDataStrings.data {
		poppedData := queue.Pop()
		if poppedData == nil {
			t.Error(
				"Expected another item: ", qItem,
				"Got nil",
			)
		}
		poppedItem, ok := poppedData.(string)
		if !ok {
			t.Error(
				"Couldn't cast data to original type: ", poppedItem,
			)
			continue
		}
		if poppedItem != qItem {
			t.Error(
				"Expected item: ", qItem,
				"Got: ", poppedItem,
			)
		}
	}
}

func TestQueueLength(t *testing.T) {
	queue := NewQueue()
	for _, qItem := range queueTestDataStrings.data {
		queue.Push(qItem)
	}
	for _, qItem := range queueTestDataStrings2.data {
		queue.Push(qItem)
	}

	totalLength := queueTestDataStrings.length + queueTestDataStrings2.length

	if queue.Length() != totalLength {
		t.Error(
			"New Queue Length of: ", queue.Length(),
			"Should be: ", totalLength,
		)
	}
	// Pop 3 items and check length again
	queue.Pop()
	queue.Pop()
	queue.Pop()
	if queue.Length() != totalLength-3 {
		t.Error(
			"New Queue Length of: ", queue.Length(),
			"Should be: ", totalLength-3,
		)
	}
}

func TestQueueConcurrencyLength(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	queue := NewQueue()
	go func() {
		defer wg.Done()
		for _, qItem := range queueTestDataStrings.data {
			queue.Push(qItem)
		}
	}()
	go func() {
		defer wg.Done()
		for _, qItem := range queueTestDataStrings2.data {
			queue.Push(qItem)
		}
	}()
	wg.Wait()

	totalLength := queueTestDataStrings.length + queueTestDataStrings2.length

	if queue.Length() != totalLength {
		t.Error(
			"New Queue Length of: ", queue.Length(),
			"Should be: ", totalLength,
		)
	}

	wg.Add(4)
	go func() {
		defer wg.Done()
		queue.Pop()
	}()
	go func() {
		defer wg.Done()
		queue.Pop()
	}()
	go func() {
		defer wg.Done()
		queue.Pop()
	}()
	go func() {
		defer wg.Done()
		queue.Pop()
	}()
	wg.Wait()

	if queue.Length() != totalLength-4 {
		t.Error(
			"New Queue Length of: ", queue.Length(),
			"Should be: ", totalLength-4,
		)
	}
}
