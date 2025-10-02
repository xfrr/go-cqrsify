package apix

import (
	"encoding/json"
)

type ContentType string

func (ct ContentType) String() string {
	return string(ct)
}

// Media type constants (per spec).
const (
	ContentTypeJSONAPI     ContentType = "application/vnd.api+json"
	ContentTypeProblemJSON ContentType = "application/problem+json"
	ContentTypeMergePatch  ContentType = "application/merge-patch+json"
)

// Links models top-level or relationship links as per JSON:API.
type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

// Meta allows attaching non-standard metadata.
type Meta map[string]any

// Resource is a minimal JSON:API resource wrapper.
// T should be the concrete attributes struct (e.g., UserAttrs).
type Resource[T any] struct {
	Type          string                  `json:"type"`
	ID            string                  `json:"id,omitempty"`
	Attributes    T                       `json:"attributes"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
	Links         *Links                  `json:"links,omitempty"`
	Meta          Meta                    `json:"meta,omitempty"`
}

// Relationship represents a JSON:API resource relationship.
type Relationship struct {
	Links *Links `json:"links,omitempty"`
	Data  any    `json:"data,omitempty"` // single or array of resource identifier(s)
	Meta  Meta   `json:"meta,omitempty"`
}

// SingleDocument is a JSON:API document for one resource.
//
// NOTE: JSON:API says "data" is either an object or array. This provides two distinct
// doc types for clarity and encoding efficiency.
type SingleDocument[T any] struct {
	Data     Resource[T] `json:"data"`
	Links    *Links      `json:"links,omitempty"`
	Included []any       `json:"included,omitempty"` // advanced use (compound docs)
	Meta     Meta        `json:"meta,omitempty"`
}

// ManyDocument for multiple resources of the same type.
type ManyDocument[T any] struct {
	Data     []Resource[T] `json:"data"`
	Links    *Links        `json:"links,omitempty"`
	Included []any         `json:"included,omitempty"`
	Meta     Meta          `json:"meta,omitempty"`
}

// NewSingle creates a minimal single resource document.
func NewSingle[T any](typ, id string, attrs T) SingleDocument[T] {
	return SingleDocument[T]{
		Data: Resource[T]{
			Type:       typ,
			ID:         id,
			Attributes: attrs,
		},
	}
}

// NewMany creates a minimal many-resource document from a slice of ID+Attrs pairs.
func NewMany[T any](typ string, items []struct {
	ID    string
	Attrs T
}) ManyDocument[T] {
	out := ManyDocument[T]{}
	out.Data = make([]Resource[T], 0, len(items))
	for _, it := range items {
		out.Data = append(out.Data, Resource[T]{
			Type:       typ,
			ID:         it.ID,
			Attributes: it.Attrs,
		})
	}
	return out
}

func (d SingleDocument[T]) MarshalJSON() ([]byte, error) {
	type alias SingleDocument[T]
	return json.Marshal(alias(d))
}

func (d ManyDocument[T]) MarshalJSON() ([]byte, error) {
	type alias ManyDocument[T]
	return json.Marshal(alias(d))
}

func UnmarshalSingleDocument[T any](r []byte) (SingleDocument[T], error) {
	var doc SingleDocument[T]
	err := json.Unmarshal(r, &doc)
	return doc, err
}

func UnmarshalManyDocument[T any](data []byte) (ManyDocument[T], error) {
	var doc ManyDocument[T]
	err := json.Unmarshal(data, &doc)
	return doc, err
}
