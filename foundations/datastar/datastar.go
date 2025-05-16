package datastar

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Responder struct {
	rc *http.ResponseController
	w  http.ResponseWriter
}

func NewResponder(w http.ResponseWriter) Responder {
	rc := http.NewResponseController(w)

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	// flush headers
	if err := rc.Flush(); err != nil {
		// Below panic is a deliberate choice as it should never occur and is an environment issue.
		// https://crawshaw.io/blog/go-and-sqlite
		// In Go, errors that are part of the standard operation of a program are returned as values.
		// Programs are expected to handle errors.
		panic(fmt.Sprintf("response writer failed to flush: %v", err))
	}

	return Responder{
		rc: rc,
		w:  w,
	}
}

func (r Responder) MergeFragment(frag []byte) {
	defer r.finalize()

	r.writeEvent(mergeFragments)
	for ln := range bytes.Lines(frag) {
		r.writeData("fragments", bytes.TrimSuffix(ln, []byte{'\n'}))
	}
}

type templComponent interface {
	Render(ctx context.Context, w io.Writer) error
}

func (r Responder) MergeFragmentTempl(ctx context.Context, frag templComponent) error {
	buf := bytes.NewBuffer(nil)
	err := frag.Render(ctx, buf)
	if err != nil {
		return err
	}

	r.MergeFragment(buf.Bytes())
	return nil
}

func (r Responder) ExecuteScript(script string) {
	defer r.finalize()

	r.writeEvent(executeScript)
	for ln := range strings.Lines(script) {
		r.writeData("script", []byte(strings.TrimSuffix(ln, "\n")))
	}
}

type event struct {
	s string
}

func (e event) String() string {
	return e.s
}

var (
	mergeFragments = event{"datastar-merge-fragments"}
	executeScript  = event{"datastar-execute-script"}
)

func (r Responder) writeEvent(event event) {
	fmt.Fprintf(r.w, "event: %s\n", event)
}

func (r Responder) writeData(k string, v []byte) {
	fmt.Fprintf(r.w, "data: %s %s\n", k, v)
}

func (r Responder) finalize() {
	r.w.Write([]byte{'\n'})
}
