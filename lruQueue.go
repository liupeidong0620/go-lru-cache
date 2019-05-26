package lrucache

import (
	"time"
	"unsafe"
)

type lruCacheQueue struct {
	expire     time.Time
	next, prev *lruCacheQueue
}

func queueInit(size int64) *lruCacheQueue {
	var prev *lruCacheQueue
	var last *lruCacheQueue
	var next *lruCacheQueue
	var i int64

	if size < 0 {
		return nil
	}

	q := make([]lruCacheQueue, size+1, size+1)
	if size == 0 {
		q[0].next = &(q[0])
		q[0].prev = &(q[0])
	} else {
		prev = &(q[0])
		for i = 1; i <= size; i++ {
			next = &(q[i])
			prev.next = next
			next.prev = prev
			prev = next
		}
		last = &(q[size])
		last.next = &(q[0])
		q[0].prev = last
	}

	return &(q[0])
}

func queueIsEmpty(head *lruCacheQueue) bool {
	return unsafe.Pointer(head) == unsafe.Pointer(head.next)
}

func queueRemove(node *lruCacheQueue) {
	var prev *lruCacheQueue
	var next *lruCacheQueue

	next = node.next
	prev = node.prev

	next.prev = prev
	prev.next = next

	node.next = nil
	node.prev = nil
}

func queueInsertHead(head, node *lruCacheQueue) {
	var next *lruCacheQueue
	next = head.next

	next.prev = node
	node.next = next

	node.prev = head
	head.next = node
}

func queueInsertTail(head, node *lruCacheQueue) {
	var tail *lruCacheQueue
	var next *lruCacheQueue

	tail = head.prev

	tail.next = node
	node.prev = tail

	node.next = head
	head.prev = next
}

func queueLast(head *lruCacheQueue) *lruCacheQueue {
	return head.prev
}

func queueHead(head *lruCacheQueue) *lruCacheQueue {
	return head.next
}
