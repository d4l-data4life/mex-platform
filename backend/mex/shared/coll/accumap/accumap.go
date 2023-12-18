package accumap

import (
	"fmt"
	"strings"
)

type Accumap[D any, B any] interface {
	Push(value D)
	PushWithKey(key string, value D)

	GetByKey(key string) (*B, error)
	GetByKeyOrNil(key string) *B

	Keys() []string

	Size() int
	String() string

	ToMap() map[string]B
}

type myGathermap[D any, B any] struct {
	accu  func(bucket B, data D) B
	zero  func(data D) B
	keyer func(data D) string

	m    map[string]*B
	keys []string
}

const defaultKeySliceCap = 32

func NewAccumap[D any, B any](accu func(bucket B, data D) B, zero func(data D) B, keyer func(data D) string) Accumap[D, B] {
	return &myGathermap[D, B]{
		accu:  accu,
		zero:  zero,
		keyer: keyer,
		m:     make(map[string]*B),
		keys:  make([]string, 0, defaultKeySliceCap),
	}
}

func (g *myGathermap[D, B]) Push(value D) {
	g.PushWithKey(g.keyer(value), value)
}

func (g *myGathermap[D, B]) PushWithKey(key string, value D) {
	if bucket, ok := g.m[key]; ok {
		*bucket = g.accu(*bucket, value)
	} else {
		b := g.zero(value)
		g.m[key] = &b
		g.keys = append(g.keys, key)
	}
}

func (g *myGathermap[D, B]) GetByKey(key string) (*B, error) {
	if bucket, ok := g.m[key]; ok {
		return bucket, nil
	}
	return nil, fmt.Errorf("key not found: %s", key)
}

func (g *myGathermap[D, B]) GetByKeyOrNil(key string) *B {
	if bucket, ok := g.m[key]; ok {
		return bucket
	}
	return nil
}

func (g *myGathermap[D, B]) Size() int {
	return len(g.m)
}

func (g *myGathermap[D, B]) Keys() []string {
	return g.keys
}

func (g *myGathermap[D, B]) ToMap() map[string]B {
	m := make(map[string]B)
	for k, v := range g.m {
		m[k] = *v
	}
	return m
}

func (g *myGathermap[D, B]) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("size: %d\n", len(g.m)))

	for k, v := range g.m {
		sb.WriteString(fmt.Sprintf("%-10s  ---> %#v\n", k, v))
	}

	return sb.String()
}
