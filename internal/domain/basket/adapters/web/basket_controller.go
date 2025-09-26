package web

import (
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/adapters/common"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases"
	"github.com/gin-gonic/gin"
)

type BasketController interface {
	ShowBasket(c *gin.Context)
}

type BasketControllerImpl struct {
	usecases.ShowBasketUseCase
}

func NewBasketController(showBasketUseCase usecases.ShowBasketUseCase) BasketController {
	return &BasketControllerImpl{
		ShowBasketUseCase: showBasketUseCase,
	}
}

func (controller *BasketControllerImpl) ShowBasket(c *gin.Context) {
	userID := common.GetUserID()

	output, err := controller.ShowBasketUseCase.Execute(
		&usecases.ShowBasketUseCaseInput{
			UserID: userID,
		},
	)
	if err != nil {
		c.HTML(500, "index.html", gin.H{
			"message": err.Error(),
		})
		return
	}

	c.HTML(200, "index.html", gin.H{
		"userBasket": output.UserBasket,
	})
}
