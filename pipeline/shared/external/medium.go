package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/beeceej/posts/pipeline/shared/post"
)

const domain = "https://api.medium.com"

var (
	meEndpoint = fmt.Sprintf("%s/v1/me", domain)
)

func createPostEndpoint(authorID string) string {
	return fmt.Sprintf("%s/v1/users/%s/posts", domain, authorID)
}

// MediumAPI is an http client which interacts with the Medium API
type MediumAPI struct {
	*http.Client
	IntegrationToken string
}

// GetAuthorID hits the Medium API and returns the User ID,
// this userID is needed to make subsequent calls to the API
func (m *MediumAPI) GetAuthorID() (string, error) {
	req, err := http.NewRequest(http.MethodGet, meEndpoint, http.NoBody)
	if err != nil {
		return "", err
	}
	m.addAuthHeader(req)
	response, err := m.Do(req)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var t struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	err = json.Unmarshal(b, &t)
	if err != nil {
		return "", err
	}

	return t.Data.ID, nil
}

func (m *MediumAPI) String() string {
	return ""
}

// CreatePost will convert a post into the medium post domain model
func (m *MediumAPI) CreatePost(post post.Post, authorID string) error {
	var (
		b              []byte
		mediumPostBody struct {
			Title         string   `json:"title"`
			ContentFormat string   `json:"contentFormat"`
			Content       string   `json:"content"`
			Tags          []string `json:"tags"`
			CanonicalURL  string   `json:"canonicalUrl"`
			PublishStatus string   `json:"publishStatus"`
		}
		err error
	)
	mediumPostBody.Title = post.Title
	mediumPostBody.ContentFormat = "markdown"
	mediumPostBody.Content = post.Body
	mediumPostBody.Tags = []string{} // I don't have tags yet :(
	mediumPostBody.CanonicalURL = fmt.Sprintf("https://blog.beeceej.com/blog/%s", post.NormalizedTitle)
	mediumPostBody.PublishStatus = "draft"

	b, err = json.Marshal(mediumPostBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, createPostEndpoint(authorID), bytes.NewReader(b))
	if err != nil {
		return err
	}
	m.addAuthHeader(req)
	resp, err := m.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)

	fmt.Println(string(b))
	return nil
}

func (m *MediumAPI) addAuthHeader(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.IntegrationToken))
	req.Header.Add("Content-Type", "application/json")
}
