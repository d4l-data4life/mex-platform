package forest

type NodeID = string

type Reader[T any] interface {
	Size() int

	Depth(id NodeID) (int, error)
	MustDepth(id NodeID) int

	Parent(id NodeID) (NodeID, error)
	MustParent(id NodeID) NodeID

	RootPath(id NodeID) ([]NodeID, error)
	MustRootPath(id NodeID) []NodeID

	GetByIndex(index int) (*T, error)
	GetByID(id NodeID) (*T, error)

	String() string
}

type Writer[T any] interface {
	Add(id NodeID, parentID NodeID, payload T) bool
	Seal() Reader[T]
	Size() int
}
