package main

import (
	"fmt"
	"sync"
)

// PubSub - pub/sub pattern
type PubSub[T any] struct {
	Subscribers []chan T
	mu          sync.RWMutex
	closed      bool
}

// NewPubSub - creates a new PubSub
func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		mu: sync.RWMutex{},
	}
}

// Subscribe - creates a new subscriber
func (s *PubSub[T]) Subscribe() <-chan T {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	r := make(chan T)

	s.Subscribers = append(s.Subscribers, r)

	return r
}

// Publish - sends the value to all subscribers
func (s *PubSub[T]) Publish(value T) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return
	}

	for _, ch := range s.Subscribers {
		ch <- value
	}
}

// Close - closes all subscribers
func (s *PubSub[T]) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	for _, ch := range s.Subscribers {
		close(ch)
	}

	s.closed = true
}

func main() {
	ps := NewPubSub[string]()

	wg := sync.WaitGroup{}

	s1 := ps.Subscribe()

	wg.Add(1)

	go func() {
		defer wg.Done()
		for val := range s1 {
			fmt.Println("Subscriber 1's value:", val)
		}

		fmt.Print("Subscriber 1, exiting\n")
	}()

	s2 := ps.Subscribe()

	wg.Add(1)

	go func() {
		defer wg.Done()
		for val := range s2 {
			fmt.Println("Subscriber 2's value:", val)
		}

		fmt.Print("Subscriber 2, exiting\n")
	}()

	ps.Publish("one")
	ps.Publish("two")
	ps.Publish("three")

	ps.Close()

	wg.Wait()

	fmt.Println("Done")
}
