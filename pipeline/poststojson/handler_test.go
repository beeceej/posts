package poststojson

import "testing"

var handler = &Handler{
	PostsRepositoryURL: "https://github.com/beeceej/iGo",
}

func TestHandler(t *testing.T) {
	handler.Handle(nil)
}
