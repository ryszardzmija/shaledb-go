package membership

import (
	"github.com/ryszardzmija/shaledb-go/internal/identity"
)

func (m File) IdentityNodes() []identity.Node {
	nodes := make([]identity.Node, len(m.Nodes))

	for i, node := range m.Nodes {
		nodes[i] = identity.Node{
			ID: identity.NodeID(node.ID),
			Address: identity.NodeAddress{
				Host: node.Address.Host,
				Port: node.Address.Port,
			},
		}
	}

	return nodes
}
