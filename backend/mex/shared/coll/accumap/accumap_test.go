package accumap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var add = func(b int, d int) int { return b + d }
var id = func(d int) int { return d }
var dummy = func(d int) string { return "" }

func TestEmptyGathermap(t *testing.T) {
	require.Equal(t, 0, NewAccumap[int, int](add, id, dummy).Size())
}

func TestSimplePayloads(t *testing.T) {
	m := NewAccumap[int, int](add, id, dummy)

	m.PushWithKey("foo", 1)
	require.Equal(t, 1, *m.GetByKeyOrNil("foo"))

	m.PushWithKey("foo", 2)
	require.Equal(t, 3, *m.GetByKeyOrNil("foo"))

	m.PushWithKey("bar", 1000)
	m.PushWithKey("bar", 1000)
	require.Equal(t, 2000, *m.GetByKeyOrNil("bar"))

	require.Equal(t, 2, m.Size())

	require.ElementsMatch(t, []string{"foo", "bar"}, m.Keys())
}

type service struct {
	name    string
	replica string
	memory  uint32
}

func TestGroupBy(t *testing.T) {
	accu := func(bucket []service, data service) []service {
		return append(bucket, data)
	}

	zero := func(data service) []service {
		return []service{data}
	}

	keyer := func(data service) string {
		return data.name
	}

	m := NewAccumap(accu, zero, keyer)

	m.Push(service{name: "metadata", replica: "met-1", memory: 1000})
	m.Push(service{name: "metadata", replica: "met-2", memory: 2000})
	m.Push(service{name: "metadata", replica: "met-3", memory: 3000})
	m.Push(service{name: "index", replica: "idx-1", memory: 10000})
	m.Push(service{name: "query", replica: "qry-1", memory: 100})
	m.Push(service{name: "query", replica: "qry-2", memory: 200})

	t.Logf(m.String())

	_, err := m.GetByKey("metadata")
	require.Nil(t, err)
}

type row struct {
	key     string
	lang    string
	display string
}

type node struct {
	key      string
	displays map[string]string
}

func accuNode(b node, d row) node {
	b.displays[d.lang] = d.display
	return b
}

func zeroNode(d row) node {
	b := node{
		key:      d.key,
		displays: make(map[string]string),
	}
	b.displays[d.lang] = d.display
	return b
}

func keyerNode(d row) string {
	return d.key
}

func TestComplexePayloads(t *testing.T) {
	m := NewAccumap(accuNode, zeroNode, keyerNode)

	m.Push(row{key: "D4L", lang: "de", display: "Hauptquartier"})
	m.Push(row{key: "MFI", lang: "de", display: "M F I - Deutsch"})
	m.Push(row{key: "MFI", lang: "en", display: "M F I - English"})
	m.Push(row{key: "D4L", lang: "en", display: "Headquarter"})
}
