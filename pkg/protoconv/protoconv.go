package protoconv

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

func UnmarshalYAML(b []byte, v proto.Message) error {
	tmp := map[string]any{}
	if err := yaml.Unmarshal(b, tmp); err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	jb, err := json.Marshal(tmp)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	if err := protojson.Unmarshal(jb, v); err != nil {
		return fmt.Errorf("failed to unmarshal proto: %w", err)
	}
	return nil
}

type ProtoMessageWrapper[T any] interface {
	proto.Message
	*T
}

type Processor[P proto.Message] interface {
	Process(context.Context, P) error
}

type Handler[T any, P ProtoMessageWrapper[T]] struct {
	processors []Processor[P]
}

func NewHandler[T any, P ProtoMessageWrapper[T]](ps []Processor[P]) *Handler[T, P] {
	return &Handler[T, P]{
		processors: ps,
	}
}

func (h *Handler[T, P]) Do(ctx context.Context, path string) error {
	// Read bytes from path
	// fake data
	yb := []byte(`foo: bar
abc: xyz
isOK: true`)

	p := P(new(T))
	if err := UnmarshalYAML(yb, p); err != nil {
		return err
	}

	for _, processor := range h.processors {
		if err := processor.Process(ctx, p); err != nil {
			return fmt.Errorf("process error: %w", err)
		}
	}

	// send things to downstream.

	// Used for testing.
	// txt, err := prototext.Marshal(p)
	// return fmt.Errorf("inject: %v:\n%v", err, string(txt))
	return nil
}
