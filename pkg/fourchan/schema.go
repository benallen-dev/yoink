package fourchan

import "fmt"

type Page struct {
	Threads []Thread `json:"threads"`
}

type Thread struct {
	Posts []Post `json:"posts"`
}

func (t Thread) PostCount() int {
	return len(t.Posts)
}

func (t Thread) String() string {
	return fmt.Sprintf("Thread with %d posts", t.PostCount())
}

// TODO: Add field descriptions
type Post struct {
	No            int64  `json:"no"`
	Now           string `json:"now"`
	Name          string `json:"name"`
	Sub           string `json:"sub,omitempty"`
	Com           string `json:"com,omitempty"`
	Filename      string `json:"filename"`
	Ext           string `json:"ext,omitempty"`
	W             int    `json:"w,omitempty"`
	H             int    `json:"h,omitempty"`
	TnW           int    `json:"tn_w,omitempty"`
	TnH           int    `json:"tn_h,omitempty"`
	Tim           int64  `json:"tim,omitempty"`
	Time          int64  `json:"time"`
	Md5           string `json:"md5,omitempty"`
	Fsize         int    `json:"fsize,omitempty"`
	Resto         int64  `json:"resto"`
	BumpLimit     int    `json:"bumplimit,omitempty"`
	ImageLimit    int    `json:"imagelimit,omitempty"`
	SemanticURL   string `json:"semantic_url,omitempty"`
	Replies       int    `json:"replies,omitempty"`
	Images        int    `json:"images,omitempty"`
	OmittedPosts  int    `json:"omitted_posts,omitempty"`
	OmittedImages int    `json:"omitted_images,omitempty"`
	Sticky      int    `json:"sticky,omitempty"`
	Closed      int    `json:"closed,omitempty"`
	Capcode     string `json:"capcode"`
	UniqueIPs   int    `json:"unique_ips,omitempty"`

} 

func (p Page) ThreadCount() int {
	return len(p.Threads)
}
