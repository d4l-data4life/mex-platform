package forest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyForest(t *testing.T) {
	require.Equal(t, 0, NewForestWriter[any]().Size())
}

func TestAddNode(t *testing.T) {
	f := NewForestWriter[int]()
	require.Equal(t, 0, f.Size())

	require.Equal(t, true, f.Add("foo", "", 1234))
	require.Equal(t, 1, f.Size())

	require.Equal(t, false, f.Add("foo", "", 666))
	require.Equal(t, 1, f.Size())
}

func TestDepth(t *testing.T) {
	f := NewForestWriter[int]()

	require.Equal(t, true, f.Add("foo", "bar", 200))
	require.Equal(t, true, f.Add("bar", "", 100))
	require.Equal(t, true, f.Add("hello", "foo", 300))

	require.Equal(t, 3, f.Size())

	r := f.Seal()
	fmt.Println(r)

	require.Equal(t, tuple(0, nil), tuple(r.Depth("bar")))
	require.Equal(t, tuple(1, nil), tuple(r.Depth("foo")))
	require.Equal(t, tuple(2, nil), tuple(r.Depth("hello")))

	require.Equal(t, tuple("", nil), tuple(r.Parent("bar")))
	require.Equal(t, tuple("bar", nil), tuple(r.Parent("foo")))
	require.Equal(t, tuple("foo", nil), tuple(r.Parent("hello")))

	s, err := r.Parent("bla")
	require.Equal(t, "", s)
	require.Error(t, err)
}

func TestInvalidForest(t *testing.T) {
	f := NewForestWriter[int]()

	require.Equal(t, true, f.Add("foo", "bar", 1234))
	require.Nil(t, f.Seal())

	require.Equal(t, true, f.Add("bar", "", 1234))
	require.NotNil(t, f.Seal())
}

func TestCannotAddToSealedForest(t *testing.T) {
	f := NewForestWriter[int]()

	f.Add("foo", "", 1234)

	r := f.Seal()
	require.NotNil(t, r)

	require.Equal(t, false, f.Add("bar", "", 666))
}

func TestRootPath(t *testing.T) {
	f := NewForestWriter[int]()

	f.Add("D4L", "", 0)
	f.Add("MFI", "D4L", 0)
	f.Add("MF", "MFI", 0)
	f.Add("MF-1", "MF", 0)
	f.Add("MF-2", "MF", 0)
	f.Add("MF-3", "MF", 0)
	f.Add("MF-4", "MF", 0)

	r := f.Seal()
	require.NotNil(t, r)

	fmt.Println(r)

	require.Equal(t, tuple([]NodeID{"MF", "MFI", "D4L"}, nil), tuple(r.RootPath("MF-1")))
	require.Equal(t, tuple([]NodeID{"MF", "MFI", "D4L"}, nil), tuple(r.RootPath("MF-2")))
	require.Equal(t, tuple([]NodeID{"MF", "MFI", "D4L"}, nil), tuple(r.RootPath("MF-2")))
	require.Equal(t, tuple([]NodeID{"MF", "MFI", "D4L"}, nil), tuple(r.RootPath("MF-3")))

	require.Equal(t, tuple([]NodeID{"MFI", "D4L"}, nil), tuple(r.RootPath("MF")))
	require.Equal(t, tuple([]NodeID{"D4L"}, nil), tuple(r.RootPath("MFI")))
	require.Equal(t, tuple([]NodeID{}, nil), tuple(r.RootPath("D4L")))

	s, err := r.RootPath("foo")
	require.Nil(t, s)
	require.Error(t, err)
}

type info struct {
	display string
}

func TestIteration(t *testing.T) {
	f := NewForestWriter[info]()

	f.Add("D4L", "", info{"Headquarter"})
	f.Add("MFI", "D4L", info{"MFI..."})
	f.Add("MF", "MFI", info{"MF"})
	f.Add("MF-1", "MF", info{"MF-1"})

	r := f.Seal()
	require.NotNil(t, r)

	for i := 0; i < r.Size(); i++ {
		data, err := r.GetByIndex(i)
		require.Nil(t, err)
		require.NotNil(t, data)
	}
}

func tuple(x ...any) []any {
	return x
}
