package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goods-manager/internal/domain"
	"goods-manager/internal/domain/entity"
	"net/http"
	"strconv"
)

type GoodController struct {
	goodUsecase domain.GoodUsecase
}

// getGoodFromRequest retrieves a Good entity from the request context.
//
// It expects 'id' and 'projectId' query parameters in the request.
// If any required parameter is missing or if conversion fails, it returns nil.
// If the Good is not found or if there's an internal server error, appropriate JSON responses are sent.
func (g *GoodController) getGoodFromRequest(c *gin.Context) *entity.Good {
	id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return nil
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}

	projectId, ok := c.GetQuery("projectId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "projectId is required"})
		return nil
	}
	projectIdInt, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}

	goodNotFound := func(details string) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    3,
			"message": "errors.good.notFound",
			"details": details,
		})
	}

	good, err := g.goodUsecase.Get(c, idInt)
	if err != nil {
		if errors.Is(err, domain.ErrorGoodNotFound) {
			goodNotFound(err.Error())
			return nil
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil
	}

	if good.ProjectId != projectIdInt {
		goodNotFound("projectId not match")
		return nil
	}

	return good
}

// Create this function is used to create a good.
//
// @Summary		Add a new good to the store
// @Tags		good
// @Accept		json
// @Produce		json
//
// @Param		projectId	query		int				true	"Project ID"
// @Param		good	body		entity.Good			true	"Good object that needs to be added to the store"
//
// @Success		200		{object}	entity.Good			"Good object that was added"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		500		{string}	string				"Server error"
// @Router		/good/create 	[post]
func (g *GoodController) Create(c *gin.Context) {
	var good entity.Good
	if err := c.BindJSON(&good); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if good.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	// if exists projectId param, then set it
	projectId, ok := c.GetQuery("projectId")
	if ok {
		projectIdInt, err := strconv.Atoi(projectId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		good.ProjectId = projectIdInt
	}

	err := g.goodUsecase.Create(c, &good)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, good)
}

// List this function is used for get goods.
//
// @Summary		Get list goods
// @Tags		good
// @Accept		json
// @Produce		json
//
// @Param		offset	query		int				true	"Offset of select"
// @Param		limit	query		int				true	"Limit of rows"
//
// @Success		200		{object}	ListResponse		"Goods objects and metadata"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		500		{string}	string				"Server error"
// @Router		/good/list			[get]
func (g *GoodController) List(c *gin.Context) {
	limit := 10
	limitQuery, ok := c.GetQuery("limit")
	if ok {
		limitInt, err := strconv.Atoi(limitQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		limit = limitInt
	}

	offset := 1
	offsetQuery, ok := c.GetQuery("offset")
	if ok {
		offsetInt, err := strconv.Atoi(offsetQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		offset = offsetInt
	}

	goods, err := g.goodUsecase.List(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := MetaFromGoods(goods, limit, offset)
	c.JSON(200, ListResponse{Meta: meta, Goods: goods})
}

// Update this function update good.
//
// @Summary		Update good
// @Tags		good
// @Accept		json
// @Produce		json
//
// @Param		projectId	query		int				true	"Project ID"
// @Param		id			query		int				true	"ID of good"
// @Param		good		body		entity.Good		true	"Good object that needs update"
//
// @Success		200		{object}	entity.Good			"Good that was updated"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"Good not found"
// @Failure		500		{string}	string				"Server error"
// @Router		/good/update		[patch]
func (g *GoodController) Update(c *gin.Context) {
	var goodUpdate entity.Good
	if err := c.BindJSON(&goodUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if goodUpdate.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	good := g.getGoodFromRequest(c)
	if good == nil {
		return
	}

	good.Name = goodUpdate.Name
	if goodUpdate.Description != "" {
		good.Description = goodUpdate.Description
	}

	err := g.goodUsecase.Update(c, good)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, good)
}

// Delete this function delete good.
//
// @Summary		Delete good
// @Tags		good
// @Produce		json
//
// @Param		projectId	query		int				true	"Project ID"
// @Param		id			query		int				true	"ID of good"
//
// @Success		200		{object}	entity.Good			"Good that was deleted"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"Good not found"
// @Failure		500		{string}	string				"Server error"
// @Router		/good/remove		[delete]
func (g *GoodController) Delete(c *gin.Context) {
	good := g.getGoodFromRequest(c)
	if good == nil {
		return
	}

	err := g.goodUsecase.Delete(c, good)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, good)
}

// Reprioritize this function update good priority.
//
// @Summary		Reprioritize good priority
// @Tags		good
// @Produce		json
//
// @Param		projectId	query		int				true	"Project ID"
// @Param		id			query		int				true	"ID of good"
// @Param		good		body		PrioritizeRequest		true	"New priority"
//
// @Success		200		{object}	PrioritizeResponse	"List goods where was update priority"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"Good not found"
// @Failure		500		{string}	string				"Server error"
// @Router		/good/reprioritiize		[patch]
func (g *GoodController) Reprioritize(c *gin.Context) {
	var priorityRequest PrioritizeRequest
	if err := c.BindJSON(&priorityRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if priorityRequest.NewPriority < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "newPriority must be greater than 0"})
		return
	}

	good := g.getGoodFromRequest(c)
	if good == nil {
		return
	}

	newPriorities, err := g.goodUsecase.Reprioritize(c, good.Id, priorityRequest.NewPriority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := PrioritizeResponseFromMap(newPriorities)
	c.JSON(200, resp)
}

func NewGoodController(goodUsecase domain.GoodUsecase) *GoodController {
	return &GoodController{goodUsecase: goodUsecase}
}
