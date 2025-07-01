package app

import (
	"OnlieStore/internal/auth"
	"OnlieStore/internal/config"
	"OnlieStore/internal/data"
	"OnlieStore/internal/model"
	"OnlieStore/internal/service"
	"OnlieStore/internal/util"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

type App struct {
	orderHandler *service.OrderService
	productStore *service.ProductStore
	userManager  *service.UserManager
	userAuth     *auth.UserAuth
	loader       *data.Loader
}

func NewApp() *App {
	return &App{
		orderHandler: service.NewOrderService(),
		productStore: service.NewProductStore(),
		userManager:  service.NewUserManager(),
		userAuth:     auth.NewUserAuth(config.GetConfig().Secret),
		loader:       data.NewLoader(),
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
	if !app.productStore.IsProductAvailableToBuy(order.ProductID, order.Quantity) {
		err := errors.New("Product is not available in the store to buy ")
		logrus.WithError(err).Error("Failed to add order")
		return err
	}

	// process order
	app.orderHandler.AddOrder(order)

	// update the balance is store
	err := app.productStore.UpdateProductQuantity(order.ProductID, util.ActionProductDecrease, order.Quantity)
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

func (app *App) LoadData() error {
	err := app.loadUsers()
	if err != nil {
		logrus.WithError(err).Error("Failed to load users")
		return err
	}

	err = app.loadProducts()
	if err != nil {
		logrus.WithError(err).Error("Failed to load products")
		return err
	}

	return nil
}

func (app *App) loadUsers() error {
	filePath := fmt.Sprintf("%s/users.csv", config.GetConfig().DataFilePath)
	users, err := app.loader.LoadUsers(filePath)
	if err != nil {
		logrus.WithError(err).Error("Failed to load users")
		return err
	}

	for _, user := range users {
		err = app.userManager.AddUser(user)
		if err != nil {
			logrus.WithError(err).Error("Failed to add user")
			return err
		}

		logrus.WithField("user", user).Info("Added user")
	}

	return nil
}

func (app *App) loadProducts() error {
	filePath := fmt.Sprintf("%s/products.csv", config.GetConfig().DataFilePath)
	products, err := app.loader.LoadProducts(filePath)
	if err != nil {
		logrus.WithError(err).Error("Failed to load users")
		return err
	}

	for _, p := range products {
		app.productStore.AddProduct(p)
		logrus.WithField("product", p).Info("Added product")
	}

	return nil
}
