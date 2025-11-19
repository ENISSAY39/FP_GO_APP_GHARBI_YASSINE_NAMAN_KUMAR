package controllers

import (
	"net/http"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
	"github.com/gin-gonic/gin"
)

// ProjectsCreate creates a project owned by the authenticated user
func ProjectsCreate(c *gin.Context) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user := c.MustGet("user").(models.User)
	project := models.Project{
		Name:        body.Name,
		Description: body.Description,
		OwnerID:     user.ID,
	}
	if err := initializers.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"project": project})
}

// ProjectsIndex returns all projects (preloads tasks)
func ProjectsIndex(c *gin.Context) {
	var projects []models.Project
	initializers.DB.Preload("Tasks").Find(&projects)
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// ProjectsShow returns a single project by id (with tasks)
func ProjectsShow(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := initializers.DB.Preload("Tasks").First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"project": project})
}

// ProjectsUpdate updates a project (owner only)
func ProjectsUpdate(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var project models.Project
	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	user := c.MustGet("user").(models.User)
	if project.OwnerID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not owner"})
		return
	}
	project.Name = body.Name
	project.Description = body.Description
	initializers.DB.Save(&project)
	c.JSON(http.StatusOK, gin.H{"project": project})
}

// ProjectsDelete deletes a project (owner only)
func ProjectsDelete(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	user := c.MustGet("user").(models.User)
	if project.OwnerID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not owner"})
		return
	}
	initializers.DB.Delete(&project)
	c.Status(http.StatusNoContent)
}
