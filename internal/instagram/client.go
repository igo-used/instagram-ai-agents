package instagram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client handles interactions with Instagram Graph API
type Client struct {
	AccessToken string
	UserID      string
	BaseURL     string
}

// MediaInsights represents insights for a media object
type MediaInsights struct {
	ID          string `json:"id"`
	Engagement  int    `json:"engagement"`
	Impressions int    `json:"impressions"`
	Reach       int    `json:"reach"`
	Saved       int    `json:"saved"`
	Timestamp   string `json:"timestamp"`
}

// Media represents an Instagram media object
type Media struct {
	ID        string `json:"id"`
	Caption   string `json:"caption"`
	MediaType string `json:"media_type"`
	MediaURL  string `json:"media_url"`
	Permalink string `json:"permalink"`
	Timestamp string `json:"timestamp"`
	Username  string `json:"username"`
}

// NewClient creates a new Instagram client
func NewClient() (*Client, error) {
	accessToken := os.Getenv("INSTAGRAM_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("INSTAGRAM_ACCESS_TOKEN environment variable not set")
	}

	userID := os.Getenv("INSTAGRAM_USER_ID")
	if userID == "" {
		return nil, fmt.Errorf("INSTAGRAM_USER_ID environment variable not set")
	}

	return &Client{
		AccessToken: accessToken,
		UserID:      userID,
		BaseURL:     "https://graph.instagram.com/v12.0",
	}, nil
}

// GetRecentMedia gets the most recent media from the user's Instagram account
func (c *Client) GetRecentMedia() ([]Media, error) {
	endpoint := fmt.Sprintf("%s/%s/media", c.BaseURL, c.UserID)

	// Build query parameters
	params := url.Values{}
	params.Add("fields", "id,caption,media_type,media_url,permalink,thumbnail_url,timestamp,username")
	params.Add("access_token", c.AccessToken)

	// Make the request
	resp, err := http.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var response struct {
		Data   []Media `json:"data"`
		Paging struct {
			Cursors struct {
				Before string `json:"before"`
				After  string `json:"after"`
			} `json:"cursors"`
			Next string `json:"next"`
		} `json:"paging"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetMediaInsights gets insights for a specific media
func (c *Client) GetMediaInsights(mediaID string) (*MediaInsights, error) {
	endpoint := fmt.Sprintf("%s/%s/insights", c.BaseURL, mediaID)

	// Build query parameters
	params := url.Values{}
	params.Add("metric", "engagement,impressions,reach,saved")
	params.Add("access_token", c.AccessToken)

	// Make the request
	resp, err := http.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// We're not using the response body in this mock implementation
	// Just read it to properly close the connection
	_, _ = ioutil.ReadAll(resp.Body)

	// In a real implementation, you would parse the actual response
	// For now, we'll return mock data
	insights := &MediaInsights{
		ID:          mediaID,
		Engagement:  100 + (time.Now().Day() * 10), // Random-ish number
		Impressions: 1000 + (time.Now().Day() * 100),
		Reach:       800 + (time.Now().Day() * 80),
		Saved:       50 + (time.Now().Day() * 5),
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	return insights, nil
}

// PostContent posts content to Instagram (note: this is a simplified example)
// In reality, posting to Instagram requires a Facebook Page, Instagram Business Account,
// and content hosting. This is just a placeholder for the API structure.
func (c *Client) PostContent(caption string, imageURL string) (string, error) {
	// This is a placeholder - in reality, posting to Instagram via API is more complex
	// and requires multiple steps including uploading media to Facebook first
	return fmt.Sprintf("Post created with caption: %s", caption), nil
}
