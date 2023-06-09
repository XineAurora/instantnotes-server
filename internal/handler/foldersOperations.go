package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/XineAurora/instantnotes-server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) CreateFolder(c *gin.Context) {
	//get request body
	//TODO: Add parent param
	var body struct {
		Name     string
		GroupID  uint
		ParentID uint
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	currentUser, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var gid *uint
	if body.GroupID == 0 {
		gid = nil
	} else {
		//check if user are in group
		var count int64
		h.DB.Model(&models.GroupMember{}).
			Where("user_id = ? AND group_id = ?", currentUser.(models.User).ID, body.GroupID).
			Count(&count)
		if count != 1 {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		gid = &body.GroupID
	}
	//create folder
	folder := models.Folder{Name: body.Name, UserID: currentUser.(models.User).ID, GroupID: gid}

	res := h.DB.Create(&folder)
	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	//if parent folder isn't 0 create link
	folderLink := models.FolderLink{ParentFolderID: body.ParentID, ChildFolderID: folder.ID}
	res = h.DB.Create(&folderLink)
	if res.Error != nil {
		h.DB.Delete(&folder)
		c.Status(http.StatusBadRequest)
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

func (h *Handler) GetFolderParent(c *gin.Context) {
	id := c.Param("id")

	var folderLink models.FolderLink

	res := h.DB.Where("child_folder_id = ?", id).First(&folderLink)
	if res.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"folderLink": models.FolderLink{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"folderLink": folderLink,
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

	res := h.DB.Delete(&models.Folder{}, id)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
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

	notes, folders, err := getFolderContent(uint(folderID), user.(models.User).ID, h.DB.Session(&gorm.Session{NewDB: true}))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"folders": folders,
		"notes":   notes,
	})
}

func getFolderContent(folderID uint, userID uint, tx *gorm.DB) ([]models.Note, []models.Folder, error) {
	var notes []models.Note
	var folders []models.Folder
	var groupFolders []models.Folder

	if folderID == 0 {
		//find notes
		if tx.Where("folder_id IS NULL AND user_id = ?", userID).Find(&notes).Error != nil {
			return nil, nil, errors.New("db error")
		}
		//find folders without parent
		if tx.Joins("LEFT JOIN folder_links ON folders.id = folder_links.child_folder_id").
			Where("folder_links.parent_folder_id IS NULL AND folders.group_id IS NULL").
			Find(&folders).
			Error != nil {
			return nil, nil, errors.New("db error")
		}
		//find group folders without parent
		if tx.Joins("LEFT JOIN folder_links ON folder_links.child_folder_id = folders.id").
			Where("group_members.user_id = ? AND folder_links.child_folder_id IS NULL", userID).
			Joins("JOIN group_members ON folders.group_id = group_members.group_id").
			Find(&groupFolders).
			Error != nil {
			return nil, nil, errors.New("db error")
		}
		folders = append(folders, groupFolders...)
	} else {
		//find notes
		if tx.Where("folder_id=?", folderID).Find(&notes).Error != nil {
			return nil, nil, errors.New("db error")
		}
		//find inner folders
		if tx.Table("folders").
			Joins("JOIN folder_links ON folder_links.child_folder_id = folders.id").
			Where("folder_links.parent_folder_id = ?", folderID).
			Find(&folders).Error != nil {
			return nil, nil, errors.New("db error")
		}
	}
	return notes, folders, nil
}

func (h *Handler) RequireFolderPremisson(c *gin.Context) {
	//get authorized user
	id := c.Param("id")
	currentUser, exist := c.Get("user")
	if !exist {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	if id == "0" {
		c.Next()
		return
	}

	//get folder
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
