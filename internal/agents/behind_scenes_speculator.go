package agents

import (
	"fmt"
	"os"
	"strings"
)

// BehindScenesSpeculator generates speculative "insider" content about tech companies
type BehindScenesSpeculator struct {
	OpenAIKey string
	Companies []string
}

// SpeculationResult represents the result of a speculation
type SpeculationResult struct {
	Company     string   `json:"company"`
	Headline    string   `json:"headline"`
	Speculation string   `json:"speculation"`
	Disclaimer  string   `json:"disclaimer"`
	Sources     []string `json:"sources"`
}

// NewBehindScenesSpeculator creates a new behind scenes speculator
func NewBehindScenesSpeculator() (*BehindScenesSpeculator, error) {
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	// Default list of tech companies to speculate about
	companies := []string{
		"Apple", "Google", "Microsoft", "Meta", "Amazon",
		"Tesla", "Twitter", "Netflix", "Spotify", "Uber",
	}

	return &BehindScenesSpeculator{
		OpenAIKey: openAIKey,
		Companies: companies,
	}, nil
}

// GenerateSpeculation generates speculative content about a tech company
func (b *BehindScenesSpeculator) GenerateSpeculation(company string, topic string) (*SpeculationResult, error) {
	// Validate company
	companyValid := false
	for _, c := range b.Companies {
		if strings.EqualFold(c, company) {
			companyValid = true
			break
		}
	}

	if !companyValid {
		return nil, fmt.Errorf("invalid company: %s", company)
	}

	// In a real implementation, you would call the OpenAI API to generate the speculation
	// For now, we'll return mock data

	// Generate a headline based on the company and topic
	headline := fmt.Sprintf("What's Really Happening Inside %s's %s Division", company, topic)

	// Generate the speculation
	speculation := fmt.Sprintf(`
# Behind the Scenes at %s: %s

According to sources familiar with the matter, %s has been working on a secret project in their %s division. 
While the company has publicly stated that they're focused on incremental improvements, our insider sources 
suggest that there's much more happening behind closed doors.

## What We've Heard

1. The team has reportedly been working nights and weekends on a breakthrough technology
2. Several key engineers were pulled from other projects to work on this initiative
3. The company has been quietly acquiring smaller startups with expertise in this area

## Why This Matters

If these rumors are true, %s could be positioning themselves to disrupt the entire industry. 
Competitors should be worried, as this move could significantly shift market dynamics.

## Timeline

Our sources suggest we might see an announcement within the next 6-12 months, though the 
company is known for keeping projects under wraps until the last minute.
`, company, topic, company, topic, company)

	// Add a disclaimer
	disclaimer := "DISCLAIMER: This content is speculative and based on rumors and analysis. It should not be taken as confirmed fact or used for investment decisions."

	// Add some mock sources
	sources := []string{
		"Anonymous industry insiders",
		"Pattern analysis of recent job postings",
		"Supply chain observations",
		"Recent patent filings",
	}

	result := &SpeculationResult{
		Company:     company,
		Headline:    headline,
		Speculation: speculation,
		Disclaimer:  disclaimer,
		Sources:     sources,
	}

	return result, nil
}

// ListAvailableCompanies returns the list of companies that can be speculated about
func (b *BehindScenesSpeculator) ListAvailableCompanies() []string {
	return b.Companies
}

// GenerateTopics generates potential topics for speculation based on a company
func (b *BehindScenesSpeculator) GenerateTopics(company string) ([]string, error) {
	// Validate company
	companyValid := false
	for _, c := range b.Companies {
		if strings.EqualFold(c, company) {
			companyValid = true
			break
		}
	}

	if !companyValid {
		return nil, fmt.Errorf("invalid company: %s", company)
	}

	// In a real implementation, you would call the OpenAI API to generate topics
	// For now, we'll return mock data based on the company

	topics := []string{}

	switch strings.ToLower(company) {
	case "apple":
		topics = []string{"AR/VR Headset", "Electric Vehicle", "AI Strategy", "Next iPhone"}
	case "google":
		topics = []string{"Quantum Computing", "Search Algorithm", "Android Future", "AI Ethics"}
	case "microsoft":
		topics = []string{"Windows Next Gen", "Azure Strategy", "Gaming Division", "OpenAI Partnership"}
	case "meta":
		topics = []string{"Metaverse Plans", "AI Research", "Content Moderation", "VR Hardware"}
	case "amazon":
		topics = []string{"AWS Expansion", "Retail Strategy", "Logistics Innovation", "Space Ventures"}
	case "tesla":
		topics = []string{"Full Self-Driving", "Battery Technology", "Robotics Division", "Mars Plans"}
	case "twitter":
		topics = []string{"Algorithm Changes", "Monetization Strategy", "Content Policies", "User Growth"}
	case "netflix":
		topics = []string{"Content Strategy", "Gaming Expansion", "AI Recommendations", "Ad-Supported Tier"}
	case "spotify":
		topics = []string{"Podcast Strategy", "Creator Tools", "AI Music Generation", "Subscription Models"}
	case "uber":
		topics = []string{"Autonomous Vehicles", "Delivery Expansion", "Labor Relations", "Urban Mobility"}
	default:
		topics = []string{"New Products", "R&D Initiatives", "Leadership Changes", "Strategic Pivots"}
	}

	return topics, nil
}
