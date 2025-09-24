package rest

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket/usecases"
)

type BasketController interface {
	ShowBasket(c *gin.Context)
	ClearBasket(c *gin.Context)
	AddProduct(c *gin.Context)
	RemoveProduct(c *gin.Context)
}

var _ BasketController = (*BasketControllerImpl)(nil)

type BasketControllerImpl struct {
	usecases.ShowBasketUseCase
	usecases.ClearBasketUseCase
	usecases.AddProductUseCase
	usecases.RemoveProductUseCase
}

func NewBasketController(
	showBasketUseCase usecases.ShowBasketUseCase,
	clearBasketUseCase usecases.ClearBasketUseCase,
	addProductUseCase usecases.AddProductUseCase,
	removeProductUseCase usecases.RemoveProductUseCase,
) *BasketControllerImpl {
	return &BasketControllerImpl{
		ShowBasketUseCase:    showBasketUseCase,
		ClearBasketUseCase:   clearBasketUseCase,
		AddProductUseCase:    addProductUseCase,
		RemoveProductUseCase: removeProductUseCase,
	}
}

func (controller *BasketControllerImpl) ShowBasket(c *gin.Context) {
	userID := getUserID(c)

	output, err := controller.ShowBasketUseCase.Execute(
		&usecases.ShowBasketUseCaseInput{
			UserID: userID,
		},
	)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, output.UserBasket)
}

func (controller *BasketControllerImpl) ClearBasket(c *gin.Context) {
	userID := getUserID(c)

	output, err := controller.ClearBasketUseCase.Execute(
		&usecases.ClearBasketUseCaseInput{
			UserID: userID,
		},
	)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, output.UserBasket)
}

func (controller *BasketControllerImpl) AddProduct(c *gin.Context) {
	userID := getUserID(c)
	productID := c.Param("productID")
	count := c.Param("count")
	if count == "" {
		count = "1"
	}

	countInteger, err := strconv.Atoi(count)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	output, err := controller.AddProductUseCase.Execute(
		&usecases.AddProductUseCaseInput{
			UserID:    userID,
			ProductID: productID,
			Count:     countInteger,
		},
	)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, output.UserBasket)
}

func (controller *BasketControllerImpl) RemoveProduct(c *gin.Context) {
	userID := getUserID(c)
	productID := c.Param("productID")
	count := c.Param("count")
	if count == "" {
		count = "1"
	}

	countInteger, err := strconv.Atoi(count)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	output, err := controller.RemoveProductUseCase.Execute(
		&usecases.RemoveProductUseCaseInput{
			UserID:    userID,
			ProductID: productID,
			Count:     countInteger,
		},
	)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, output.UserBasket)
}
