package handler

import (
	"net/http"

	"github.com/XineAurora/instantnotes-server/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateNote ...
func (h *Handler) CreateNote(c *gin.Context) {
	//Get user, group, folder, title and content of note from request
	var body struct {
		Title    string
		Content  string
		UserID   uint
		FolderID uint
		GroupID  uint
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
	}
	var fid, gid *uint
	if body.FolderID == 0 {
		fid = nil
	} else {
		fid = &body.FolderID
	}
	if body.GroupID == 0 {
		gid = nil
	} else {
		gid = &body.GroupID
	}
	//Create a note
	note := models.Note{Title: body.Title, Content: body.Content, UserID: body.UserID, FolderID: fid, GroupID: gid}

	res := h.DB.Create(&note)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	//Return it
	c.JSON(http.StatusOK, gin.H{
		"note": note,
	})
}

func (h *Handler) ReadNote(c *gin.Context) {
	id := c.Param("id")
	//TODO add user verification
	var body struct {
		UserId uint
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
	}

	var note models.Note
	res := h.DB.First(&note, id)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"note": note,
	})
}

func (h *Handler) UpdateNote(c *gin.Context) {
	id := c.Param("id")

	//Get user, group, folder, title and content of note from request
	var body struct {
		Title    string
		Content  string
		UserID   uint
		FolderID uint
		GroupID  uint
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
	}
	var fid, gid *uint
	if body.FolderID == 0 {
		fid = nil
	} else {
		fid = &body.FolderID
	}
	if body.GroupID == 0 {
		gid = nil
	} else {
		gid = &body.GroupID
	}
	//Find a note
	var note models.Note
	res := h.DB.First(&note, id)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//Update a note
	h.DB.Model(&note).Updates(models.Note{Title: body.Title, Content: body.Content, UserID: body.UserID, FolderID: fid, GroupID: gid})
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"note": note,
	})
}

func (h *Handler) DeleteNote(c *gin.Context) {
	id := c.Param("id")

	res := h.DB.Delete(&models.Note{}, id)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}
