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
	if err = h.Uploader.Upload(postIndex); err != nil {
		return nil, err
	}
	return h.Inflight.Write(b)

}

// HandleSiteMap is
func (h Handler) HandleSiteMap(ref inflight.Ref) (*inflight.Ref, error) {
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
	if err = h.Uploader.UploadSiteMap(postIndex); err != nil {
		return nil, err
	}
	return h.Inflight.Write(b)
}
