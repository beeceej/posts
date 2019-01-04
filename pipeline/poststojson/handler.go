package poststojson

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/beeceej/posts/pipeline/shared/state"

	"github.com/beeceej/posts/pipeline/shared/inflight"

	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Handler is the entrypoint into the posts-to-json service layer logic
type Handler struct {
	s3iface.S3API
	PostsRepositoryURL string
	Inflight           *inflight.Inflight
}

// Handle will clone the contents of the git repository into memory
// then for all markdown files extract the metadata from them,
// and write them to S3 where the next step will take over
func (h *Handler) Handle(event interface{}) (*state.InflightRef, error) {
	fmt.Println("Handle Begin")
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: h.PostsRepositoryURL,
	})
	if err != nil {
		return nil, err
	}
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()

	if err != nil {
		return nil, err
	}

	postTree, err := tree.Tree("posts")

	if err != nil {
		return nil, err
	}

	converter := new(postConverter)

	if err != nil {
		return nil, err
	}

	postTree.Files().ForEach(converter.convert)
	
	b, err := json.Marshal(converter.posts)
	if err != nil {
		return nil, err
	}

	return h.Inflight.Write(bytes.NewReader(b))
}
