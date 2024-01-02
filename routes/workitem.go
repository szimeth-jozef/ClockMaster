package routes

import (
	"szimeth-jozef/clockmaster/handlers"
	"szimeth-jozef/clockmaster/services/workitem"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func AddWorkItemRoutes(root *echo.Group, db *gorm.DB) {
	workItemHandler := handlers.WorkItemHandler{
		DB:              db,
		WorkItemService: workitem.WorkItemService{DB: db},
	}

	workItemGroup := root.Group("/workitem")

	workItemGroup.GET("", workItemHandler.GetWorkItems)
	workItemGroup.POST("", workItemHandler.CreateWorkItem)
	workItemGroup.PATCH("/:id/start", workItemHandler.StartWorkItem)
	workItemGroup.PATCH("/stop", workItemHandler.StopWorkItem)
	workItemGroup.DELETE("/:id", workItemHandler.DeleteWorkItem)
}
