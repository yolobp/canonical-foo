package protoconv

import (
	"context"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

func TestUnmarshalYAML(t *testing.T) {
	yb := []byte(`foo: bar
abc: xyz
isOK: true`)

	v := &structpb.Struct{}

	if err := UnmarshalYAML(yb, v); err != nil {
		t.Error(err)
	}
}

type fakeProcessor struct{}

func (p *fakeProcessor) Process(_ context.Context, m *structpb.Struct) error {
	m.Fields["new_attr"] = structpb.NewNumberValue(100.0)
	return nil
}

// Example of Handler initialization.
// I could replace *structpb.Struct with any other proto message type.
func TestHandler(t *testing.T) {
	h := NewHandler([]Processor[*structpb.Struct]{&fakeProcessor{}})
	if err := h.Do(context.Background(), ""); err != nil {
		t.Error(err)
	}

	// This fails because string is not a proto message.
	// h2 := NewHandler([]Processor[string]{})
}
