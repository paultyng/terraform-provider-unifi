package progress

import (
	"time"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/identity"
	"github.com/opencontainers/go-digest"
)

type Writer interface {
	Write(*client.SolveStatus)
	WriteBuildRef(string, string)
	ValidateLogSource(digest.Digest, interface{}) bool
	ClearLogSource(interface{})
}

func Write(w Writer, name string, f func() error) {
	dgst := digest.FromBytes([]byte(identity.NewID()))
	tm := time.Now()

	vtx := client.Vertex{
		Digest:  dgst,
		Name:    name,
		Started: &tm,
	}

	w.Write(&client.SolveStatus{
		Vertexes: []*client.Vertex{&vtx},
	})

	err := f()

	tm2 := time.Now()
	vtx2 := vtx
	vtx2.Completed = &tm2
	if err != nil {
		vtx2.Error = err.Error()
	}
	w.Write(&client.SolveStatus{
		Vertexes: []*client.Vertex{&vtx2},
	})
}

func WriteBuildRef(w Writer, target string, ref string) {
	w.WriteBuildRef(target, ref)
}

func NewChannel(w Writer) (chan *client.SolveStatus, chan struct{}) {
	ch := make(chan *client.SolveStatus)
	done := make(chan struct{})
	go func() {
		for {
			v, ok := <-ch
			if !ok {
				close(done)
				w.ClearLogSource(done)
				return
			}

			if len(v.Logs) > 0 {
				logs := make([]*client.VertexLog, 0, len(v.Logs))
				for _, l := range v.Logs {
					if w.ValidateLogSource(l.Vertex, done) {
						logs = append(logs, l)
					}
				}
				v.Logs = logs
			}

			w.Write(v)
		}
	}()
	return ch, done
}

type tee struct {
	Writer
	ch chan *client.SolveStatus
}

func (t *tee) Write(v *client.SolveStatus) {
	v2 := *v
	t.ch <- &v2
	t.Writer.Write(v)
}

func Tee(w Writer, ch chan *client.SolveStatus) Writer {
	if ch == nil {
		return w
	}
	return &tee{
		Writer: w,
		ch:     ch,
	}
}
