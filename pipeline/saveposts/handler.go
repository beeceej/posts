package saveposts

import (
	"encoding/json"

	"github.com/beeceej/posts/pipeline/shared/domain"
	"github.com/beeceej/posts/pipeline/shared/inflight"
	"github.com/beeceej/posts/pipeline/shared/state"
)

type Handler struct {
	Inflight *inflight.Inflight
	Saver    *PostSaver
}

// Handle is
func (h Handler) Handle(ref state.InflightRef) error {
	b, err := h.Inflight.Get(ref.Object)
	if err != nil {
		return err
	}

	posts := new(domain.PostIndex)
	if err = json.Unmarshal(b, &posts.Posts); err != nil {
		return err
	}
	return h.Saver.SavePosts(posts.Posts)
}
