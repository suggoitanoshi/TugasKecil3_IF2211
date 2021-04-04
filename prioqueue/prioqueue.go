package prioqueue

type PrioQueueEl struct {
	info interface{}
	prio float64
}

type PrioQueue struct {
	els []PrioQueueEl
}

func NewPQ() PrioQueue {
	pq := PrioQueue{}
	pq.els = make([]PrioQueueEl, 0)
	return pq
}

func (pq *PrioQueue) Enqueue(el interface{}, prio float64) {
	newEl := PrioQueueEl{el, prio}
	if len(pq.els) == 0 {
		pq.els = append(pq.els, newEl)
	} else {
		queued := false
		for i, iter := range pq.els {
			if iter.prio >= prio {
				pq.els = append(pq.els[:i], append([]PrioQueueEl{newEl}, pq.els[i:]...)...)
				queued = true
				break
			}
		}
		if !queued {
			pq.els = append(pq.els, newEl)
		}
	}
}

func (pq PrioQueue) GetElements() []interface{} {
	infos := make([]interface{}, 0)
	for _, v := range pq.els {
		infos = append(infos, v.info)
	}
	return infos
}

func (pq *PrioQueue) DequeueElementAt(idx int) interface{} {
	old := pq.els[idx]
	pq.els = append(pq.els[:idx], pq.els[idx+1:]...)
	return old.info
}

func (pq *PrioQueue) Dequeue() interface{} {
	old := pq.els[0]
	pq.els = pq.els[1:]
	return old.info
}

func (pq PrioQueue) Peek() interface{} {
	return pq.els[0].info
}

func (pq PrioQueue) IsEmpty() bool {
	return len(pq.els) == 0
}
