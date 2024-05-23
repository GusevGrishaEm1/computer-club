package service

import (
	"testing"
)

func TestQueue(t *testing.T) {
	queue := NewQueue()

	// Test enqueue
	queue.Enqueue("item1")
	queue.Enqueue("item2")
	if queue.Size != 2 {
		t.Errorf("Expected size 2, got %d", queue.Size)
	}

	// Test dequeue
	item, ok := queue.Dequeue()
	if !ok || item != "item1" {
		t.Errorf("Expected item1, got %v", item)
	}
	if queue.Size != 1 {
		t.Errorf("Expected size 1, got %d", queue.Size)
	}

	// Test empty queue
	item, ok = queue.Dequeue()
	if !ok || item != "item2" {
		t.Errorf("Expected item2, got %v", item)
	}
	if queue.Size != 0 {
		t.Errorf("Expected size 0, got %d", queue.Size)
	}

	item, ok = queue.Dequeue()
	if ok {
		t.Errorf("Expected nil item, got %v", item)
	}
	if queue.Size != 0 {
		t.Errorf("Expected size 0, got %d", queue.Size)
	}
}
