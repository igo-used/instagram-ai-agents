package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/igo-used/instagram-ai-agents/internal/agents"
)

func main() {
	r := gin.Default()
	
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Instagram AI Agents API",
		})
	})
	
	r.GET("/api/tech-trends", func(c *gin.Context) {
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
			"data": news,
		})
	})
	
	r.GET("/api/content-ideas", func(c *gin.Context) {
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
			"data": ideas,
		})
	})
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}
