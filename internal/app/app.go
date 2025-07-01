package app

import (
	"OnlieStore/internal/auth"
	"OnlieStore/internal/config"
	"OnlieStore/internal/model"
	"OnlieStore/internal/service"
	"OnlieStore/internal/util"
	"errors"
	"github.com/sirupsen/logrus"
)

type App struct {
	orderHandler *service.OrderService
	productStore *service.ProductStore
	userManager  *service.UserManager
	userAuth     *auth.UserAuth
}

func NewApp() *App {
	return &App{
		orderHandler: service.NewOrderService(),
		productStore: service.NewProductStore(),
		userManager:  service.NewUserManager(),
		userAuth:     auth.NewUserAuth(config.GetConfig().Secret),
	}
}

func (app *App) GetProducts(params *model.PaginationParams) ([]*model.Stock, error) {
	products, err := app.productStore.GetProducts(params)
	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve products")
	}

	return products, err
}

func (app *App) AddProduct(product *model.ProductDetails) {
	app.productStore.AddProduct(product)
}

func (app *App) GetOrder(id string) (*model.Order, error) {
	order, err := app.orderHandler.GetOrder(id)
	if err != nil {
		logrus.WithError(err).Error("Failed to get order")
	}

	return order, err
}

func (app *App) AddOrder(order *model.Order) error {
	// validate
	if !app.productStore.IsProductAvailableToBuy(order.ID, order.Quantity) {
		err := errors.New("Product is not available in the store to buy ")
		logrus.WithError(err).Error("Failed to add order")
		return err
	}

	// process order
	app.orderHandler.AddOrder(order)

	// update the balance is store
	err := app.productStore.UpdateProductQuantity(order.ID, util.ActionProductDecrease, order.Quantity)
	if err != nil {
		logrus.WithError(err).Error("Failed to add order")
		// an error has occured. We should update the order status as cancelled
		_ = app.orderHandler.UpdateOrderStatus(order.ID, util.OrderStatusError)
		return err
	}

	return err
}

func (app *App) GenerateJWTToken(userName string, password string) (string, error) {
	user, err := app.userManager.ValidateAndGetUser(userName, password)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate JWT token as user is invalid")
		return "", err
	}

	token, err := app.userAuth.GenerateToken(user.ID, userName)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate JWT token")
	}

	return token, err
}

func (app *App) UpdateOrderStatus(orderId string, status util.OrderStatus) error {
	err := app.orderHandler.UpdateOrderStatus(orderId, status)
	if err != nil {
		logrus.WithError(err).Error("Failed to update order status")
	}

	return err
}
