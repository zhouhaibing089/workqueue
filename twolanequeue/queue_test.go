package twolanequeue

import "testing"

func TestTowLaneQueue(t *testing.T) {
	l := NewStaticLane()

	// prioritize odd numbers
	for i := 0; i < 7; i++ {
		if i%2 == 1 {
			l.Prioritize(i)
		}
	}
	q := NewNamed("", l)

	if q.Len() != 0 {
		t.Errorf("expected 0, got %d", q.Len())
	}

	q.Push(0)
	q.Push(1)
	// [1] -> [div] -> [0]

	val, ok := q.Pop().(int)
	if !ok {
		t.Errorf("expected int, got %T", val)
	}
	// since 1 is odd and has higher priority
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	q.Push(2)
	// [div] -> [0] -> [2]
	val, ok = q.Pop().(int)
	if !ok {
		t.Errorf("expected int, got %T", val)
	}
	// since 1 is odd and has higher priority
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}

	q.Push(3)
	q.Push(4)
	// [3] -> [div] -> [2] -> [4]
	l.Reset(3) // mark 3 as slow
	q.Touch(3)
	// [div] -> [2] -> [4] -> [3]
	val, ok = q.Pop().(int)
	if !ok {
		t.Errorf("expected int, got %T", val)
	}
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}

	q.Touch(4)
	val, ok = q.Pop().(int)
	if !ok {
		t.Errorf("expected int, got %T", val)
	}
	if val != 4 {
		t.Errorf("expected 4, got %d", val)
	}

	q.Push(5)
	q.Push(6)
	// [5] -> [div] -> [3] -> [6]
	l.Prioritize(6) // mark 6 as fast
	q.Touch(6)
	// [5] -> [6] -> [div] -> [3]
	q.Touch(5)
	val, ok = q.Pop().(int)
	if !ok {
		t.Errorf("expected int, got %T", val)
	}
	if val != 5 {
		t.Errorf("expected 5, got %d", val)
	}

	val, ok = q.Pop().(int)
	if !ok {
		t.Errorf("expected int, got %T", val)
	}
	if val != 6 {
		t.Errorf("expected 6, got %d", val)
	}

	if q.Len() != 1 {
		t.Errorf("expected 1, got %d", q.Len())
	}
}
