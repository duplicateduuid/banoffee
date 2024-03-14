package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateResourcePayload struct {
	Url         string  `db:"url" json:"url"`
	Name        string  `db:"name" json:"name"`
	ImageUrl    *string `db:"image_url" json:"image_url"`
	Author      *string `db:"author" json:"author"`
	Description *string `db:"description" json:"description"`
}

func (s *API) handleCreateResource() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := CreateResourcePayload{}

		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		resource := NewResource(req.Url, req.Name, req.ImageUrl, req.Author, req.Description)

		err := s.repositories.resourceRepository.CreateResource(resource)

		if err != nil {
			fmt.Println(err)
			ctx.JSON(400, gin.H{"error": "Cannot create resource"})
			return
		}

		ctx.JSON(200, gin.H{"message": "Resource created with success!"})
	}
}

type GetResourcePayload struct {
	Id   *uuid.UUID `db:"id" form:"id"`
	Url  *string    `db:"url" form:"url"`
	Name *string    `db:"name" form:"name"`
}

func (s *API) handleGetResource() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := GetResourcePayload{}

		if ctx.ShouldBindQuery(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if req.Id != nil {
			resource, err := s.repositories.resourceRepository.GetResourceById(*req.Id)

			if err != nil {
				ctx.JSON(400, gin.H{"error": "Invalid URL"})
				return
			}

			ctx.JSON(200, gin.H{"resource": resource})
			return
		}

		if req.Url != nil {
			resource, err := s.repositories.resourceRepository.GetResourceByUrl(*req.Url)

			if err != nil {
				ctx.JSON(400, gin.H{"error": "Invalid URL"})
				return
			}

			ctx.JSON(200, gin.H{"resource": resource})
			return
		}

		if req.Name != nil {
			resource, err := s.repositories.resourceRepository.GetResourceByName(*req.Name)

			if err != nil {
				ctx.JSON(400, gin.H{"error": "Invalid URL"})
				return
			}

			ctx.JSON(200, gin.H{"resource": resource})
			return
		}

		ctx.JSON(400, gin.H{"error": "Invalid input"})
	}
}
