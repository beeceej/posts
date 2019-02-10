package post

type PostGetter interface {
	Get(id, md5 string) (*Post, error)
	// GetAllVersions(id string) ([]*Post, error)
	// BatchGet([]Post) ([]*Post, error)
}

type PostWriter interface {
	Write(post []Post) error
}

type PostRepository interface {
	PostGetter
	PostWriter
}
