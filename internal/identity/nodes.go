package identity

import "fmt"

type NodeID string

type NodeAddress struct {
	Host string
	Port uint16
}

type Node struct {
	ID      NodeID
	Address NodeAddress
}

type ClusterView struct {
	self  Node
	nodes []Node
}

func (v ClusterView) Self() Node {
	return v.self
}

func (v ClusterView) Nodes() []Node {
	nodes := make([]Node, len(v.nodes))
	copy(nodes, v.nodes)
	return nodes
}

func NewClusterView(selfID NodeID, nodes []Node) (ClusterView, error) {
	if err := validateClusterViewInput(selfID, nodes); err != nil {
		return ClusterView{}, err
	}

	nodesCopy := make([]Node, len(nodes))
	copy(nodesCopy, nodes)

	selfNode, err := getSelfNode(selfID, nodesCopy)
	if err != nil {
		return ClusterView{}, err
	}

	return ClusterView{
		self:  selfNode,
		nodes: nodesCopy,
	}, nil
}

func getSelfNode(selfID NodeID, nodes []Node) (Node, error) {
	for _, node := range nodes {
		if selfID == node.ID {
			return node, nil
		}
	}

	return Node{}, fmt.Errorf("identity invariant violated: selfID %q not found after validation", selfID)
}
