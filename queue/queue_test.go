package queue

import (
	"sync"
	"testing"
	"time"
)

func TestEnqueueDequeue(t *testing.T) {
	expecteds := []string{"a", "b", "c", "d", "e"}

	q := New[string]()
	for expectedIndex, v := range expecteds {
		actualIndex := q.Enqueue(v)
		if expectedIndex != actualIndex {
			t.Errorf("index should be equal to the expected value, expected: %d, actual: %d", expectedIndex, actualIndex)
		}
	}

	for _, expected := range expecteds {
		actual, ok := q.Dequeue()
		if !ok {
			t.Errorf("item should be dequeued")
		}
		if expected != actual {
			t.Errorf("item should be equal to the expected value, expected: %s, actual: %s", expected, actual)
		}
	}
}

func TestAsyncEnqueueDequeue(t *testing.T) {
	expected := "A"
	q := New[string]()
	go func() {
		time.Sleep(3 * time.Second)
		q.Enqueue(expected)
	}()

	// expected to block until enqueue
	actual, ok := q.Dequeue()
	if !ok {
		t.Errorf("item should be dequeued")
	}
	if expected != actual {
		t.Errorf("item should be equal to the expected value, expected: %s, actual: %s", expected, actual)
	}
}

func TestCloseConcurrency(t *testing.T) {
	var (
		q  = New[int]()
		wg sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, ok := q.Dequeue(); ok {
			t.Error("should not dequeue from an empty queue")
		}
	}()

	time.Sleep(1 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		q.Close()
	}()
	wg.Wait()

	actualIndex := q.Enqueue(1)
	actualItem, ok := q.Dequeue()

	if actualIndex != -1 {
		t.Errorf("should not enqueue to a closed queue")
	}
	if actualItem != 0 {
		t.Errorf("should not dequeue from a closed queue")
	}
	if ok {
		t.Errorf("should not dequeue from a closed queue")
	}
}

func TestClear(t *testing.T) {
	q := New[int]()

	expectedSize := 10
	for i := 0; i < expectedSize; i++ {
		q.Enqueue(i)
	}

	itemsCap := cap(q.items)
	actualSize := q.Clear()

	if q.Len() != 0 {
		t.Errorf("queue size should be zero")
	}
	if expectedSize != actualSize {
		t.Errorf("flushed size should be equal to the enqueued size")
	}
	if itemsCap != cap(q.items) {
		t.Errorf("capacity should not change after flush")
	}
}

func TestRemoveAt(t *testing.T) {
	q := New[int]()
	size := 10
	expected := make([]int, 0, size)

	for i := 0; i < size; i++ {
		value := i + 1
		expected = append(expected, value)
		q.Enqueue(value)
	}

	removeValue := 5
	removeIndex := q.IndexOf(removeValue)
	newSize := q.RemoveAt(removeIndex)

	// assert.Equal(t, size-1, newSize, "new size should be equal to the old size minus 1")
	// assert.Equal(t, -1, q.IndexOf(removedValue), "index of removed item should be -1")
	if size-1 != newSize {
		t.Errorf("new size should be equal to the old size minus 1")
	}
	if -1 != q.IndexOf(removeValue) {
		t.Errorf("index of removed item should be -1")
	}

	for i := 0; i < q.Len(); i++ {
		item, ok := q.Dequeue()
		if !ok {
			t.Errorf("item should be dequeued")
		}
		if i != removeIndex && expected[i] != item {
			t.Errorf("item should be equal to the expected value, expected: %d, actual: %d", expected[i], item)
		}
		if removeValue == item {
			t.Errorf("removed item should not be in the queue")
		}
	}
}
