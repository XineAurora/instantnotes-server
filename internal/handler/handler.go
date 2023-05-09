package handler

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO: add dependency injection
type Handler struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func New() *Handler {
	router := gin.Default()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB_NAME"),
		os.Getenv("POSTGRES_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	h := &Handler{Router: router, DB: db}

	// TODO: Add middleware to check if user own note or have premissions
	noteManipulation := router.Group("/api")
	{
		noteManipulation.POST("/notes", h.RequireAuth, h.CreateNote)
		noteManipulation.GET("/notes/:id", h.RequireAuth, h.RequirePremisson, h.ReadNote)
		noteManipulation.PUT("/notes/:id", h.RequireAuth, h.RequirePremisson, h.UpdateNote)
		noteManipulation.DELETE("/notes/:id", h.RequireAuth, h.RequirePremisson, h.DeleteNote)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.SignUp)
		auth.POST("/signin", h.SignIn)
	}

	folderManipulation := router.Group("/folder")
	{
		folderManipulation.POST("/", h.RequireAuth, h.CreateFolder)
		folderManipulation.GET("/:id", h.RequireAuth, h.RequireFolderPremisson, h.ReadFolder)
		folderManipulation.GET("/content/:id", h.RequireAuth, h.RequireFolderPremisson, h.ReadFolderContent)
		folderManipulation.PUT("/:id", h.RequireAuth, h.RequireFolderPremisson, h.UpdateFolder)
		folderManipulation.DELETE("/:id", h.RequireAuth, h.RequireFolderPremisson, h.DeleteFolder)
	}

	groupManipulation := router.Group("/group")
	{
		groupManipulation.POST("/", h.RequireAuth, h.CreateGroup)
		groupManipulation.GET("/members/:id", h.RequireAuth, h.RequireMembership, h.GetGroupMembers)
		groupManipulation.POST("/member/:id", h.RequireAuth, h.RequireGroupOwn, h.AddGroupMember)
		groupManipulation.DELETE("/member/:id", h.RequireAuth, h.RequireGroupOwnOrSelf, h.RemoveGroupMember)
		groupManipulation.PUT("/:id", h.RequireAuth, h.RequireGroupOwn, h.UpdateGroup)
		groupManipulation.DELETE("/:id", h.RequireAuth, h.RequireGroupOwn, h.DeleteGroup)
	}
	return h
}
