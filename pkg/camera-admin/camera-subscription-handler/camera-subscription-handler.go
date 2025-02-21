package camerasubscriptionhandler

import (
	"sync"
)

type CameraSubscriptionHandler[T any] struct {
	mu            sync.Mutex
	subscriptions []chan T
}

// newCameraSubscriptionHandler returns a new camera-subscription-handler.
func newCameraSubscriptionHandler[T any]() *CameraSubscriptionHandler[T] {
	return &CameraSubscriptionHandler[T]{}
}

// Subscribe creates a new subscriber channel and returns it.
func (sh *CameraSubscriptionHandler[T]) Subscribe() <-chan T {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	return sh.subscribe()
}

// Unsubscribe removes the given subscriber channel and closes it.
func (sh *CameraSubscriptionHandler[T]) Unsubscribe(logCh <-chan T) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.unsubscribe(logCh)
}

// UnsubscribeAll unsubscribes all subscriber channels and closes them.
func (sh *CameraSubscriptionHandler[T]) UnsubscribeAll() {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.unsubscribeAll()
}

// Publish sends the given item to all subscriber channels. If a subscriber is not reading its channel, it may be unsubscribe and no longer notified.
func (sh *CameraSubscriptionHandler[T]) Publish(item T) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.publish(item)
}

// Subscriptions return the current amount of subscriptions.
func (sh *CameraSubscriptionHandler[T]) Subscriptions() int {
	return len(sh.subscriptions)
}

func (sh *CameraSubscriptionHandler[T]) subscribe() <-chan T {
	logCh := make(chan T) // cache up to 0 items before the subscriber channel blocks --> improvement of frame pushing
	sh.subscriptions = append(sh.subscriptions, logCh)
	return logCh
}

func (sh *CameraSubscriptionHandler[T]) unsubscribe(logCh <-chan T) {
	var newSubs []chan T
	for _, curCh := range sh.subscriptions {
		if curCh == logCh {
			close(curCh)
		} else {
			newSubs = append(newSubs, curCh)
		}
	}
	sh.subscriptions = newSubs
}

func (sh *CameraSubscriptionHandler[T]) unsubscribeAll() {
	for _, curCh := range sh.subscriptions {
		close(curCh)
	}
	sh.subscriptions = nil
}

func (sh *CameraSubscriptionHandler[T]) publish(item T) {
	for _, curCh := range sh.subscriptions {
		select {
		case curCh <- item:
			// --> channel is able to receive item
		default:
			// --> channel is not able to receive item --> channel is blocked
		forLoop:
			for {
				select {
				case <-curCh:
					// --> unblock the channel
				default:
					// --> channel is unblocked --> break
					break forLoop
				}
			}
		}
	}
}
