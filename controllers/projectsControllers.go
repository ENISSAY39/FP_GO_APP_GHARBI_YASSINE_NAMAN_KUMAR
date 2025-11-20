package controllers

import (
	"net/http"
	"strings"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	body.Name = strings.TrimSpace(body.Name)
	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	currentUser, ok := uAny.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user in context"})
		return
	}

	project := models.Project{
		Name:        body.Name,
		Description: strings.TrimSpace(body.Description),
		OwnerID:     currentUser.ID,
	}

	if err := initializers.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}

	// Re-query to preload Owner and Tasks
	if err := initializers.DB.Preload("Owner").Preload("Tasks").First(&project, project.ID).Error; err != nil {
		// fallback: still return created project without relations
		c.JSON(http.StatusCreated, gin.H{"project": project})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"project": project})
}

// ProjectsIndex returns all projects (preloads owner and tasks)
func ProjectsIndex(c *gin.Context) {
	var projects []models.Project
	if err := initializers.DB.Preload("Owner").Preload("Tasks").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query projects"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// ProjectsShow returns a single project by id (with owner and tasks)
func ProjectsShow(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := initializers.DB.Preload("Owner").Preload("Tasks").First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	var project models.Project
	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	currentUser, ok := uAny.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user in context"})
		return
	}
	if project.OwnerID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	project.Name = strings.TrimSpace(body.Name)
	project.Description = strings.TrimSpace(body.Description)

	if err := initializers.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update project"})
		return
	}

	// preload relations for response
	if err := initializers.DB.Preload("Owner").Preload("Tasks").First(&project, project.ID).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"project": project})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// ProjectsDelete deletes a project (owner only)
func ProjectsDelete(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	currentUser, ok := uAny.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user in context"})
		return
	}
	if project.OwnerID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if err := initializers.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete project"})
		return
	}
	c.Status(http.StatusNoContent)
}
