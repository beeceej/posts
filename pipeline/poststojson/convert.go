package poststojson

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/beeceej/posts/pipeline/shared/domain"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	indexVal         = 4
	indexID          = 0
	indexTitle       = 1
	indexAuthor      = 2
	indexPostedAt    = 3
	indexUpdatedLast = 4
	indexVisible     = 5
)

var (
	// Used to parse out data from the comments in each markdown file
	removeCommentsRegex = regexp.MustCompile(`((\<\!--)(.+):(.+)(--\>)\n)`)
	captureMetaRegex    = regexp.MustCompile(`((\<\!--)(.+):(.+)(--\>))`)

	// Used to normalize titles for pretty URL's
	titleblackList  = regexp.MustCompile(`[ #$%&@/:;<=>?[\]^{|}~“‘+,]`)
	catchDoubleDash = regexp.MustCompile(`-(-)+`)
)

type postConverter struct {
	posts []*domain.Post
}

func (p *postConverter) convert(f *object.File) error {
	var (
		md  string
		err error
	)

	if md, err = f.Contents(); err != nil {
		return err
	}

	p.posts = append(p.posts, toPost(md))

	return nil
}

func toPost(md string) (post *domain.Post) {
	post = captureMeta(md)
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

func captureMeta(md string) *domain.Post {
	matches := captureMetaRegex.FindAllStringSubmatch(md, -1)

	id := strings.TrimSpace(matches[indexID][indexVal])
	title := strings.TrimSpace(matches[indexTitle][indexVal])
	normalizedTitle := normalizeTitle(title)
	author := strings.TrimSpace(matches[indexAuthor][indexVal])
	postedAt := strings.TrimSpace(matches[indexPostedAt][indexVal])
	updatedLast := strings.TrimSpace(matches[indexUpdatedLast][indexVal])
	visible := strings.TrimSpace(matches[indexVisible][indexVal])
	isVisible, err := strconv.ParseBool(visible)

	if err != nil {
		fmt.Println("Couldn't parse value for isVisible, defaulting to false")
		isVisible = false
	}

	return &domain.Post{
		ID:              id,
		Title:           title,
		NormalizedTitle: normalizedTitle,
		Author:          author,
		PostedAt:        postedAt,
		UpdatedAt:       updatedLast,
		Visible:         isVisible,
		MD5:             getmd5(md),
	}
}
