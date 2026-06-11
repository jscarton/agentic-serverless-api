// internal/mcp/registry.go
package mcp

import (
	"context"
	"errors"
)

var ErrToolNotFound = errors.New("tool not found")

type Tool struct {
	Name        string
	Description string
	InputSchema map[string]any
	Handler     func(ctx context.Context, args map[string]any) (any, error)
}

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{tools: make(map[string]Tool)}
}

func (r *Registry) Register(t Tool) {
	r.tools[t.Name] = t
}

func (r *Registry) List() []Tool {
	list := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		list = append(list, t)
	}
	return list
}

func (r *Registry) Call(ctx context.Context, name string, args map[string]any) (any, error) {
	t, ok := r.tools[name]
	if !ok {
		return nil, ErrToolNotFound
	}
	return t.Handler(ctx, args)
}
