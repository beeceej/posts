package saveposts

import (
	"github.com/beeceej/posts/pipeline/shared/post"
)

type PostSaver struct {
	post.PostWriter
	TableName string
}

func (p *PostSaver) SavePosts(posts []post.Post) error {
	return p.PostWriter.Write(posts)
}
