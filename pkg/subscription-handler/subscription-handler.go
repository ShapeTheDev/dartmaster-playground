package subscriptionhandler

import (
	"sync"
)

type SubscriptionHandler[T any] struct {
	mu            sync.Mutex
	subscriptions []chan T
}

// NewSubscriptionHandler returns a new subscription handler
func NewSubscriptionHandler[T any]() *SubscriptionHandler[T] {
	return &SubscriptionHandler[T]{}
}

// Subscribe creates a new subscriber channel and returns it.
func (sh *SubscriptionHandler[T]) Subscribe() <-chan T {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	return sh.subscribe()
}

// Unsubscribe removes the given subscriber channel and closes it.
func (sh *SubscriptionHandler[T]) Unsubscribe(logCh <-chan T) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.unsubscribe(logCh)
}

// UnsubscribeAll unsubscribes all subscriber channels and closes them.
func (sh *SubscriptionHandler[T]) UnsubscribeAll() {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.unsubscribeAll()
}

// Publish send the given item to all subscriber channels. If a subscriber is not reading its channel, it may be unsubscribed and no loger notified.
func (sh *SubscriptionHandler[T]) Publish(item T) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.publish(item)
}

// Subscriptions return the current amount of subscriptions.
func (sh *SubscriptionHandler[T]) Subscriptions() int {
	return len(sh.subscriptions)
}

func (sh *SubscriptionHandler[T]) subscribe() <-chan T {
	logCh := make(chan T, 10) // cache up to 10 item before the subscriber channel blocks and get unsubscribed
	sh.subscriptions = append(sh.subscriptions, logCh)
	return logCh
}

func (sh *SubscriptionHandler[T]) unsubscribe(logCh <-chan T) {
	var newsubs []chan T
	for _, curCh := range sh.subscriptions {
		if curCh == logCh {
			close(curCh)
		} else {
			newsubs = append(newsubs, curCh)
		}
	}
	sh.subscriptions = newsubs
}

func (sh *SubscriptionHandler[T]) unsubscribeAll() {
	for _, curCh := range sh.subscriptions {
		close(curCh)
	}
	sh.subscriptions = nil
}

func (sh *SubscriptionHandler[T]) publish(item T) {
	var unsubscribe []chan T
	for _, curCh := range sh.subscriptions {
		select {
		case curCh <- item:
		default:
			unsubscribe = append(unsubscribe, curCh)
		}
	}
	for _, blockingCh := range unsubscribe {
		sh.unsubscribe(blockingCh)
	}
}
