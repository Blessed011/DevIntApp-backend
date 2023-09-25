package api

import (
	"lab1/internal/mdstr"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func search(query string, cards []mdstr.Card) []mdstr.Card {
	results := []mdstr.Card{}

	for _, card := range cards {
		if strings.Contains(card.Title, query) {
			results = append(results, card)
		}
	}

	return results
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.LoadHTMLGlob("../../templates/*")

	pipe := mdstr.GetModule()

	r.GET("/full", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"pipeline": pipe,
		})
	})

	r.GET("/full/:page", func(c *gin.Context) {
		page := c.Param("page")

		number, err := strconv.Atoi(page)
		if err != nil || number > len(pipe) {
			number = 0
		}

		number -= 1
		c.HTML(http.StatusOK, "cards.tmpl", gin.H{
			"Title":       pipe[number].Title,
			"Description": pipe[number].Description,
			"Image":       pipe[number].Image,
			"LaunchDate":  pipe[number].LaunchDate,
		})
	})

	r.GET("/search", func(c *gin.Context) {
		query := c.Query("query")

		results := search(query, pipe)

		c.HTML(200, "index.tmpl", gin.H{
			"pipeline": results,
		})
	})

	r.Static("/image", "../../resources/images")
	r.Static("/styles", "../../resources/css")

	r.Run()

	log.Println("Server down")
}
