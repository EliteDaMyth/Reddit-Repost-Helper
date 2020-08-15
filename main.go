package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Reddit Repost Helper V0.0.1")
	first := true
	c := widget.NewVBox(widget.NewLabel("Reddit Repost helper 101"))
	subentry := widget.NewEntry()
	subentry.SetPlaceHolder("Subreddit")
	threshentry := widget.NewEntry()
	threshentry.SetPlaceHolder("Score Threshold (Get Posts above this amount of score)")
	c.Append(subentry)
	c.Append(threshentry)
	c.Append(widget.NewButton("Get New Post", func() {
		log.Println("Subreddit:", subentry.Text)
		log.Println("Score Threshold (Get Posts above this amount of score):", threshentry.Text)
		p := getPost(subentry.Text, "25", threshentry.Text)
		ourl, _ := url.Parse(p.OriginalURL)
		iurl, _ := url.Parse(p.URL)
		if p.Found {
			if first {
				c.Append(widget.NewHyperlink("Original Post Link", ourl))
				c.Append(widget.NewLabel("Original Title: " + p.Title))
				c.Append(widget.NewHyperlink("Image URL", iurl))
				first = false
			} else {
				c.Children[len(c.Children)-3] = widget.NewHyperlink("Original Post Link", ourl)
				c.Children[len(c.Children)-2] = widget.NewLabel("Original Title: " + p.Title)
				c.Children[len(c.Children)-1] = widget.NewHyperlink("Image URL", iurl)
				c.Refresh()
			}
		} else {
			if first {
				c.Append(widget.NewLabel("No Post Found."))
				c.Append(widget.NewLabel("Please try searching for another subreddit"))
				c.Append(widget.NewLabel("or lowering the score threshold"))
				first = false
			} else {
				c.Children[len(c.Children)-3] = widget.NewLabel("No Post Found.")
				c.Children[len(c.Children)-2] = widget.NewLabel("Please try searching for another subreddit")
				c.Children[len(c.Children)-1] = widget.NewLabel("or lowering the score threshold")
			}
		}

	}))

	w.SetContent(c)
	w.ShowAndRun()
}

func getPost(sub, limit, score string) Post {
	// Start Making  URL
	baseURL, _ := url.Parse("https://api.pushshift.io")
	baseURL.Path += "reddit/submission/search/"
	before := fmt.Sprintf("%d", time.Now().AddDate(-1, 0, 0).Unix())
	after := fmt.Sprintf("%d", time.Now().AddDate(-1, -1, 0).Unix())

	params := url.Values{}
	sc := fmt.Sprintf(">%s", score)
	params.Add("score", sc)
	params.Add("limit", limit)
	params.Add("subreddit", sub)
	params.Add("allow_images", "true")
	params.Add("before", before)
	params.Add("after", after)

	baseURL.RawQuery = params.Encode()

	URL := baseURL.String()

	fmt.Printf("Encoded URL is %q\n", URL)
	// End making URL to query

	// Start fetching API
	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var submissions SubmissionListing
	_ = json.Unmarshal(responseData, &submissions)
	rand.Seed(time.Now().Unix())
	if len(submissions.Data) == 0 {
		return Post{"", "", "", false}
	}

	rp := submissions.Data[rand.Intn(len(submissions.Data))]
	fmt.Println(rp)
	return Post{rp.URL, rp.Title, rp.FullLink, true}
}

// TYPES HERE

// Post
type Post struct {
	URL         string
	Title       string
	OriginalURL string
	Found       bool
}

// SubmissionListing .
type SubmissionListing struct {
	Data []Submission `json:"data"`
}

// SubmissionGilding .
type SubmissionGilding struct {
}

// Resolutions : defines resolutions of image
type Resolutions struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

// Source : defines the source

type Source struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

// Variants .
type Variants struct {
}

// Images : Slice of images
type Images struct {
	ID          string        `json:"id"`
	Resolutions []Resolutions `json:"resolutions"`
	Source      Source        `json:"source"`
	Variants    Variants      `json:"variants"`
}

// Preview .
type Preview struct {
	Enabled bool     `json:"enabled"`
	Images  []Images `json:"images"`
}

