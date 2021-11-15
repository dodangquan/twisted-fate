package lucky

import (
	"crypto/rand"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/rs/zerolog/log"
)

func NewNumber(value int64) *Number {
	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal().Err(err).Int64("value", value).Msg("Cannot create new number")
	}

	priority, err := hashstructure.Hash(b, hashstructure.FormatV2, nil)
	if err != nil {
		log.Fatal().Err(err).Int64("value", value).Msg("Cannot create new number")
	}

	return &Number{
		Value:    value,
		priority: priority,
	}
}

type Number struct {
	Value    int64
	priority uint64
}

type NumberHeap []*Number

func (h NumberHeap) Len() int {
	return len(h)
}

func (h NumberHeap) Less(i, j int) bool {
	return h[i].priority > h[j].priority
}

func (h NumberHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *NumberHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length, not just its contents.
	*h = append(*h, x.(*Number))
}

func (h *NumberHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	*h = old[0 : n-1]
	return item
}
