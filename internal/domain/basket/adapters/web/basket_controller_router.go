package web

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templatesFS embed.FS

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

	templ := template.Must(template.New("").ParseFS(templatesFS, "templates/*.html"))
	router.SetHTMLTemplate(templ)

	router.GET("/", controllerRouter.basketController.ShowBasket)

	return nil
}
