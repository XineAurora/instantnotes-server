package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/XineAurora/instantnotes-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateFolder(c *gin.Context) {
	//get request body
	var body struct {
		Name    string
		GroupID uint
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	var gid *uint
	if body.GroupID == 0 {
		gid = nil
	} else {
		gid = &body.GroupID
	}
	//TODO add check if user are in group

	//create folder
	currentUser, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	folder := models.Folder{Name: body.Name, UserID: currentUser.(models.User).ID, GroupID: gid}

	res := h.DB.Create(&folder)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	//return folder
	c.JSON(http.StatusOK, gin.H{
		"folder": folder,
	})
}

func (h *Handler) ReadFolder(c *gin.Context) {
	id := c.Param("id")
	var folder models.Folder
	res := h.DB.First(&folder, id)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"folder": folder,
	})
}

func (h *Handler) UpdateFolder(c *gin.Context) {
	//TODO: Add change parent folder
	id := c.Param("id")

	var body struct {
		Name string
	}
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var folder models.Folder
	res := h.DB.First(&folder, id)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	//add channge parent folder
	res = h.DB.Model(&folder).Updates(models.Folder{Name: body.Name})
	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"folder": folder,
	})
}

func (h *Handler) DeleteFolder(c *gin.Context) {
	//delete folder and all information it contains
	id := c.Param("id")
	folderID, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	user, exist := c.Get("user")
	if !exist {
		c.Status(http.StatusInternalServerError)
		return
	}

	res := h.DB.Begin()
	if res.Error != nil {
		res.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{})
	}

	err = h.deleteFolderContent(uint(folderID), user.(models.User).ID)
	if err != nil {
		res.Rollback()
		c.Status(http.StatusInternalServerError)
		return
	} else {
		res.Commit()
	}

	c.Status(http.StatusOK)
}

func (h *Handler) ReadFolderContent(c *gin.Context) {
	id := c.Param("id")
	//if id is 0 it means that need to return content without parent folder

	folderID, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user, exist := c.Get("user")
	if !exist {
		c.Status(http.StatusInternalServerError)
		return
	}

	notes, folders, err := h.getFolderContent(uint(folderID), user.(models.User).ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"folders": folders,
		"notes":   notes,
	})
}

func (h *Handler) getFolderContent(folderID uint, userID uint) ([]models.Note, []models.Folder, error) {
	var notes []models.Note
	var folders []models.Folder
	res := h.DB.Where("folder_id=?", folderID).Find(&notes)
	if res.Error != nil {
		return nil, nil, errors.New("db error")
	}
	if folderID == 0 {
		res = h.DB.Table("folders").
			Joins("LEFT JOIN folder_links ON folder_links.child_folder_id = folders.id").
			Where("folders.user_id = ?", userID).
			Where("NOT EXISTS (SELECT 1 FROM folder_links fl WHERE fl.child_folder_id = folders.id)").
			Where("folder_links.parent_folder_id IS NULL").
			Find(&folders)
		if res.Error != nil {
			return nil, nil, errors.New("db error")
		}
	} else {
		res = h.DB.Table("folders").
			Joins("JOIN folder_links ON folder_links.child_folder_id = folders.id").
			Where("folder_links.parent_folder_id = ?", folderID).
			Find(&folders)
		if res.Error != nil {
			return nil, nil, errors.New("db error")
		}
	}
	return notes, folders, nil
}

func (h *Handler) deleteFolderContent(folderID uint, userID uint) error {
	notes, folders, err := h.getFolderContent(folderID, userID)
	if err != nil {
		return err
	}
	for _, note := range notes {
		if res := h.DB.Delete(&note); res.Error != nil {
			return res.Error
		}
	}
	for _, folder := range folders {
		err = h.deleteFolderContent(folder.ID, userID)
		if err != nil {
			return err
		}
	}
	if res := h.DB.Delete(&models.Folder{ID: folderID}); res.Error != nil {
		return res.Error
	}
	return nil
}

func (h *Handler) RequireFolderPremisson(c *gin.Context) {
	//get authorized user
	id := c.Param("id")
	currentUser, exist := c.Get("user")
	if !exist {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	//get note
	var folder models.Folder
	res := h.DB.First(&folder, id)
	if res.Error != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	//check if user is owner or in group
	if folder.UserID != currentUser.(models.User).ID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "no premissions",
		})
	}

	c.Next()
}
