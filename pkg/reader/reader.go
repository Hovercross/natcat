package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/hovercross/natcat/pkg/util"
)

// A Reader reads from input, and calls the publish function
type Reader struct {
	Input      io.Reader
	Publish    func(data []byte) error
	Wrap       bool
	JSONInput  bool
	ReaderID   string
	ReaderName string

	TimeGenerator     func() time.Time // For testing with deterministic output
	MessageIDFunction func() string    // For testing with determinstic output
	messageCount      int

	m sync.Mutex
}

// Run will execute the reader, with a lock, until the input is closed
func (r *Reader) Run() error {
	r.m.Lock()
	defer r.m.Unlock()

	r.messageCount = 0

	// If we're wrapping, we're running in intelligent mode and will emit startup/shutdown
	if r.Wrap {
		bm := r.getBaseMessage()
		bm.MessageType = "startup"
		r.Publish(util.MustMarshal(bm))
	}

	scan := bufio.NewScanner(r.Input)

	for scan.Scan() {
		raw := scan.Bytes()

		if err := r.Publish(r.translate(raw)); err != nil {
			return fmt.Errorf("Unable to publish message: %v", err)
		}
	}

	if err := scan.Err(); err != nil {
		if r.Wrap {
			// If we only had an input error, create a synthetic message
			msg := r.getError(r.getBaseMessage(), nil, err, "Error: Input Reader")

			r.Publish(msg)
		}

		return fmt.Errorf("Unable to read input stream: %v", err)
	}

	if r.Wrap {
		bm := r.getBaseMessage()
		bm.MessageType = "shutdown"
		r.Publish(util.MustMarshal(bm))
	}

	return nil
}

func (r *Reader) translate(data []byte) []byte {
	if !r.Wrap {
		return data
	}

	bm := r.getBaseMessage()

	if r.JSONInput {
		bm.MessageType = "json"

		msg := wrappedJSONMessage{
			baseMessage: bm,
			Data:        json.RawMessage(data),
		}

		out, err := json.Marshal(msg)

		if err != nil {
			log.Printf("Error marshaling: %v", err)
			return r.getError(bm, data, err, "Error: JSON Encode")
		}

		return out
	}

	msg := wrappedMessage{
		baseMessage: bm,
		Message:     string(data),
	}

	msg.MessageType = "line"

	return util.MustMarshal(msg)
}

func (r *Reader) getError(bm baseMessage, data []byte, err error, errorType string) []byte {
	// We re-use the base messages to we don't increment the counter

	em := errorMessage{
		baseMessage: bm,
		Input:       string(data),
		Error:       err.Error(),
	}

	em.MessageType = errorType

	// Nothing in this can fail
	return util.MustMarshal(em)
}

type baseMessage struct {
	InstanceID  string     `json:"natcatInstanceID,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty"`
	Sequence    int        `json:"sequence"`
	MessageType string     `json:"type,omitempty"`
	MessageID   string     `json:"id,omitempty"`
	ReaderName  string     `json:"name,omitempty"`
}

type wrappedMessage struct {
	baseMessage
	Message string `json:"value"`
}

type wrappedJSONMessage struct {
	baseMessage
	Data json.RawMessage `json:"data"`
}

type errorMessage struct {
	baseMessage
	Input string `json:"input,omitempty"`
	Error string `json:"error"`
}

func (r *Reader) getBaseMessage() baseMessage {
	defer func() {
		r.messageCount++
	}()

	bm := baseMessage{
		InstanceID: r.ReaderID,
		Sequence:   r.messageCount,
	}

	if r.TimeGenerator != nil {
		t := r.TimeGenerator()
		bm.Timestamp = &t
	}

	if r.MessageIDFunction != nil {
		bm.MessageID = r.MessageIDFunction()
	}

	bm.ReaderName = r.ReaderName

	return bm
}
