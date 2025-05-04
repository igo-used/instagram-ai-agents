package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/igo-used/instagram-ai-agents/internal/agents"
	"github.com/igo-used/instagram-ai-agents/internal/database"
	"github.com/igo-used/instagram-ai-agents/internal/instagram"
)

func main() {
	// Initialize database
	db, err := database.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve static files
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// Web UI routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Instagram AI Agents",
		})
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "Dashboard | Instagram AI Agents",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Tech Trend Analyzer routes
		api.GET("/tech-trends", func(c *gin.Context) {
			analyzer, err := agents.NewTechTrendAnalyzer()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			news, err := analyzer.FetchTechNews()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   news,
			})
		})

		api.GET("/content-ideas", func(c *gin.Context) {
			analyzer, err := agents.NewTechTrendAnalyzer()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			news, err := analyzer.FetchTechNews()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			ideas, err := analyzer.GenerateContentIdeas(news)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   ideas,
			})
		})

		// Sarcasm Enhancer routes
		api.POST("/enhance-content", func(c *gin.Context) {
			var req struct {
				Content      string `json:"content" binding:"required"`
				SarcasmLevel int    `json:"sarcasmLevel" binding:"required,min=1,max=10"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			enhancer, err := agents.NewSarcasmEnhancer()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			enhanced, err := enhancer.EnhanceContent(req.Content, req.SarcasmLevel)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   enhanced,
			})
		})

		// Behind the Scenes Speculator routes
		api.GET("/companies", func(c *gin.Context) {
			speculator, err := agents.NewBehindScenesSpeculator()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			companies := speculator.ListAvailableCompanies()

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   companies,
			})
		})

		api.GET("/topics/:company", func(c *gin.Context) {
			company := c.Param("company")

			speculator, err := agents.NewBehindScenesSpeculator()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			topics, err := speculator.GenerateTopics(company)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   topics,
			})
		})

		api.POST("/speculate", func(c *gin.Context) {
			var req struct {
				Company string `json:"company" binding:"required"`
				Topic   string `json:"topic" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			speculator, err := agents.NewBehindScenesSpeculator()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			result, err := speculator.GenerateSpeculation(req.Company, req.Topic)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   result,
			})
		})

		// Instagram routes
		api.GET("/instagram/media", func(c *gin.Context) {
			client, err := instagram.NewClient()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			media, err := client.GetRecentMedia()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   media,
			})
		})

		api.GET("/instagram/insights/:mediaId", func(c *gin.Context) {
			mediaID := c.Param("mediaId")

			client, err := instagram.NewClient()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			insights, err := client.GetMediaInsights(mediaID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   insights,
			})
		})

		// Database routes
		api.GET("/content-ideas/db", func(c *gin.Context) {
			ideas, err := db.GetContentIdeas()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   ideas,
			})
		})

		api.POST("/content-ideas/db", func(c *gin.Context) {
			var idea database.ContentIdea

			if err := c.ShouldBindJSON(&idea); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			err := db.SaveContentIdea(&idea)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   idea,
			})
		})

		api.GET("/posts", func(c *gin.Context) {
			posts, err := db.GetPosts()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   posts,
			})
		})

		api.POST("/posts", func(c *gin.Context) {
			var post database.Post

			if err := c.ShouldBindJSON(&post); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			err := db.SavePost(&post)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   post,
			})
		})

		api.GET("/analytics/:postId", func(c *gin.Context) {
			postIDStr := c.Param("postId")
			postID, err := strconv.Atoi(postIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid post ID",
				})
				return
			}

			analytics, err := db.GetAnalyticsForPost(postID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   analytics,
			})
		})
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}
