package queue

import (
	"sync"
)

// Queue a instance of thread-safe and unlimited-size generics in-memory message queue
// and alternative way to communicate between goroutines compared to channel.
type Queue[T comparable] struct {
	cond     *sync.Cond
	items    []T
	mu       sync.RWMutex
	isClosed bool
}

// New creates a new message queue.
func New[T comparable]() *Queue[T] {
	q := &Queue[T]{
		items: make([]T, 0),
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Close closes the queue.
// queue will be permanent unusable after this method is called.
func (q *Queue[T]) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.isClosed = true
	q.items = nil
	q.cond.Broadcast()
}

// IsClosed returns true if the queue is closed.
func (q *Queue[T]) IsClosed() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.isClosed
}

// Enqueue adds an item to the end of the queue. returns the index of the item.
func (q *Queue[T]) Enqueue(item T) (index int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isClosed {
		return -1
	}

	q.items = append(q.items, item)
	q.cond.Signal()

	return len(q.items) - 1
}

// Dequeue removes an item from the front of the queue.
//
// If the queue is empty, this method blocks(wait) until an item is available.
//
// returns second value as false if the queue is closed.
func (q *Queue[T]) Dequeue() (val T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 {
		if q.isClosed {
			return val, false
		}
		q.cond.Wait()
	}

	// Recheck the length of the queue after waking up from wait.
	// or another goroutine might have dequeued the last item.
	if q.isClosed || len(q.items) == 0 {
		return val, false
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item, true
}

// TryDequeue removes an item from the front of the queue, similar to Dequeue.
//
// However, TryDequeue does NOT block if the queue is empty. It returns second value as false immediately if the queue is empty or closed.
func (q *Queue[T]) TryDequeue() (val T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isClosed || len(q.items) == 0 {
		return val, false
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item, true
}

// IsEmpty returns true if the queue is empty.
func (q *Queue[T]) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items) == 0
}

// IndexOf returns the index of the first item that matches the target.
func (q *Queue[T]) IndexOf(target T) (index int) {
	return q.IndexOfIter(func(item T) bool {
		return item == target
	})
}

// IndexOfIter returns the index of the first item that matches the callback.
// If no item matches the callback, -1 is returned.
func (q *Queue[T]) IndexOfIter(cb func(item T) bool) (index int) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	for i, item := range q.items {
		if cb(item) {
			return i
		}
	}
	return -1
}

// Len returns the length of the queue.
func (q *Queue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items)
}

// Items returns a copy of the all items in the queue.
func (q *Queue[T]) Items() []T {
	q.mu.RLock()
	defer q.mu.RUnlock()
	items := make([]T, len(q.items))
	copy(items, q.items)
	return items
}

// RemoveAt removes an item at the given index and returns the new size of the queue.
func (q *Queue[T]) RemoveAt(index int) (size int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items[:index], q.items[index+1:]...)
	return len(q.items)
}

// Clear removes all items from the queue and returns the size of the removed items.
func (q *Queue[T]) Clear() (size int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	size = len(q.items)

	// clean up the slice without changing the capacity and allocation.
	q.items = q.items[:0:cap(q.items)]

	return size
}
