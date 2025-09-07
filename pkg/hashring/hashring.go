package hashring

import (
	"sort"
	"strconv"

	xx "github.com/cespare/xxhash/v2"
)

type Node struct{ ID, Addr string }

type vnode struct {
	hash uint64
	node Node
}

type Ring struct {
	vnodes   []vnode
	replicas int
}

func New(nodes []Node, replicas int) *Ring {
	r := &Ring{replicas: replicas}
	for _, n := range nodes {
		for i := 0; i < replicas; i++ {
			h := xx.Sum64String(n.ID + ":" + strconv.Itoa(i))
			r.vnodes = append(r.vnodes, vnode{hash: h, node: n})
		}
	}
	sort.Slice(r.vnodes, func(i, j int) bool { return r.vnodes[i].hash < r.vnodes[j].hash })
	return r
}

// N 个不同节点（去重）
func (r *Ring) PickN(key string, n int) []Node {
	if len(r.vnodes) == 0 || n <= 0 {
		return nil
	}
	h := xx.Sum64String(key)
	i := sort.Search(len(r.vnodes), func(i int) bool { return r.vnodes[i].hash >= h })
	res := make([]Node, 0, n)
	seen := map[string]struct{}{}
	for j := 0; len(res) < n && j < len(r.vnodes); j++ {
		vn := r.vnodes[(i+j)%len(r.vnodes)]
		if _, ok := seen[vn.node.ID]; ok {
			continue
		}
		seen[vn.node.ID] = struct{}{}
		res = append(res, vn.node)
	}
	return res
}
