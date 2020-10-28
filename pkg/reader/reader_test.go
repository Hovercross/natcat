package reader_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hovercross/natcat/pkg/reader"
)

func Test_Unwrapped(t *testing.T) {
	testData := []byte("Hello world")
	publisher := &mockPublisher{}

	r := reader.Reader{
		Input:   bytes.NewBuffer(testData),
		Publish: publisher.Publish,
	}

	if err := r.Run(); err != nil {
		t.Error(err)
	}

	expected := &mockPublisher{
		Messages: []string{string(testData)},
	}

	for _, difference := range deep.Equal(publisher, expected) {
		t.Error(difference)
	}
}

func Test_Wrapped(t *testing.T) {
	testData := "Hello world"
	publisher := &mockPublisher{}

	r := reader.Reader{
		Input:   bytes.NewBufferString(testData),
		Publish: publisher.Publish,
		Wrap:    true,
	}

	if err := r.Run(); err != nil {
		t.Error(err)
	}

	expected := []string{
		`{"sequence":0,"type":"startup"}`,
		`{"sequence":1,"type":"line","value":"Hello world"}`,
		`{"sequence":2,"type":"shutdown"}`,
	}

	for _, difference := range deep.Equal(publisher.Messages, expected) {
		t.Error(difference)
	}
}

func Test_WrappedJSON(t *testing.T) {
	testData := `{"hello": "world"}`
	publisher := &mockPublisher{}

	r := reader.Reader{
		Input:     bytes.NewBufferString(testData),
		Publish:   publisher.Publish,
		Wrap:      true,
		JSONInput: true,
	}

	if err := r.Run(); err != nil {
		t.Error(err)
	}

	expected := []string{
		`{"sequence":0,"type":"startup"}`,
		`{"sequence":1,"type":"json","data":{"hello":"world"}}`,
		`{"sequence":2,"type":"shutdown"}`,
	}

	for _, difference := range deep.Equal(publisher.Messages, expected) {
		t.Error(difference)
	}
}

func Test_WrappedJSONError(t *testing.T) {
	testData := `{"hello": "world"`
	publisher := &mockPublisher{}

	r := reader.Reader{
		Input:     bytes.NewBufferString(testData),
		Publish:   publisher.Publish,
		Wrap:      true,
		JSONInput: true,
	}

	if err := r.Run(); err != nil {
		t.Error(err)
	}

	expected := []string{
		`{"sequence":0,"type":"startup"}`,
		`{"sequence":1,"type":"Error: JSON Encode","input":"{\"hello\": \"world\"","error":"json: error calling MarshalJSON for type json.RawMessage: unexpected end of JSON input"}`,
		`{"sequence":2,"type":"shutdown"}`,
	}

	for _, difference := range deep.Equal(publisher.Messages, expected) {
		t.Error(difference)
	}
}

func Test_Timestamp(t *testing.T) {
	testData := "Hello world"
	publisher := &mockPublisher{}

	r := reader.Reader{
		Input:         bytes.NewBufferString(testData),
		Publish:       publisher.Publish,
		Wrap:          true,
		TimeGenerator: jan1,
	}

	if err := r.Run(); err != nil {
		t.Error(err)
	}

	expected := []string{
		`{"timestamp":"2020-01-01T00:00:00Z","sequence":0,"type":"startup"}`,
		`{"timestamp":"2020-01-01T00:00:00Z","sequence":1,"type":"line","value":"Hello world"}`,
		`{"timestamp":"2020-01-01T00:00:00Z","sequence":2,"type":"shutdown"}`,
	}

	for _, difference := range deep.Equal(publisher.Messages, expected) {
		t.Error(difference)
	}
}

func Test_InvalidReader(t *testing.T) {
	ir := &invalidReader{}
	publisher := &mockPublisher{}

	r := reader.Reader{
		Input:   ir,
		Publish: publisher.Publish,
		Wrap:    true,
	}

	if err := r.Run(); err == nil {
		t.Error("err is nil")
	}

	expected := []string{
		`{"sequence":0,"type":"startup"}`,
		`{"sequence":1,"type":"Error: Input Reader","error":"Invalid reader test"}`,
	}

	for _, difference := range deep.Equal(publisher.Messages, expected) {
		t.Error(difference)
	}
}

func Test_PublishError(t *testing.T) {
	testData := "Hello world"
	publisher := &mockPublisher{
		err: fmt.Errorf("test"),
	}

	r := reader.Reader{
		Input:         bytes.NewBufferString(testData),
		Publish:       publisher.Publish,
		Wrap:          true,
		TimeGenerator: jan1,
	}

	if err := r.Run(); err == nil {
		t.Error("err is nil")
	}

}

type mockPublisher struct {
	Messages []string // String makes it easier to compare
	err      error
}

func (mp *mockPublisher) Publish(data []byte) error {
	mp.Messages = append(mp.Messages, string(data))

	return mp.err
}

func jan1() time.Time {
	return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
}

type invalidReader struct {
}

func (ir *invalidReader) Read(buf []byte) (int, error) {
	return 0, fmt.Errorf("Invalid reader test")
}
