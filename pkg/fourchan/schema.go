package fourchan

import (
	"encoding/json"
	//"fmt"
)

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

	f := t.Posts[0]

	j, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return "Invalid json object"
	}
	
	return string(j)
	//return fmt.Sprintf("Thread with %d posts", t.PostCount())
}

type Post struct {
	// The numeric post ID
	No            int  `json:"no"`
	// For replies: this is the ID of the thread being replied to. For OP: this value is zero
	Resto         int  `json:"resto"`
	// If the thread is being pinned to the top of the page (OP only, if thread is currently stickied)	
	Sticky        int    `json:"sticky,omitempty"`
	// If the thread is closed to replies (OP only, if thread is currently closed)
	Closed        int    `json:"closed,omitempty"`
	// MM/DD/YY(Day)HH:MM (:SS on some boards), EST/EDT timezone
	Now           string `json:"now"`
	// UNIX timestamp the post was created
	Time          int  `json:"time"`
	// Name user posted with. Defaults to Anonymous
	Name          string `json:"name"`
	// The user's tripcode, in format: !tripcode or !!securetripcode (if post has tripcode)
	Trip		  string `json:"trip,omitempty"`
	// The poster's ID (if post has ID)
	Id            string `json:"id,omitempty"`
	// The capcode identifier for a post (if post has capcode)
	Capcode       string `json:"capcode"`
	// 	Poster's ISO 3166-1 alpha-2 country code (if country flags are enabled)
	Country       string `json:"country,omitempty"`
	// Poster's country name (if country flags are enabled)
	CountryName   string `json:"country_name,omitempty"`
	// Poster's board flag code (if board flags are enabled)
	BoardFlag	 string `json:"board_flag,omitempty"`
	// Poster's board flag name (if board flags are enabled)
	FlagName      string `json:"flag_name,omitempty"`
	// OP Subject text (OP only, if subject was included)
	Sub           string `json:"sub,omitempty"`
	// Comment (HTML escaped) (if comment was included)
	Com           string `json:"com,omitempty"`
	// Unix timestamp + microtime that an image was uploaded	
	Tim           int  `json:"tim,omitempty"`
	// Filename as it appeared on the poster's device
	Filename      string `json:"filename"`
	// Filetype
	Ext           string `json:"ext,omitempty"`
	// Size of uploaded file in bytes
	Fsize         int    `json:"fsize,omitempty"`
	// 24 character, packed base64 MD5 hash of file
	Md5           string `json:"md5,omitempty"`
	// Image width dimension
	W             int    `json:"w,omitempty"`
	// Image height dimension
	H             int    `json:"h,omitempty"`
	// Thumbnail image width dimension
	TnW           int    `json:"tn_w,omitempty"`
	// Thumbnail image height dimension
	TnH           int    `json:"tn_h,omitempty"`
	// If the file was deleted (if post had attachment and attachment is deleted)
	FileDeleted   int    `json:"filedeleted,omitempty"`
	// If the image was spoilered or not (if post has attachment and attachment is spoilered)
	Spoiler       int    `json:"spoiler,omitempty"`
	// The custom spoiler ID for a spoilered image (1-10 or not set)
	CustomSpoiler int    `json:"custom_spoiler,omitempty"`
	// Total number of replies to a thread (OP only)
	Replies       int    `json:"replies,omitempty"`
	// Total number of image replies to a thread (OP only)
	Images        int    `json:"images,omitempty"`
	// If a thread has reached bumplimit, it will no longer bump (OP only, only if bump limit has been reached)
	BumpLimit     int    `json:"bumplimit,omitempty"`
	// If an image has reached image limit, no more image replies can be made (OP only, only if image limit has been reached)
	ImageLimit    int    `json:"imagelimit,omitempty"`
	// The category of .swf upload (OP only, /f/ only)
	Tag           string `json:"tag,omitempty"`
	// SEO URL slug for thread (OP only)
	SemanticURL   string `json:"semantic_url,omitempty"`
	// Year 4chan pass bought (if poster put 'since4pass' in the options field)
	Since4pass    int    `json:"since4pass,omitempty"`
	// Number of unique posters in the thread (OP only, only if post has NOT been archived)
	UniqueIPs     int    `json:"unique_ips,omitempty"`
	// Mobile optimized image exists for post
	MImg		  int    `json:"m_img,omitempty"`
	// Thread has reached the board's archive (OP only, if thread has been archived)
	Archived      int    `json:"archived,omitempty"`
	// UNIX timestamp the post was archived (OP only, if thread has been archived)
	ArchivedOn    int  `json:"archived_on,omitempty"`

	// Number of replies minus the number of previewed replies (OP only, Page only)
	OmittedPosts  int    `json:"omitted_posts,omitempty"`
	// Number of image replies minus the number of previewed image replies	(OP only, Page only)
	OmittedImages int    `json:"omitted_images,omitempty"`
}

func (p Page) ThreadCount() int {
	return len(p.Threads)
}
