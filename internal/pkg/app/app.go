package app

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"lab1/internal/app/config"
	"lab1/internal/app/dsn"
	"lab1/internal/app/redis"
	"lab1/internal/app/repository"
	"lab1/internal/app/role"

	_ "lab1/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
	redisClient *redis.Client
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	api := r.Group("/api")
	{
		modules := api.Group("/modules")
		{
			modules.GET("", app.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), app.GetAllModules)
			modules.GET("/:id", app.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), app.GetModule)
			modules.DELETE("/:id", app.WithAuthCheck(role.Moderator), app.DeleteModule)
			modules.PUT("/:id", app.WithAuthCheck(role.Moderator), app.ChangeModule)
			modules.POST("", app.WithAuthCheck(role.Moderator), app.AddModule)
			modules.POST("/:id/add_to_mission", app.WithAuthCheck(role.Customer, role.Moderator), app.AddToMission)
		}

		missions := api.Group("/missions")
		{
			missions.GET("", app.WithAuthCheck(role.Customer, role.Moderator), app.GetAllMissions)
			missions.GET("/:id", app.WithAuthCheck(role.Customer, role.Moderator), app.GetMission)
			missions.PUT("", app.WithAuthCheck(role.Customer, role.Moderator), app.UpdateMission)
			missions.DELETE("", app.WithAuthCheck(role.Customer, role.Moderator), app.DeleteMission)
			missions.DELETE("/delete_module/:id", app.WithAuthCheck(role.Customer, role.Moderator), app.DeleteFromMission)
			missions.PUT("/user_confirm", app.WithAuthCheck(role.Customer, role.Moderator), app.UserConfirm)
			missions.PUT("/:id/moderator_confirm", app.WithAuthCheck(role.Moderator), app.ModeratorConfirm)
		}

		user := api.Group("/user")
		{
			user.POST("/sign_up", app.Register)
			user.POST("/login", app.Login)
			user.GET("/logout", app.Logout)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("Host and Port: ", app.config.ServiceHost, app.config.ServicePort)
	r.Run(fmt.Sprintf("localhost:80")) //......................................changeable

	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	app.minioClient, err = minio.New(app.config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	app.redisClient, err = redis.New(app.config.Redis)
	if err != nil {
		return nil, err
	}

	return &app, nil
}
