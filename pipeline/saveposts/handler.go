package saveposts

import (
	"encoding/json"

	"github.com/beeceej/inflight"
	"github.com/beeceej/posts/pipeline/shared/post"
)

type Handler struct {
	Inflight   *inflight.Inflight
	PostWriter post.PostWriter
}

// Handle is
func (h Handler) Handle(ref inflight.Ref) error {
	b, err := h.Inflight.Get(ref.Object)
	if err != nil {
		return err
	}

	postIndex := new(post.PostIndex)
	if err = json.Unmarshal(b, &postIndex); err != nil {
		return err
	}

	return h.PostWriter.Write(postIndex.Posts)
}
