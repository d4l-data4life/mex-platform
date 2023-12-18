package forest

import "fmt"

func (f *myForest[T]) Depth(id NodeID) (int, error) {
	if !f.valid {
		panic("invalid forest")
	}

	if idx, ok := f.vidx[id]; ok {
		return f.nodes[idx].depth, nil
	}

	return -1, fmt.Errorf("unknown node ID: %s", id)
}

func (f *myForest[T]) Parent(id NodeID) (NodeID, error) {
	if !f.valid {
		panic("invalid forest")
	}

	if idx, ok := f.vidx[id]; ok {
		return f.nodes[idx].parentID, nil
	}
	return "", fmt.Errorf("unknown node ID: %s", id)
}

func (f *myForest[T]) Size() int { return len(f.nodes) }

func (f *myForest[T]) RootPath(id NodeID) ([]NodeID, error) {
	if !f.valid {
		panic("invalid forest")
	}

	if path, ok := f.rootpathCache[id]; ok {
		return path, nil
	}

	idx, ok := f.vidx[id]
	if !ok {
		return nil, fmt.Errorf("unknown node ID: %s", id)
	}

	path := []NodeID{}
	n := f.nodes[idx]
	for {
		if n.parentID == "" {
			f.rootpathCache[id] = path
			return path, nil
		}
		path = append(path, n.parentID)
		n = f.nodes[f.vidx[n.parentID]]
	}
}

func (f *myForest[T]) GetByIndex(index int) (*T, error) {
	if index < 0 {
		return nil, fmt.Errorf("index out of bounds: %d", index)
	}

	if index > len(f.nodes)-1 {
		return nil, fmt.Errorf("index out of bounds: %d", index)
	}

	return &f.nodes[index].payload, nil
}

func (f *myForest[T]) GetByID(id NodeID) (*T, error) {
	idx, ok := f.vidx[id]
	if !ok {
		return nil, fmt.Errorf("unknown ID: %s", id)
	}

	return &f.nodes[idx].payload, nil
}

// MustX versions of methods

func (f *myForest[T]) MustParent(id NodeID) NodeID {
	x, err := f.Parent(id)
	if err != nil {
		panic(err)
	}
	return x
}

func (f *myForest[T]) MustDepth(id NodeID) int {
	d, err := f.Depth(id)
	if err != nil {
		panic(err)
	}
	return d
}

func (f *myForest[T]) MustRootPath(id NodeID) []NodeID {
	x, err := f.RootPath(id)
	if err != nil {
		panic(err)
	}
	return x
}
