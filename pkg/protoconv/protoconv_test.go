package protoconv

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
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

func BenchmarkHandler(b *testing.B) {
	ps := []Processor[*structpb.Struct]{}
	for k := 0; k < 10; k++ {
		ps = append(ps, &fakeProcessor{})
	}
	h := NewHandler(ps)
	for n := 0; n < b.N; n++ {
		h.Do(context.Background(), "")
	}
}

type fakeProcessor2 struct{}

func (p *fakeProcessor2) Process(_ context.Context, v proto.Message) error {
	m, ok := v.(*structpb.Struct)
	if !ok {
		return fmt.Errorf("not struct can't handle")
	}
	m.Fields["new_attr"] = structpb.NewNumberValue(100.0)
	return nil
}

func TestHandler2(t *testing.T) {
	h := NewHandler2(&structpb.Struct{}, []Processor2{&fakeProcessor2{}})
	if err := h.Do(context.Background(), ""); err != nil {
		t.Error(err)
	}
}

func BenchmarkHandler2(b *testing.B) {
	ps := []Processor2{}
	for k := 0; k < 10; k++ {
		ps = append(ps, &fakeProcessor2{})
	}
	h := NewHandler2(&structpb.Struct{}, ps)
	for n := 0; n < b.N; n++ {
		h.Do(context.Background(), "")
	}
}

// Benchmark result:
// go test ./pkg/protoconv/... -bench=.
// goos: linux
// goarch: amd64
// pkg: github.com/yolobp/canonical-foo/pkg/protoconv
// cpu: Intel(R) Xeon(R) CPU @ 2.00GHz
// BenchmarkHandler-96                55726             21258 ns/op
// BenchmarkHandler2-96               55501             21556 ns/op
// PASS
// ok      github.com/yolobp/canonical-foo/pkg/protoconv   2.852s
