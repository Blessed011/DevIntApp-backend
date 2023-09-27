package api

import (
	"lab1/internal/mdstr"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func search(query string, modules []mdstr.Module) []mdstr.Module {
	var results []mdstr.Module

	query = strings.ToLower(query)

	for _, module := range modules {
		if strings.Contains(strings.ToLower(module.Title), query) {
			results = append(results, module)
		}
	}

	return results
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.LoadHTMLGlob("../../templates/*")

	modules := mdstr.GetModule()

	r.GET("/search", func(c *gin.Context) {
		searchModule := c.Query("moduleName")
		results := search(searchModule, modules)
		if len(results) == 0 {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"modules": modules,
			})
		} else {
			c.HTML(200, "index.tmpl", gin.H{
				"modules":    results,
				"searchName": searchModule,
			})
		}
	})

	r.GET("/full/:page", func(c *gin.Context) {
		page := c.Param("page")

		number, err := strconv.Atoi(page)
		if err != nil || number > len(modules) {
			number = 0
		}

		number -= 1
		c.HTML(http.StatusOK, "cards.tmpl", gin.H{
			"Title":       modules[number].Title,
			"Description": modules[number].Description,
			"Image":       modules[number].Image,
			"LaunchDate":  modules[number].LaunchDate,
			"Mass":        modules[number].Mass,
			"Diameter":    modules[number].Diameter,
			"Length":      modules[number].Length,
		})
	})

	r.Static("/image", "../../resources/images")
	r.Static("/styles", "../../resources/css")

	r.Run()

	log.Println("Server down")
}
