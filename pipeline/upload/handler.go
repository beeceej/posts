package upload

import (
	"encoding/json"

	"github.com/beeceej/posts/pipeline/shared/post"

	"github.com/beeceej/inflight"
)

// Handler is the entrypoint into the posts-to-json service layer logic
type Handler struct {
	Inflight *inflight.Inflight
	Uploader *Uploader
}

// Handle is
func (h Handler) Handle(ref inflight.Ref) error {
	b, err := h.Inflight.Get(ref.Object)
	if err != nil {
		return err
	}

	posts := new(post.PostIndex)
	if err = json.Unmarshal(b, &posts.Posts); err != nil {
		return err
	}
	return h.Uploader.Upload(posts)
}
