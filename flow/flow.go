// Package flow provides the core pipeline and flow orchestration primitives for eino.
// It enables composing LLM components into directed acyclic graphs (DAGs) for
// complex AI workflows.
package flow

import (
	"context"
	"fmt"
	"sync"
)

// NodeType represents the type of a node in the flow graph.
type NodeType string

const (
	// NodeTypeLLM represents a large language model node.
	NodeTypeLLM NodeType = "llm"
	// NodeTypeTool represents a tool/function call node.
	NodeTypeTool NodeType = "tool"
	// NodeTypeRetriever represents a retriever/search node.
	NodeTypeRetriever NodeType = "retriever"
	// NodeTypeTransform represents a data transformation node.
	NodeTypeTransform NodeType = "transform"
)

// Node represents a single processing unit within a flow.
type Node struct {
	// ID is the unique identifier for this node.
	ID string
	// Type describes the kind of processing this node performs.
	Type NodeType
	// Handler is the function invoked when this node is executed.
	Handler func(ctx context.Context, input any) (any, error)
}

// Edge represents a directed connection between two nodes.
type Edge struct {
	From string
	To   string
}

// Flow is a directed acyclic graph of processing nodes.
type Flow struct {
	mu    sync.RWMutex
	nodes map[string]*Node
	edges []Edge
}

// New creates and returns an empty Flow.
func New() *Flow {
	return &Flow{
		nodes: make(map[string]*Node),
	}
}

// AddNode registers a node with the flow. Returns an error if a node
// with the same ID has already been added.
func (f *Flow) AddNode(n *Node) error {
	if n == nil {
		return fmt.Errorf("flow: node must not be nil")
	}
	if n.ID == "" {
		return fmt.Errorf("flow: node ID must not be empty")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, exists := f.nodes[n.ID]; exists {
		return fmt.Errorf("flow: node %q already registered", n.ID)
	}
	f.nodes[n.ID] = n
	return nil
}

// AddEdge creates a directed edge from the node identified by fromID to
// the node identified by toID. Both nodes must be registered before
// calling AddEdge.
func (f *Flow) AddEdge(fromID, toID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.nodes[fromID]; !ok {
		return fmt.Errorf("flow: source node %q not found", fromID)
	}
	if _, ok := f.nodes[toID]; !ok {
		return fmt.Errorf("flow: destination node %q not found", toID)
	}
	f.edges = append(f.edges, Edge{From: fromID, To: toID})
	return nil
}

// Run executes the flow starting from the node identified by startID,
// passing initialInput as the first input. Nodes are executed in
// topological order following the registered edges.
func (f *Flow) Run(ctx context.Context, startID string, initialInput any) (any, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	startNode, ok := f.nodes[startID]
	if !ok {
		return nil, fmt.Errorf("flow: start node %q not found", startID)
	}

	output, err := startNode.Handler(ctx, initialInput)
	if err != nil {
		return nil, fmt.Errorf("flow: node %q failed: %w", startID, err)
	}

	current := startID
	for {
		next := f.nextNode(current)
		if next == "" {
			break
		}
		n := f.nodes[next]
		output, err = n.Handler(ctx, output)
		if err != nil {
			return nil, fmt.Errorf("flow: node %q failed: %w", next, err)
		}
		current = next
	}

	return output, nil
}

// nextNode returns the ID of the first node connected from fromID,
// or an empty string if no outgoing edge exists.
func (f *Flow) nextNode(fromID string) string {
	for _, e := range f.edges {
		if e.From == fromID {
			return e.To
		}
	}
	return ""
}
