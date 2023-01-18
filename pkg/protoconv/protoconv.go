package protoconv

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

// General func to umarshal yaml bytes to proto.
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

// Wrap the proto message interface. I didn't expect it's needed. But this
// constraint is essential in order to initialize a new _generic_ proto message.
type ProtoMessageWrapper[T any] interface {
	proto.Message
	*T
}

// A generic interface for processing proto messages.
type Processor[P proto.Message] interface {
	Process(context.Context, P) error
}

// Assume the handler has common logic to retrieve the proto message from place.
// The proto message could be any proto message type. But an instance of Handler
// can only handle one type of proto message.
//
// Handler will then call a list of processors for that type of proto message
// to process the messages.
//
// After the processing, there will also be some common logic to marshal the
// proto message and send it somewhere else.
type Handler[T any, P ProtoMessageWrapper[T]] struct {
	processors []Processor[P]
}

// Create a new Handler with the given processors. The proto message type needs
// to match.
func NewHandler[T any, P ProtoMessageWrapper[T]](ps []Processor[P]) *Handler[T, P] {
	return &Handler[T, P]{
		processors: ps,
	}
}

func (h *Handler[T, P]) Do(ctx context.Context, path string) error {
	// Common logic to read yaml bytes from somewhere.

	// Here we use fake data.
	yb := []byte(`foo: bar
abc: xyz
isOK: true`)

	// This is why ProtoMessageWrapper is needed.
	// I cannot create a new proto message of the type without it.
	// Referenced: https://stackoverflow.com/questions/72090387/what-is-the-generic-type-for-a-pointer-that-implements-an-interface
	p := P(new(T))
	if err := UnmarshalYAML(yb, p); err != nil {
		return err
	}

	for _, processor := range h.processors {
		if err := processor.Process(ctx, p); err != nil {
			return fmt.Errorf("process error: %w", err)
		}
	}

	// jb, err := protojson.Marshal(p)
	// if err != nil {
	// 	return err
	// }
	// send things to downstream.
	return nil
}
