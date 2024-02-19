package twolanequeue

import (
	"container/list"
	"sync"

	"github.com/zhouhaibing089/workqueue"
)

type StaticLane interface {
	Lane

	Prioritize(item interface{})
}

type staticLane struct {
	lock      sync.Mutex
	fasttrack map[interface{}]struct{}
}

func (l *staticLane) Slow(item interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	_, ok := l.fasttrack[item]
	return !ok
}

func (l *staticLane) Reset(item interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.fasttrack, item)
}

func (l *staticLane) Prioritize(item interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.fasttrack[item] = struct{}{}
}

func NewStaticLane() StaticLane {
	return &staticLane{
		fasttrack: make(map[interface{}]struct{}),
	}
}

type Lane interface {
	Slow(item interface{}) bool
	Reset(item interface{})
}

func NewNamed(name string, lane Lane) workqueue.Queue {
	return &queue{
		name: name,

		slow: list.New(),
		fast: list.New(),
		lane: lane,

		fastset: make(map[interface{}]*list.Element),
		slowset: make(map[interface{}]*list.Element),
	}
}

type queue struct {
	name string

	slow *list.List
	fast *list.List
	lane Lane

	fastset map[interface{}]*list.Element
	slowset map[interface{}]*list.Element
}

func (q *queue) Touch(item interface{}) {
	slow := q.lane.Slow(item)
	if slow {
		ele, ok := q.fastset[item]
		if !ok {
			return
		}
		// move to slow lane
		q.fast.Remove(ele)
		q.slowset[item] = q.slow.PushBack(item)
		delete(q.fastset, item)
	} else {
		ele, ok := q.slowset[item]
		if !ok {
			return
		}
		// move to fast lane
		q.slow.Remove(ele)
		q.fastset[item] = q.fast.PushBack(item)
		delete(q.slowset, item)
	}
	twoLaneQueueDepth.WithLabelValues(q.name, "slow").Set(float64(len(q.slowset)))
	twoLaneQueueDepth.WithLabelValues(q.name, "fast").Set(float64(len(q.fastset)))
}

func (q *queue) Push(item interface{}) {
	slow := q.lane.Slow(item)
	if slow {
		q.slowset[item] = q.slow.PushBack(item)
	} else {
		q.fastset[item] = q.fast.PushBack(item)
	}
	twoLaneQueueDepth.WithLabelValues(q.name, "slow").Set(float64(len(q.slowset)))
	twoLaneQueueDepth.WithLabelValues(q.name, "fast").Set(float64(len(q.fastset)))
}

func (q *queue) Len() int {
	return q.fast.Len() + q.slow.Len()
}

func (q *queue) Pop() interface{} {
	var item interface{}
	if q.fast.Len() > 0 {
		item = q.fast.Remove(q.fast.Front())
		delete(q.fastset, item)
	} else {
		item = q.slow.Remove(q.slow.Front())
		delete(q.slowset, item)
	}
	twoLaneQueueDepth.WithLabelValues(q.name, "slow").Set(float64(len(q.slowset)))
	twoLaneQueueDepth.WithLabelValues(q.name, "fast").Set(float64(len(q.fastset)))
	q.lane.Reset(item)
	return item
}
