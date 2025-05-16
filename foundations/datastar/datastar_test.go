package datastar

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testServer struct {
	s *httptest.Server
}

func newTestServer(t *testing.T, handler func(Responder)) testServer {
	t.Helper()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responder := NewResponder(w)
		handler(responder)
	}))
	t.Cleanup(s.Close)

	return testServer{s: s}
}

func (ts testServer) emptyRequest(t *testing.T) string {
	t.Helper()

	resp, err := http.Get(ts.s.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return string(bs)
}

func TestDatastar_MergeFragments(t *testing.T) {
	s := newTestServer(t, func(r Responder) {
		r.MergeFragment([]byte("<div>\nHey!\n</div>"))
	})

	resp := s.emptyRequest(t)
	assert.Equal(t, `event: datastar-merge-fragments
data: fragments <div>
data: fragments Hey!
data: fragments </div>

`, resp)
}
