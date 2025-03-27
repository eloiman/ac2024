package utils

import (
	"errors"
	"log"
	"sync"
)

type sharedState[T any] struct {
	sync.Mutex
	size        uint32
	resultsLeft uint32
	result      chan T
	isOpen      bool
}

type Promise[T any] struct {
	state *sharedState[T]
}

type Future[T any] struct {
	state *sharedState[T]
}

type Result[T any] struct {
	result T
	err    error
}

func NewPromise[T any](n uint32) Promise[T] {
	return Promise[T]{state: &sharedState[T]{result: make(chan T, n), size: n, resultsLeft: n, isOpen: true}}
}

func (promise *Promise[T]) Put(value T) bool {
	promise.state.Lock()
	defer promise.state.Unlock()

	if promise.state.resultsLeft == 0 {
		log.Fatalln("promise has been already out of results to return")
		return false
	}
	promise.state.result <- value
	promise.state.resultsLeft--

	return true
}

func (promise *Promise[T]) Close() {
	promise.state.Lock()
	defer promise.state.Unlock()

	close(promise.state.result)
	promise.state.isOpen = false
	promise.state.resultsLeft = 0
}

func (promise *Promise[T]) Future() Future[T] {
	return Future[T]{state: promise.state}
}

func (future *Future[T]) IsEmpty() bool {
	future.state.Lock()
	defer future.state.Unlock()

	return len(future.state.result) == 0
}

func (state *sharedState[T]) read() T {
	value := <-state.result
	return value
}

func (future *Future[T]) Get() T {
	future.state.Lock()
	defer future.state.Unlock()

	if !future.state.isOpen {
		panic("It's closed")
	}

	return future.state.read()
}

func (future *Future[T]) GetOr(defaultValue T) T {
	future.state.Lock()
	defer future.state.Unlock()

	if !future.state.isOpen {
		return defaultValue
	}

	return future.state.read()
}

func (future *Future[T]) Result() Result[T] {
	future.state.Lock()
	defer future.state.Unlock()

	if !future.state.isOpen {
		return Result[T]{err: errors.New("state is closed")}
	}

	if len(future.state.result) == 0 {
		return Result[T]{err: errors.New("results are empty")}
	}

	return Result[T]{result: future.state.read(), err: nil}
}

func (future *Future[T]) WaitResult() Result[T] {
	future.state.Lock()
	defer future.state.Unlock()

	if !future.state.isOpen {
		return Result[T]{err: errors.New("state is closed")}
	}

	return Result[T]{result: future.state.read(), err: nil}
}

func (future *Future[T]) IsOpen() bool {
	future.state.Lock()
	defer future.state.Unlock()

	return future.state.isOpen
}

func (future *Future[T]) Close() {
	future.state.Lock()
	defer future.state.Unlock()

	close(future.state.result)
	future.state.isOpen = false
	future.state.resultsLeft = 0
}
