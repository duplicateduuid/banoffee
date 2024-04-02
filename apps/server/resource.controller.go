package main

import (
	"fmt"
	"strconv"

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
		resource, err := s.repositories.resourceRepository.CreateResource(resource)

		if err != nil {
			fmt.Printf("[ERROR] [API.CreateResource] failed to fetch resources: %s", err)
			ctx.JSON(400, gin.H{"error": "Cannot create resource"})
			return
		}

		ctx.JSON(200, gin.H{"resource": resource})
	}
}

type SearchResourceRequest struct {
	Name   string `validate:"required" tstype:"string"`
	Limit  int    `tstype:"number"`
	Offset int    `tstype:"number"`
}

type SearchResourceResponse struct {
	Resources []Resource `json:"resources" tstype:"Resource"`
}

func (s *API) handleSearchResource() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, err := strconv.Atoi(ctx.Query("limit"))
		offset, err := strconv.Atoi(ctx.Query("offset"))

		if err != nil {
			ctx.JSON(400, gin.H{"error": fmt.Sprintf("invalid pagination")})
			return
		}

		if limit == 0 {
			limit = 10
		}

		req := SearchResourceRequest{
			Name:   ctx.Query("name"),
			Limit:  limit,
			Offset: offset,
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			ctx.JSON(400, gin.H{"error": fmt.Sprintf("valiadtion errors: %s", errors)})
			return
		}

		resources, err := s.repositories.resourceRepository.SearchByName(req.Name, req.Limit, req.Offset)
		if err != nil {
			ctx.JSON(400, gin.H{"error": fmt.Sprintf("failed to search resources: %s", err)})
			return
		}

		response := SearchResourceResponse{Resources: *resources}
		ctx.JSON(200, response)
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
