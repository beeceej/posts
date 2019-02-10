package post

import (
	"encoding/json"
	"time"
)

// Post represents the model of a post, for example:
// {
// 	"id": "1",
// 	"title": "Open Source is Awesome",
// 	"author": "beeceej",
// 	"body": "**Open source** enables people who enjoy tinkering to take ownership over what they're using on his or her machine. This rings true whether or not you're a hacker or not. As a developer I am able to easily fork a codebase, look at the internals of it, fix a problem I (or millions others) have and contribute back. If that isn't true power, I don't know what is.\n\nThe cross polination of ideas and methods is what makes open source software such a beautiful thing. As the world continues to build upon software every person should be able to feel empowered to influence the direction of technology, and the world. I am a developer so it is only natural that I enjoy digging into the internals of complex systems, but I know this isn't the case for everyone. Let's work to empower everyone to feel as if they have the ability to influence the technology around them\n\n -beeceej"
// }
type Post struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	NormalizedTitle string    `json:"normalizedTitle"`
	Author          string    `json:"author"`
	Body            string    `json:"body"`
	PostedAt        time.Time `json:"postedAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Visible         bool      `json:"visible"`
	Blurb           string    `json:"blurb"`
	MD5             string    `json:"md5"`
}

// FromBytes accepts a slice of bytes and will attempt to turn it into a Post struct
func (p *Post) FromBytes(b []byte) error {
	return json.Unmarshal(b, p)
}

// ToBytes turns the post into bytes
func (p *Post) ToBytes() []byte {
	b, _ := json.Marshal(p)
	return b
}

//PostIndex represents the model of the post index, for example:
// {
// 	"posts": [
// 	  {
// 		"title": "Open Source is Awesome",
// 		"id": "1",
// 		"blurb": "Open source enables people who enjoy tinkering to take ownership over what they're using on his or her machine...",
// 		"author": "beeceej"
// 	  }
// 	]
//   }
type PostIndex struct {
	Posts []Post `json:"posts"`
}

// FromBytes takes a slice of bytes and will attempt to Unmarshal them into a PostIndex struct
func (p *PostIndex) FromBytes(b []byte) error {
	return json.Unmarshal(b, p)
}
