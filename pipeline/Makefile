build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/convert_posts_to_json cmd/convert_posts_to_json/main.go
	zip ../bin/convert_posts_to_json.zip ../bin/convert_posts_to_json
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/upload_unchanged_posts cmd/upload_unchanged_posts/main.go
	zip ../bin/upload_unchanged_posts.zip ../bin/upload_unchanged_posts
	env GOOS=linux go build -ldflags="-s -w" -o ../bin/save_posts cmd/save_posts/main.go
	zip ../bin/save_posts.zip ../bin/save_posts