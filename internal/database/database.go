package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is a wrapper around sql.DB
type DB struct {
	*sql.DB
}

// New creates a new database connection
func New() (*DB, error) {
	// Get database connection string from environment variable
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback to individual connection parameters
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		if host == "" || port == "" || user == "" || dbname == "" {
			return nil, fmt.Errorf("database connection parameters not set")
		}

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// Initialize creates the necessary tables if they don't exist
func (db *DB) Initialize() error {
	// Create content_ideas table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS content_ideas (
			id SERIAL PRIMARY KEY,
			headline TEXT NOT NULL,
			content TEXT NOT NULL,
			talking_points TEXT[] NOT NULL,
			hashtags TEXT[] NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}

	// Create posts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			instagram_id TEXT,
			caption TEXT NOT NULL,
			media_url TEXT,
			permalink TEXT,
			status TEXT NOT NULL,
			scheduled_at TIMESTAMP,
			posted_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}

	// Create analytics table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS analytics (
			id SERIAL PRIMARY KEY,
			post_id INTEGER REFERENCES posts(id),
			engagement INTEGER NOT NULL,
			impressions INTEGER NOT NULL,
			reach INTEGER NOT NULL,
			saved INTEGER NOT NULL,
			recorded_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// ContentIdea represents a content idea in the database
type ContentIdea struct {
	ID            int       `json:"id"`
	Headline      string    `json:"headline"`
	Content       string    `json:"content"`
	TalkingPoints []string  `json:"talking_points"`
	Hashtags      []string  `json:"hashtags"`
	CreatedAt     time.Time `json:"created_at"`
}

// SaveContentIdea saves a content idea to the database
func (db *DB) SaveContentIdea(idea *ContentIdea) error {
	query := `
    INSERT INTO content_ideas (headline, content, talking_points, hashtags)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at
	`

	return db.QueryRow(
		query,
		idea.Headline,
		idea.Content,
		idea.TalkingPoints,
		idea.Hashtags,
	).Scan(&idea.ID, &idea.CreatedAt)
}

// GetContentIdeas gets all content ideas from the database
func (db *DB) GetContentIdeas() ([]ContentIdea, error) {
	query := `
		SELECT id, headline, content, talking_points, hashtags, created_at
		FROM content_ideas
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ideas []ContentIdea
	for rows.Next() {
		var idea ContentIdea
		err := rows.Scan(
			&idea.ID,
			&idea.Headline,
			&idea.Content,
			&idea.TalkingPoints,
			&idea.Hashtags,
			&idea.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		ideas = append(ideas, idea)
	}

	return ideas, nil
}

// Post represents a post in the database
type Post struct {
	ID          int        `json:"id"`
	InstagramID string     `json:"instagram_id"`
	Caption     string     `json:"caption"`
	MediaURL    string     `json:"media_url"`
	Permalink   string     `json:"permalink"`
	Status      string     `json:"status"` // draft, scheduled, posted
	ScheduledAt *time.Time `json:"scheduled_at"`
	PostedAt    *time.Time `json:"posted_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// SavePost saves a post to the database
func (db *DB) SavePost(post *Post) error {
	query := `
    INSERT INTO posts (instagram_id, caption, media_url, permalink, status, scheduled_at, posted_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, created_at
`

	return db.QueryRow(
		query,
		post.InstagramID,
		post.Caption,
		post.MediaURL,
		post.Permalink,
		post.Status,
		post.ScheduledAt,
		post.PostedAt,
	).Scan(&post.ID, &post.CreatedAt)
}

// GetPosts gets all posts from the database
func (db *DB) GetPosts() ([]Post, error) {
	query := `
		SELECT id, instagram_id, caption, media_url, permalink, status, scheduled_at, posted_at, created_at
		FROM posts
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.InstagramID,
			&post.Caption,
			&post.MediaURL,
			&post.Permalink,
			&post.Status,
			&post.ScheduledAt,
			&post.PostedAt,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Analytics represents analytics data in the database
type Analytics struct {
	ID          int       `json:"id"`
	PostID      int       `json:"post_id"`
	Engagement  int       `json:"engagement"`
	Impressions int       `json:"impressions"`
	Reach       int       `json:"reach"`
	Saved       int       `json:"saved"`
	RecordedAt  time.Time `json:"recorded_at"`
}

// SaveAnalytics saves analytics data to the database
func (db *DB) SaveAnalytics(analytics *Analytics) error {
	query := `
    INSERT INTO analytics (post_id, engagement, impressions, reach, saved)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, recorded_at
`

	return db.QueryRow(
		query,
		analytics.PostID,
		analytics.Engagement,
		analytics.Impressions,
		analytics.Reach,
		analytics.Saved,
	).Scan(&analytics.ID, &analytics.RecordedAt)
}

// GetAnalyticsForPost gets all analytics data for a post
func (db *DB) GetAnalyticsForPost(postID int) ([]Analytics, error) {
	query := `
    SELECT id, post_id, engagement, impressions, reach, saved, recorded_at
    FROM analytics
    WHERE post_id = $1
    ORDER BY recorded_at DESC
`

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyticsData []Analytics
	for rows.Next() {
		var data Analytics
		err := rows.Scan(
			&data.ID,
			&data.PostID,
			&data.Engagement,
			&data.Impressions,
			&data.Reach,
			&data.Saved,
			&data.RecordedAt,
		)
		if err != nil {
			return nil, err
		}
		analyticsData = append(analyticsData, data)
	}

	return analyticsData, nil
}

// Add these methods to database.go

// SaveSpeculation saves a speculation to the database
func (db *DB) SaveSpeculation(company, topic, headline, content string) error {
	query := `
        INSERT INTO speculations (company, topic, headline, content, created_at)
        VALUES ($1, $2, $3, $4, NOW())
    `
	_, err := db.Exec(query, company, topic, headline, content)
	return err
}

// GetRecentSpeculations gets the most recent speculations
func (db *DB) GetRecentSpeculations(limit int) ([]map[string]interface{}, error) {
	query := `
        SELECT id, company, topic, headline, created_at
        FROM speculations
        ORDER BY created_at DESC
        LIMIT $1
    `

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var company, topic, headline string
		var createdAt time.Time

		if err := rows.Scan(&id, &company, &topic, &headline, &createdAt); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"id":        id,
			"company":   company,
			"topic":     topic,
			"headline":  headline,
			"createdAt": createdAt.Format(time.RFC3339),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// InitSchema initializes the database schema
func (db *DB) InitSchema() error {
	// Add speculations table to your existing Initialize method
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS speculations (
            id SERIAL PRIMARY KEY,
            company VARCHAR(255) NOT NULL,
            topic VARCHAR(255) NOT NULL,
            headline TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        )
    `)
	if err != nil {
		return err
	}

	// Call your existing Initialize method
	return db.Initialize()
}
