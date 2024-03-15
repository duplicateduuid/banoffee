package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

type SearchResourceRequest struct {
	Name string `json:"name" validate:"required" tstype:"string"`
}

func (s *API) handleSearchResource() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Query("name")
		req := SearchResourceRequest{Name: name}

		validate := validator.New()
		err := validate.Struct(req)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			ctx.JSON(400, gin.H{"error": fmt.Sprintf("valiadtion errors: %s", errors)})
			return
		}

		resources, err := s.repositories.resourceRepository.SearchByName(req.Name)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "failed to search resources"})
			return
		}

		ctx.JSON(200, gin.H{"resources": resources})
	}
}

type GetResourcePayload struct {
	Id   string `db:"id" form:"id"`
	Url  string `db:"url" form:"url"`
	Name string `db:"name" form:"name"`
}

func (s *API) handleGetResource() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := GetResourcePayload{}

		if ctx.ShouldBindQuery(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if req.Id != "" {
			resource, err := s.repositories.resourceRepository.GetResourceById(req.Id)

			if err != nil {
				ctx.JSON(400, gin.H{"error": "Invalid ID"})
				return
			}

			ctx.JSON(200, gin.H{"resource": resource})
			return
		}

		if req.Url != "" {
			resource, err := s.repositories.resourceRepository.GetResourceByUrl(req.Url)

			if err != nil {
				ctx.JSON(400, gin.H{"error": "Invalid URL"})
				return
			}

			ctx.JSON(200, gin.H{"resource": resource})
			return
		}

		if req.Name != "" {
			resource, err := s.repositories.resourceRepository.GetResourceByName(req.Name)

			if err != nil {
				ctx.JSON(400, gin.H{"error": "Invalid Name"})
				return
			}

			ctx.JSON(200, gin.H{"resource": resource})
			return
		}

		ctx.JSON(400, gin.H{"error": "Invalid input"})
	}
}

type GetMyResourcesPayload struct {
	Limit  int `db:"limit" form:"limit"`
	Offset int `db:"offset" form:"offset"`
}

func (s *API) handleGetMyResources() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := GetMyResourcesPayload{}

		if ctx.ShouldBindQuery(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			ctx.Abort()
			return
		}

		user := ctx.MustGet("user").(*User)

		resources, err := s.repositories.resourceRepository.GetUserResources(user, req.Limit, req.Offset)

		if err != nil {
			fmt.Println(err)
			ctx.JSON(400, gin.H{"error": "Cannot retrieve resources"})
			ctx.Abort()
			return
		}

		if len(*resources) <= 0 {
			ctx.JSON(200, gin.H{"resources": []*Resource{}})
			ctx.Abort()
			return
		}

		ctx.JSON(200, gin.H{"resources": resources})
	}
}
