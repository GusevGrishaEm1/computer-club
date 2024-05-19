package util

type Node struct {
	Value interface{}
	Next  *Node
}

type Queue struct {
	Head *Node
	Tail *Node
	Size int
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Enqueue(value interface{}) {
	node := &Node{Value: value}
	if q.Head == nil {
		q.Head = node
		q.Tail = node
	} else {
		q.Tail.Next = node
		q.Tail = node
	}
	q.Size++
}

func (q *Queue) Dequeue() (interface{}, bool) {
	if q.Head == nil {
		return nil, false
	}
	value := q.Head.Value
	q.Head = q.Head.Next
	if q.Head == nil {
		q.Tail = nil
	}
	q.Size--
	return value, true
}
