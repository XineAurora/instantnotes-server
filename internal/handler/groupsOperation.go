package handler

import (
	"net/http"
	"strconv"

	"github.com/XineAurora/instantnotes-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateGroup(c *gin.Context) {
	// get data from the body
	var body struct {
		Name string
	}
	err := c.Bind(&body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// get user info
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	// create the group
	group := models.Group{Name: body.Name, OwnerID: user.(models.User).ID}
	res := h.DB.Create(&group)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	res = h.DB.Create(&models.GroupMember{GroupID: group.ID, Group: group, UserID: user.(models.User).ID})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": 3, "id": group.ID})
		return
	}

	//return it
	c.JSON(http.StatusOK, gin.H{
		"group": group,
	})
}

func (h *Handler) GetGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.Group
	res := h.DB.First(&group, id)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, group)
}

func (h *Handler) GetAllGroups(c *gin.Context) {
	//get user id
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	// find all groups
	var groups []models.Group
	if h.DB.
		Joins("INNER JOIN group_members ON groups.id = group_members.group_id").
		Joins("INNER JOIN users ON group_members.user_id = users.id").
		Where("users.id = ?", user.(models.User).ID).
		Find(&groups).
		Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	//response with them
	c.JSON(http.StatusOK, &groups)
}

func (h *Handler) GetGroupMembers(c *gin.Context) {
	//get group id from request params
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//get all group members
	var users []models.User
	res := h.DB.Table("users").
		Select("users.id, users.name, users.email").
		Joins("INNER JOIN group_members ON users.id = group_members.user_id").
		Where("group_members.group_id = ?", id).
		Scan(&users)

	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	//return it
	c.JSON(http.StatusOK, gin.H{
		"members": users,
	})
}

func (h *Handler) AddGroupMember(c *gin.Context) {
	//get user email from body and groupID from params
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var body struct {
		Email string
	}

	err = c.Bind(&body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//find user in DB
	var user models.User
	res := h.DB.First(&user, "email=?", body.Email)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	//add user to group
	res = h.DB.Create(&models.GroupMember{GroupID: uint(groupID), UserID: user.ID})
	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) RemoveGroupMember(c *gin.Context) {
	//get userID from body and groupID from params
	groupID := c.Param("id")

	var body struct {
		UserID uint
	}
	err := c.Bind(&body)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	//check if user in group
	var count int64
	res := h.DB.Table("group_members").
		Where("user_id = ? AND group_id = ?", body.UserID, groupID).
		Count(&count)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if count == 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	//remove user from group
	res = h.DB.Where("group_id = ? AND user_id = ?", groupID, body.UserID).Delete(&models.GroupMember{})
	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateGroup(c *gin.Context) {
	//get data from body, id from param
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var body struct {
		Name string
	}
	err = c.Bind(&body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//find group
	var group models.Group
	res := h.DB.First(&group, groupID)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//update it
	res = h.DB.Model(&group).Updates(models.Group{Name: body.Name})
	if res.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	//return updated
	c.JSON(http.StatusOK, gin.H{
		"group": group,
	})
}

func (h *Handler) DeleteGroup(c *gin.Context) {
	//should delete group and all folders and notes
	//get groupId from param
	id := c.Param("id")

	res := h.DB.Delete(&models.Group{}, id)
	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) RequireGroupOwn(c *gin.Context) {
	//get authorized user
	id := c.Param("id")
	currentUser, exist := c.Get("user")
	if !exist {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	//get group
	var group models.Group
	res := h.DB.First(&group, id)
	if res.Error != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	//check if user is owner
	if group.OwnerID != currentUser.(models.User).ID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "no premissions",
		})
	}

	c.Next()
}

func (h *Handler) RequireGroupOwnOrSelf(c *gin.Context) {
	//get authorized user
	id := c.Param("id")
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	//get group
	var group models.Group
	res := h.DB.First(&group, id)
	if res.Error != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	var groupMember models.GroupMember
	res = h.DB.First(&groupMember, "group_id = ? AND user_id = ?", id, user.(models.User).ID)
	if res.Error != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	//check if user is owner or itself
	if group.OwnerID != user.(models.User).ID && groupMember.UserID != user.(models.User).ID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "no premissions",
		})
	}

	c.Next()
}

func (h *Handler) RequireMembership(c *gin.Context) {
	//get authorized user
	id := c.Param("id")
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	//check if user in group
	var count int64
	res := h.DB.Table("group_members").
		Where("user_id = ? AND group_id = ?", user.(models.User).ID, id).
		Count(&count)
	if res.Error != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if count == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	c.Next()
}
