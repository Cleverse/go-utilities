/*
# Example

Basic usage:

	q := New[string]()
	q.Enqueue("Foo")

	item, ok := q.Dequeue()
	if !ok {
	  panic("item should be dequeued")
	}
	fmt.Println(item) // Output: Foo

Concurrency usage:

		q := New[string]()

		go func() {
		  for {
		    item, ok := q.Dequeue()
		    if !ok {
		      break
		    }
		    fmt.Println(item)
		  }
		}()

		data := []string{"Foo", "Bar", "Baz"}
		for _, item := range data {
		  q.Enqueue(item)
		}

	  q.Close() // close queue to stop goroutine

Queue is unlimited capacity, so you can enqueue as many as you want without blocking or dequeue required:

	q := New[string]()
	for i := 0; i < 1000000; i++ {
	  q.Enqueue("Foo")
	}
*/
package queue
