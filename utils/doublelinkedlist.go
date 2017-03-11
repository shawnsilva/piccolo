package utils

import (
	"sync"
)

type (
	data struct {
		name string
		data interface{}
	}
	// Node contains links to next and previous node, and data of current node
	Node struct {
		data
		lock *sync.Mutex
		next *Node
		prev *Node
	}
	// DoubleLinkedList implements a double linked list data structure
	DoubleLinkedList struct {
		lock   *sync.Mutex
		length int
		head   *Node
		tail   *Node
	}
)

// NewDoubleLinkedList returns a new empty DoubleLinkedList
func NewDoubleLinkedList() *DoubleLinkedList {
	l := &DoubleLinkedList{}
	l.lock = &sync.Mutex{}
	l.length = 0
	return l
}

// NewNode creates a new node for a DoubleLinkedList
func NewNode(dataName string, dataContents interface{}) *Node {
	n := &Node{}
	n.lock = &sync.Mutex{}
	n.SetData(dataName, dataContents)
	return n
}

// GetData will return the data stored in the node
func (n *Node) GetData() (string, interface{}) {
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.data.name, n.data.data
}

// SetData will add data to the node
func (n *Node) SetData(dataName string, dataContents interface{}) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.data = data{dataName, dataContents}
}

// Next returns the next node
func (n *Node) Next() *Node {
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.next
}

// Prev returns the previous node
func (n *Node) Prev() *Node {
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.prev
}

// First returns the first node in the double linked list
func (l *DoubleLinkedList) First() *Node {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.head
}

// Last returns the last node in the double linked list
func (l *DoubleLinkedList) Last() *Node {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.tail
}

// InsertBefore inserts a new node before a node
func (l *DoubleLinkedList) InsertBefore(node *Node, newNode *Node) {
	l.lock.Lock()
	defer l.lock.Unlock()

	newNode.lock.Lock()
	defer newNode.lock.Unlock()
	node.lock.Lock()
	defer node.lock.Unlock()

	newNode.next = node
	if node.prev == nil {
		newNode.prev = nil
		l.head = newNode
	} else {
		node.prev.lock.Lock()
		defer node.prev.lock.Unlock()
		newNode.prev = node.prev
		node.prev.next = newNode
	}
	node.prev = newNode
	l.length++
}

// InsertAfter inserts a new node after a node
func (l *DoubleLinkedList) InsertAfter(node *Node, newNode *Node) {
	l.lock.Lock()
	defer l.lock.Unlock()

	newNode.lock.Lock()
	defer newNode.lock.Unlock()
	node.lock.Lock()
	defer node.lock.Unlock()

	newNode.prev = node
	if node.next == nil {
		newNode.next = nil
		l.tail = newNode
	} else {
		node.next.lock.Lock()
		defer node.next.lock.Unlock()
		newNode.next = node.next
		node.next.prev = newNode
	}
	node.next = newNode
	l.length++
}

// InsertBeginning inserts a new node at the beginning
func (l *DoubleLinkedList) InsertBeginning(newNode *Node) {
	if l.head == nil {
		l.lock.Lock()
		defer l.lock.Unlock()
		newNode.lock.Lock()
		defer newNode.lock.Unlock()
		l.head = newNode
		l.tail = newNode
		newNode.prev = nil
		newNode.next = nil
		l.length++
	} else {
		l.InsertBefore(l.head, newNode)
	}
}

// InsertEnd insterts a new node at the end
func (l *DoubleLinkedList) InsertEnd(newNode *Node) {
	if l.tail == nil {
		l.InsertBeginning(newNode)
	} else {
		l.InsertAfter(l.tail, newNode)
	}
}

// Remove removes a specific node from the list
func (l *DoubleLinkedList) Remove(node *Node) {
	l.lock.Lock()
	defer l.lock.Unlock()

	node.lock.Lock()
	defer node.lock.Unlock()

	if node.prev == nil {
		l.head = node.next
	} else {
		node.prev.lock.Lock()
		defer node.prev.lock.Unlock()
		node.prev.next = node.next
	}

	if node.next == nil {
		l.tail = node.prev
	} else {
		node.next.lock.Lock()
		defer node.next.lock.Unlock()
		node.next.prev = node.prev
	}

	l.length--
}

// Find finds a specific node based on the name of the data, iterates until found
// or checked all nodes
func (l *DoubleLinkedList) Find(name string) *Node {
	var foundNode *Node
	for curNode := l.First(); curNode != nil; curNode = curNode.Next() {
		if curNode.data.name == name {
			foundNode = curNode
			break
		}
	}
	return foundNode
}

// Delete will delete a node based on the name of the data
func (l *DoubleLinkedList) Delete(name string) {
	node := l.Find(name)
	if node != nil {
		l.Remove(node)
	}
}
