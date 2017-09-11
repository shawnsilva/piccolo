package utils

import (
	"sync"
	"testing"
)

type (
	dllNodes struct {
		name string
		data interface{}
	}
	dllTest struct {
		nodes  []dllNodes
		length int
	}
)

var (
	dllTestData = dllTest{
		[]dllNodes{
			{"1", "a"},
			{"2", "b"},
			{"3", "c"},
			{"4", "d"},
			{"5", "e"},
			{"6", "f"},
		},
		6,
	}
	dllTestData2 = dllTest{
		[]dllNodes{
			{"7", "g"},
			{"8", "h"},
			{"9", "i"},
			{"10", "j"},
			{"11", "k"},
			{"12", "l"},
		},
		6,
	}
	dllTestData3 = dllTest{
		[]dllNodes{
			{"13", "m"},
			{"14", "n"},
			{"15", "o"},
			{"16", "p"},
			{"17", "q"},
			{"18", "r"},
		},
		6,
	}
	dllTestData4 = dllTest{
		[]dllNodes{
			{"19", "s"},
			{"20", "t"},
			{"21", "u"},
			{"22", "v"},
			{"23", "w"},
			{"24", "x"},
		},
		6,
	}
)

func TestDoubleLinkedListOrder(t *testing.T) {
	dll := NewDoubleLinkedList()
	for _, dllItem := range dllTestData.nodes {
		nNode := NewNode(dllItem.name, dllItem.data)
		dll.InsertEnd(nNode)
	}
	node := dll.First()
	for i := 0; i < dllTestData.length; i++ {
		nName, nData := node.GetData()
		nDataStr, ok := nData.(string)
		if !ok {
			t.Error(
				"Couldn't cast data to original type: ", nData,
			)
			continue
		}
		if nName != dllTestData.nodes[i].name {
			t.Error(
				"Node Names dont match: ", nName, ", ", dllTestData.nodes[i].name,
			)
			continue
		}
		if nDataStr != dllTestData.nodes[i].data {
			t.Error(
				"Node Data doesnt match: ", nDataStr, ", ", dllTestData.nodes[i].data,
			)
			continue
		}
		node = node.Next()
	}
}

func TestDoubleLinkedListConcurrenyLength(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(4)
	dll := NewDoubleLinkedList()
	go func() {
		defer wg.Done()
		for _, dllItem := range dllTestData.nodes {
			nNode := NewNode(dllItem.name, dllItem.data)
			dll.InsertBeginning(nNode)
		}
	}()
	go func() {
		defer wg.Done()
		for _, dllItem := range dllTestData2.nodes {
			nNode := NewNode(dllItem.name, dllItem.data)
			dll.InsertEnd(nNode)
		}
	}()
	go func() {
		defer wg.Done()
		for _, dllItem := range dllTestData3.nodes {
			nNode := NewNode(dllItem.name, dllItem.data)
			dll.InsertBeginning(nNode)
		}
	}()
	go func() {
		defer wg.Done()
		for _, dllItem := range dllTestData4.nodes {
			nNode := NewNode(dllItem.name, dllItem.data)
			dll.InsertEnd(nNode)
		}
	}()
	wg.Wait()
	totalLength := dllTestData.length + dllTestData2.length + dllTestData3.length + dllTestData4.length
	if dll.Length() != totalLength {
		t.Error(
			"New List Length of: ", dll.Length(),
			"Should be: ", totalLength,
		)
	}
}

func TestDoubleLinkedListLength(t *testing.T) {
	dll := NewDoubleLinkedList()
	for _, dllItem := range dllTestData.nodes {
		nNode := NewNode(dllItem.name, dllItem.data)
		dll.InsertBeginning(nNode)
	}
	for _, dllItem := range dllTestData2.nodes {
		nNode := NewNode(dllItem.name, dllItem.data)
		dll.InsertEnd(nNode)
	}
	for _, dllItem := range dllTestData3.nodes {
		nNode := NewNode(dllItem.name, dllItem.data)
		dll.InsertBefore(dll.First().Next().Next(), nNode)
	}
	for _, dllItem := range dllTestData4.nodes {
		nNode := NewNode(dllItem.name, dllItem.data)
		dll.InsertAfter(dll.Last().Prev().Prev(), nNode)
	}
	totalLength := dllTestData.length + dllTestData2.length + dllTestData3.length + dllTestData4.length
	if dll.Length() != totalLength {
		t.Error(
			"New List Length of: ", dll.Length(),
			"Should be: ", totalLength,
		)
	}
	// Remove a total of 3 nodes, one inthe middle, head and tail.
	dll.Delete("7")
	dll.Remove(dll.First())
	dll.Remove(dll.Last())
	if dll.Length() != totalLength-3 {
		t.Error(
			"List with removed Length of: ", dll.Length(),
			"Should be: ", totalLength-3,
		)
	}
}
