package poststojson

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/beeceej/posts/pipeline/shared/post"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	indexVal     = 4
	indexID      = 0
	indexTitle   = 1
	indexAuthor  = 2
	indexVisible = 3
)

var (
	// Used to parse out data from the comments in each markdown file
	removeCommentsRegex = regexp.MustCompile(`((\<\!--)(.+):(.+)(--\>)\n)`)
	captureMetaRegex    = regexp.MustCompile(`((\<\!--)(.+):(.+)(--\>))`)

	// Used to normalize titles for pretty URL's
	titleblackList  = regexp.MustCompile(`[ #$%&@/:;<=>?[\]^{|}~“‘+,]`)
	catchDoubleDash = regexp.MustCompile(`-(-)+`)
)

// PostConverter handles extracting meta data from the markdown documents
type PostConverter struct {
	post.PostGetter
	posts []*post.Post
}

func (p *PostConverter) convert(f *object.File) error {
	var (
		md  string
		err error
	)

	if md, err = f.Contents(); err != nil {
		return err
	}

	p.posts = append(p.posts, p.toPost(md))

	return nil
}

func (p *PostConverter) toPost(md string) (post *post.Post) {
	var err error
	post, err = p.captureMeta(md)
	if err != nil {
		fmt.Println("Unable to capture meta data for", md)
		panic(err.Error())
	}
	post.Body = removeCommentsRegex.ReplaceAllString(md, "")
	return post
}

func normalizeTitle(title string) (n string) {
	n = titleblackList.ReplaceAllLiteralString(title, "-")
	n = catchDoubleDash.ReplaceAllLiteralString(n, "-")
	n = strings.Replace(n, "!", "", -1)
	return strings.ToLower(n)
}

func getmd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func (p *PostConverter) captureMeta(md string) (*post.Post, error) {
	matches := captureMetaRegex.FindAllStringSubmatch(md, -1)

	id := strings.TrimSpace(matches[indexID][indexVal])
	title := strings.TrimSpace(matches[indexTitle][indexVal])
	normalizedTitle := normalizeTitle(title)
	author := strings.TrimSpace(matches[indexAuthor][indexVal])
	visible := strings.TrimSpace(matches[indexVisible][indexVal])
	isVisible, err := strconv.ParseBool(visible)
	md5hash := getmd5(md)

	if err != nil {
		fmt.Println("Couldn't parse value for isVisible, defaulting to false")
		isVisible = false
	}

	existingPost, err := p.PostGetter.Get(id, md5hash)
	if err != nil {
		return nil, err
	}

	var (
		postedAt    time.Time
		updatedLast time.Time
	)

	if existingPost == nil  { // If it's nil, it didn't exist before
		postedAt = time.Now().UTC()
		updatedLast = time.Now().UTC()
	} else if existingPost.MD5 != md5hash { // Only update it if the hash has changed
		postedAt = existingPost.PostedAt
		updatedLast = time.Now().UTC()
	} else {
		postedAt = existingPost.PostedAt
		updatedLast = existingPost.UpdatedAt
	}

	return &post.Post{
		ID:              id,
		Title:           title,
		NormalizedTitle: normalizedTitle,
		Author:          author,
		PostedAt:        postedAt,
		UpdatedAt:       updatedLast,
		Visible:         isVisible,
		MD5:             md5hash,
	}, nil
}
