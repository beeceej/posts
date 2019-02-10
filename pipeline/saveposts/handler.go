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
func (h Handler) Handle(ref inflight.Ref) (*inflight.Ref, error) {
	b, err := h.Inflight.Get(ref.Object)
	if err != nil {
		return nil, err
	}

	postIndex := &post.PostIndex{
		Posts: []post.Post{},
	}

	if err = json.Unmarshal(b, &postIndex.Posts); err != nil {
		return nil, err
	}

	if err = h.PostWriter.Write(postIndex.Posts); err != nil {
		return nil, err
	}

	return h.Inflight.Write(b)
}
