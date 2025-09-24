package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type BasketControllerRouter interface {
	RegisterRoutes(router gin.IRouter) error
}

var _ BasketControllerRouter = (*BasketControllerRouterImpl)(nil)

type BasketControllerRouterImpl struct {
	basketController BasketController
}

func NewBasketControllerRouter(basketController BasketController) BasketControllerRouter {
	return &BasketControllerRouterImpl{
		basketController: basketController,
	}
}

func (controllerRouter *BasketControllerRouterImpl) RegisterRoutes(router gin.IRouter) error {
	if router == nil {
		return fmt.Errorf("router is nil")
	}

	router.GET("/basket", controllerRouter.basketController.ShowBasket)
	router.DELETE("/basket", controllerRouter.basketController.ClearBasket)
	router.POST("/basket/:productID", controllerRouter.basketController.AddProduct)
	router.POST("/basket/:productID/:count", controllerRouter.basketController.AddProduct)
	router.DELETE("/basket/:productID", controllerRouter.basketController.RemoveProduct)
	router.DELETE("/basket/:productID/:count", controllerRouter.basketController.RemoveProduct)

	return nil
}
