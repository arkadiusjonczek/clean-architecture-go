package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type BasketControllerRouter interface {
	RegisterRoutes(router *gin.Engine) error
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

func (controllerRouter *BasketControllerRouterImpl) RegisterRoutes(router *gin.Engine) error {
	if router == nil {
		return fmt.Errorf("router is nil")
	}

	router.LoadHTMLGlob("internal/domain/basket/adapters/web/templates/*")
	router.GET("/", controllerRouter.basketController.ShowBasket)

	return nil
}
