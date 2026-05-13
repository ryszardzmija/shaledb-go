package identity

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field  string
	Value  string
	Reason string
	nodeID NodeID
}

func (e ValidationError) Error() string {
	msg := e.Field + ": " + e.Reason

	if e.Value != "" {
		msg += fmt.Sprintf(" (value %s)", strconv.Quote(e.Value))
	}

	if e.nodeID != "" {
		msg += fmt.Sprintf(" for node %s", strconv.Quote(string(e.nodeID)))
	}

	return msg
}

func validateClusterViewInput(selfID NodeID, nodes []Node) error {
	if err := validateNodeID("selfID", selfID); err != nil {
		return err
	}

	if err := validateNodeListNotEmpty(nodes); err != nil {
		return err
	}

	for i, node := range nodes {
		if err := validateNodeID(fmt.Sprintf("nodes[%d].id", i), node.ID); err != nil {
			return err
		}
	}

	if err := validateUniquenessOfNodeIDs(nodes); err != nil {
		return err
	}

	if err := validateHosts(nodes); err != nil {
		return err
	}

	if err := validatePorts(nodes); err != nil {
		return err
	}

	if err := validateNodeAddresses(nodes); err != nil {
		return err
	}

	if err := validateSelfExistsInNodes(selfID, nodes); err != nil {
		return err
	}

	return nil
}

func validateUniquenessOfNodeIDs(nodes []Node) error {
	seenIDs := make(map[NodeID]int)

	for i, node := range nodes {
		if _, ok := seenIDs[node.ID]; ok {
			return ValidationError{
				Field:  fmt.Sprintf("nodes[%d].id", i),
				Value:  string(node.ID),
				Reason: fmt.Sprintf("duplicates nodes[%d].id", seenIDs[node.ID]),
			}
		}
		seenIDs[node.ID] = i
	}

	return nil
}

func validateHosts(nodes []Node) error {
	for i, node := range nodes {
		trimmedHost := strings.TrimSpace(node.Address.Host)

		if trimmedHost != node.Address.Host {
			return ValidationError{
				Field:  fmt.Sprintf("nodes[%d].address.host", i),
				Value:  node.Address.Host,
				Reason: "must not contain leading or trailing whitespace",
				nodeID: node.ID,
			}
		}

		if trimmedHost == "" {
			return ValidationError{
				Field:  fmt.Sprintf("nodes[%d].address.host", i),
				Value:  node.Address.Host,
				Reason: "must not be empty",
				nodeID: node.ID,
			}
		}
	}

	return nil
}

func validatePorts(nodes []Node) error {
	for i, node := range nodes {
		if node.Address.Port == 0 {
			return ValidationError{
				Field:  fmt.Sprintf("nodes[%d].address.port", i),
				Value:  strconv.Itoa(int(node.Address.Port)),
				Reason: "must be non-zero",
				nodeID: node.ID,
			}
		}
	}

	return nil
}

func validateNodeAddresses(nodes []Node) error {
	seenAddresses := make(map[NodeAddress]int)

	for i, node := range nodes {
		if _, ok := seenAddresses[node.Address]; ok {
			return ValidationError{
				Field:  fmt.Sprintf("nodes[%d].address", i),
				Value:  net.JoinHostPort(node.Address.Host, strconv.Itoa(int(node.Address.Port))),
				Reason: fmt.Sprintf("duplicates nodes[%d].address", seenAddresses[node.Address]),
				nodeID: node.ID,
			}
		}
		seenAddresses[node.Address] = i
	}

	return nil
}

func validateSelfExistsInNodes(selfID NodeID, nodes []Node) error {
	for _, node := range nodes {
		if selfID == node.ID {
			return nil
		}
	}

	return ValidationError{
		Field:  "selfID",
		Value:  string(selfID),
		Reason: "not found among cluster members",
	}
}

func validateNodeID(field string, id NodeID) error {
	trimmedID := strings.TrimSpace(string(id))

	if trimmedID != string(id) {
		return ValidationError{
			Field:  field,
			Value:  string(id),
			Reason: "must not contain leading or trailing whitespace",
		}
	}

	if trimmedID == "" {
		return ValidationError{
			Field:  field,
			Value:  string(id),
			Reason: "must not be empty",
		}
	}

	return nil
}

func validateNodeListNotEmpty(nodes []Node) error {
	if len(nodes) == 0 {
		return ValidationError{
			Field:  "nodes",
			Value:  "",
			Reason: "must not be empty",
		}
	}

	return nil
}
