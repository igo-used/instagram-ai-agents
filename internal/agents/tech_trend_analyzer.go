package agents

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// TechTrendAnalyzer is responsible for analyzing tech trends
type TechTrendAnalyzer struct {
	OpenAIKey string
	NewsAPIKey string
}

// NewTechTrendAnalyzer creates a new tech trend analyzer
func NewTechTrendAnalyzer() (*TechTrendAnalyzer, error) {
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	newsAPIKey := os.Getenv("NEWS_API_KEY")
	if newsAPIKey == "" {
		return nil, fmt.Errorf("NEWS_API_KEY environment variable not set")
	}

	return &TechTrendAnalyzer{
		OpenAIKey: openAIKey,
		NewsAPIKey: newsAPIKey,
	}, nil
}

// FetchTechNews fetches the latest tech news
func (t *TechTrendAnalyzer) FetchTechNews() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?category=technology&language=en&apiKey=%s", t.NewsAPIKey)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	articles, ok := result["articles"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var news []map[string]interface{}
	for _, item := range articles {
		article, ok := item.(map[string]interface{})
		if ok {
			news = append(news, article)
		}
	}

	return news, nil
}

// GenerateContentIdeas generates content ideas based on tech news
func (t *TechTrendAnalyzer) GenerateContentIdeas(news []map[string]interface{}) (string, error) {
	// Prepare articles for the prompt
	var articlesText strings.Builder
	for i, article := range news {
		if i >= 5 {
			break // Limit to 5 articles
		}
		title, _ := article["title"].(string)
		description, _ := article["description"].(string)
		articlesText.WriteString(fmt.Sprintf("Title: %s\nDescription: %s\n\n", title, description))
	}

	// This is a placeholder - in a real implementation, you would call the OpenAI API
	// For now, we'll just return a mock response
	return fmt.Sprintf("Content ideas based on:\n\n%s", articlesText.String()), nil
}