// Submission .
type Submission struct {
	AllAwardings               []interface{}     `json:"all_awardings"`
	AllowLiveComments          bool              `json:"allow_live_comments"`
	Author                     string            `json:"author"`
	AuthorFlairCSSClass        interface{}       `json:"author_flair_css_class"`
	AuthorFlairRichtext        []interface{}     `json:"author_flair_richtext"`
	AuthorFlairText            interface{}       `json:"author_flair_text"`
	AuthorFlairType            string            `json:"author_flair_type"`
	AuthorFullname             string            `json:"author_fullname"`
	AuthorPatreonFlair         bool              `json:"author_patreon_flair"`
	CanModPost                 bool              `json:"can_mod_post"`
	ContestMode                bool              `json:"contest_mode"`
	CreatedUtc                 int               `json:"created_utc"`
	Domain                     string            `json:"domain"`
	FullLink                   string            `json:"full_link"`
	Gildings                   SubmissionGilding `json:"gildings"`
	ID                         string            `json:"id"`
	IsCrosspostable            bool              `json:"is_crosspostable"`
	IsMeta                     bool              `json:"is_meta"`
	IsOriginalContent          bool              `json:"is_original_content"`
	IsRedditMediaDomain        bool              `json:"is_reddit_media_domain"`
	IsRobotIndexable           bool              `json:"is_robot_indexable"`
	IsSelf                     bool              `json:"is_self"`
	IsVideo                    bool              `json:"is_video"`
	LinkFlairBackgroundColor   string            `json:"link_flair_background_color"`
	LinkFlairRichtext          []interface{}     `json:"link_flair_richtext"`
	LinkFlairTextColor         string            `json:"link_flair_text_color"`
	LinkFlairType              string            `json:"link_flair_type"`
	Locked                     bool              `json:"locked"`
	MediaOnly                  bool              `json:"media_only"`
	NoFollow                   bool              `json:"no_follow"`
	NumComments                int               `json:"num_comments"`
	NumCrossposts              int               `json:"num_crossposts"`
	Over18                     bool              `json:"over_18"`
	Permalink                  string            `json:"permalink"`
	Pinned                     bool              `json:"pinned"`
	RetrievedOn                int               `json:"retrieved_on"`
	Score                      int               `json:"score"`
	Selftext                   string            `json:"selftext"`
	SendReplies                bool              `json:"send_replies"`
	Spoiler                    bool              `json:"spoiler"`
	Stickied                   bool              `json:"stickied"`
	Subreddit                  string            `json:"subreddit"`
	SubredditID                string            `json:"subreddit_id"`
	SubredditSubscribers       int               `json:"subreddit_subscribers"`
	SubredditType              string            `json:"subreddit_type"`
	Thumbnail                  string            `json:"thumbnail"`
	Title                      string            `json:"title"`
	TotalAwardsReceived        int               `json:"total_awards_received"`
	URL                        string            `json:"url"`
	ParentWhitelistStatus      string            `json:"parent_whitelist_status,omitempty"`
	Pwls                       int               `json:"pwls,omitempty"`
	WhitelistStatus            string            `json:"whitelist_status,omitempty"`
	Wls                        int               `json:"wls,omitempty"`
	PostHint                   string            `json:"post_hint,omitempty"`
	Preview                    Preview           `json:"preview,omitempty"`
	SuggestedSort              string            `json:"suggested_sort,omitempty"`
	ThumbnailHeight            int               `json:"thumbnail_height,omitempty"`
	ThumbnailWidth             int               `json:"thumbnail_width,omitempty"`
	LinkFlairCSSClass          string            `json:"link_flair_css_class,omitempty"`
	LinkFlairTemplateID        string            `json:"link_flair_template_id,omitempty"`
	LinkFlairText              string            `json:"link_flair_text,omitempty"`
	AuthorFlairBackgroundColor string            `json:"author_flair_background_color,omitempty"`
	AuthorFlairTextColor       string            `json:"author_flair_text_color,omitempty"`
	ContentCategories          []string          `json:"content_categories,omitempty"`
}
