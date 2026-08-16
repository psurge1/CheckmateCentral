package datastructures

import (
	"errors"
	"sync"
)

type Queue[T any] struct {
	mu     sync.Mutex
	values []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{values: make([]T, 0)}
}

func (q *Queue[T]) Enqueue(value T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.values = append(q.values, value)
}

func (q *Queue[T]) Dequeue() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.values) < 1 {
		var blank T
		return blank, errors.New("queue is empty")
	}

	item := q.values[0]

	var blank T
	q.values[0] = blank
	q.values = q.values[1:]
	return item, nil
}

func (q *Queue[T]) Peek() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.values) < 1 {
		var blank T
		return blank, errors.New("queue is empty")
	}
	return q.values[0], nil
}

func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.values)
}
