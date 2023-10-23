package app

import (
	// "lab1/internal/app/config"
	"lab1/internal/app/ds"
	"lab1/internal/app/dsn"
	"lab1/internal/app/repository"

	"net/http"

	"github.com/gin-gonic/gin"

	"log"
)

type Application struct {
	repo *repository.Repository
	// config *config.Config
}

type GetModulesBack struct {
	Modules []ds.Module
	Name    string
}

func (a *Application) Run() {

	r := gin.Default()
	r.LoadHTMLGlob("../../templates/*")

	r.GET("/modules", func(c *gin.Context) {
		name := c.Query("name")
		modules, err := a.repo.GetModuleByName(name)
		if err != nil {
			log.Println("can't get modules", err)
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "index.tmpl", GetModulesBack{
			Name:    name,
			Modules: modules,
		})
	})

	r.GET("/modules/:id", func(c *gin.Context) {
		id := c.Param("id")
		module, err := a.repo.GetModuleByID(id)
		if err != nil {
			log.Printf("can't get module by id %v", err)
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "item.tmpl", *module)
	})

	r.POST("/modules", func(c *gin.Context) {
		id := c.PostForm("delete")

		a.repo.DeleteModule(id)

		modules, err := a.repo.GetModuleByName("")
		if err != nil {
			log.Println("can't get modules", err)
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "index.tmpl", GetModulesBack{
			Name:    "",
			Modules: modules,
		})
	})

	r.Static("/image", "../../resources/images")
	r.Static("/css", "../../static/css")
	r.Run("localhost:8081")
	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	app := Application{}
	// app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	return &app, nil
}
