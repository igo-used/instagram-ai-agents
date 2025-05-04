# Instagram AI Agents

AI-powered agents for growing a tech commentary Instagram account.

## Features

- **Tech Trend Analyzer**: Identifies emerging tech trends and generates content ideas
- **Sarcasm Enhancer**: Adds witty and sarcastic elements to tech commentary
- **Instagram Integration**: Connects with Instagram to fetch media and insights
- **Analytics Dashboard**: Tracks performance metrics for your content

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL
- OpenAI API key
- News API key
- Instagram Graph API access token

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/igo-used/instagram-ai-agents.git
   cd instagram-ai-agents
   ```

2. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your API keys and database configuration
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

### Using Docker

You can also run the application using Docker:

```bash
docker-compose up
```

## API Endpoints

- : Get the latest tech news
- : Generate content ideas based on tech trends
- : Add sarcastic elements to content
- : Get recent media from Instagram
- : Get insights for a specific media

## Web UI

The application includes a web UI that can be accessed at http://localhost:8080

## License

This project is licensed under the MIT License - see the LICENSE file for details.
