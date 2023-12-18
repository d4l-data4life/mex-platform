package forest

import (
	"fmt"
	"strings"
)

type myForest[T any] struct {
	nodes []myNode[T]
	vidx  map[NodeID]int

	rootpathCache map[NodeID][]NodeID

	valid bool
}

type myNode[T any] struct {
	id       NodeID
	parentID NodeID
	payload  T
	depth    int
}

const defaultNodeSliceCap = 32

func NewForestWriter[T any]() Writer[T] {
	return &myForest[T]{
		nodes: make([]myNode[T], 0, defaultNodeSliceCap),
		vidx:  make(map[string]int),

		valid: false, // not yet valid, will be after Seal
	}
}

func (f *myForest[T]) Add(id NodeID, parentID NodeID, payload T) bool {
	if f.valid {
		// cannot add to a sealed forest
		return false
	}

	if id == "" {
		return false
	}

	if _, ok := f.vidx[id]; ok {
		// Node with same ID already exists.
		return false
	}

	f.nodes = append(f.nodes, myNode[T]{
		id:       id,
		parentID: parentID,
		payload:  payload,
		depth:    -1,
	})
	f.vidx[id] = len(f.nodes) - 1

	return true
}

func (f *myForest[T]) Seal() Reader[T] {
	if !f.valid {
		f.valid = f.determineDepths()
		f.rootpathCache = make(map[NodeID][]NodeID)
	}

	if f.valid {
		return f
	}
	return nil
}

// Determine depth values for each node
func (f *myForest[T]) determineDepths() bool {
	// Preparation
	foundRoots := false
	for i := range f.nodes {
		n := &f.nodes[i]
		if n.parentID == "" {
			n.depth = 0
			foundRoots = true
		}
	}

	if !foundRoots {
		return false
	}

	// We iterate over the nodes list and set each non-root node's depth to its parent's depth +1.
	// In case the parent's depth was sure, the child node is also sure, that is. the depth value is correct.
	for {
		allSure := true
		becomeSure := false

		for i := range f.nodes {
			n := &f.nodes[i]
			if n.depth == -1 && n.parentID != "" {
				if pidx, ok := f.vidx[n.parentID]; ok {
					parentNode := f.nodes[pidx]
					if parentNode.depth > -1 {
						n.depth = parentNode.depth + 1
						becomeSure = true
					}
				} else {
					return false
				}
			}
			allSure = allSure && n.depth > -1
		}

		// If the nodes form a valid set of trees, we need to get at least one additional sure node each iteration.
		if allSure {
			return true
		}
		if !becomeSure {
			return false
		}
	}
}

func (f *myForest[T]) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("size: %d (%d)\n", len(f.nodes), len(f.vidx)))

	for k, v := range f.vidx {
		sb.WriteString(fmt.Sprintf("%-10s  ---> %3d\n", k, v))
	}

	for i := range f.nodes {
		sb.WriteString(fmt.Sprintf("%3d: %v\n", i, f.nodes[i]))
	}

	return sb.String()
}

//revive:disable-next-line:receiver-naming
func (n myNode[T]) String() string {
	if n.parentID == "" {
		return fmt.Sprintf("<'%s', parent=∅, depth=%d, data=%v>", n.id, n.depth, n.payload)
	}
	return fmt.Sprintf("<'%s', parent='%s', depth=%d, data=%v>", n.id, n.parentID, n.depth, n.payload)
}
