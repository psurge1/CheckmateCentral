// Package datastructures contains implementations of various datastructures to assist with matchmaking algorithms
// datastructures explicitly named [Name]Concurrent or Concurrent[Name] are concurrent safe. Otherwise, assume they are not
package datastructures

import (
	"errors"
	"slices"
)

type Queue[T any] struct {
	values []T
}

func (q *Queue[T]) Enqueue(value T) {
	q.values = append(q.values, value)
}

func (q *Queue[T]) Dequeue() (T, error) {
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
	if len(q.values) < 1 {
		var blank T
		return blank, errors.New("queue is empty")
	}
	return q.values[0], nil
}

func (q *Queue[T]) FindAndDelete(object T, isEqual func(T, T) bool) (T, error) {
	targetIdx := -1
	for idx, value := range q.values {
		if isEqual(object, value) {
			targetIdx = idx
			break
		}
	}
	if targetIdx == -1 {
		var blank T
		return blank, errors.New("unable to find object")
	}

	value := q.values[targetIdx]
	q.values = slices.Delete(q.values, targetIdx, targetIdx+1)
	return value, nil
}

func (q *Queue[T]) Len() int {
	return len(q.values)
}
