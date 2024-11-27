package queue

import (
	"sync"
)

type node[T any] struct {
	data T
	next *node[T]
}
type Queue[T any] struct {
	head      *node[T]
	tail      *node[T]
	headMutex sync.Mutex
	tailMutex sync.Mutex
	dataCond  *sync.Cond
}

func (q *Queue[T]) getTail() *node[T] {
	q.tailMutex.Lock()
	defer q.tailMutex.Unlock()
	return q.tail
}
func (q *Queue[T]) pop_head() *node[T] {
	old_head := q.head
	q.head = q.head.next
	return old_head
}
func (q *Queue[T]) Wait_pop() T {
	q.headMutex.Lock()
	defer q.headMutex.Unlock()
	if q.head != q.getTail() {
		return q.pop_head().data
	}
	q.dataCond.Wait()
	return q.pop_head().data
}

func (q *Queue[T]) Push(value T) {
	newNode := new(node[T])

	q.tailMutex.Lock()
	defer q.tailMutex.Unlock()
	q.tail.data = value
	q.tail.next = newNode

	q.tail = newNode
	q.dataCond.Signal()
}

func (q *Queue[T]) empty() bool {
	q.headMutex.Lock()
	defer q.headMutex.Unlock()
	return q.head == q.tail
}

func NewQueue[T any]() *Queue[T] {
	q := &Queue[T]{}
	q.dataCond = sync.NewCond(&q.headMutex)
	q.head = new(node[T])
	q.tail = q.head
	return q
}
