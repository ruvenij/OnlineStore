package api

import (
	"OnlieStore/internal/api/request"
	"OnlieStore/internal/app"
	"OnlieStore/internal/config"
	"OnlieStore/internal/model"
	"OnlieStore/internal/util"
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Api struct {
	app       *app.App
	echo      *echo.Echo
	validator *validator.Validate
}

func NewApi(app *app.App, e *echo.Echo) *Api {
	return &Api{
		app:       app,
		echo:      e,
		validator: validator.New(),
	}
}

func (api *Api) StartService() {
	logrus.Info("Starting the service at port:", config.GetConfig().Port)
	portAddress := fmt.Sprintf(":%d", config.GetConfig().Port)
	api.echo.Logger.Fatal(api.echo.Start(portAddress))
}

func (api *Api) StopService(ctx context.Context) {
	api.echo.Logger.Fatal(api.echo.Shutdown(ctx))
}

func (api *Api) RegisterFunctions() {
	logrus.Info("Registering the functions")
	// login
	api.echo.POST("/login", api.Login)

	r := api.echo.Group("/api/v1")
	r.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.GetConfig().Secret),
	}))

	// products
	r.GET("/products", api.GetProducts)
	r.POST("/products", api.AddProduct)

	// orders
	r.GET("/orders", api.GetOrder)
	r.POST("/order", api.AddNewOrder)
	r.POST("/status", api.UpdateOrderStatus)

}

func (api *Api) GetProducts(c echo.Context) error {
	limit := c.QueryParam("limit")
	page := c.QueryParam("page")

	logrus.WithFields(logrus.Fields{"path": c.Request().URL.Path, "params": c.QueryParams()}).
		Info("Incoming get products request")

	// validate the request first
	limitInt, pageInt, err := validateGetProductsRequest(limit, page)
	if err != nil {
		logrus.WithError(err).Error("Validation failed for GetProducts request")
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": err.Error()})
	}

	// process the request
	result, err := api.app.GetProducts(&model.PaginationParams{
		Limit: limitInt,
		Page:  pageInt,
	},
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to process get products request")
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}

	logrus.Debug("Retrieved result ", result)
	return c.JSON(http.StatusOK, result)
}

func (api *Api) AddProduct(c echo.Context) error {
	req := new(request.ProductDetails)
	if err := c.Bind(req); err != nil {
		logrus.WithError(err).Error("Failed to bind AddProduct request")
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": err.Error()})
	}

	// validate the request
	product, err := api.validateAddProductRequest(req)
	if err != nil {
		logrus.WithError(err).Error("Validation failed for AddProduct request")
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": err.Error()})
	}

	// process the request
	api.app.AddProduct(product)
	return c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func (api *Api) GetOrder(c echo.Context) error {
	orderId := c.QueryParam("order_id")

	if orderId == "" {
		logrus.WithField("order_id", orderId).Error("Order Id is required")
		err := errors.New("Order Id is required ")
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": err.Error()})
	}

	order, err := api.app.GetOrder(orderId)
	if err != nil {
		logrus.WithError(err).Error("Failed to get order")
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}

	return c.JSON(http.StatusOK, order)
}

func (api *Api) AddNewOrder(c echo.Context) error {
	var req request.Order
	if err := c.Bind(&req); err != nil {
		logrus.WithError(err).Error("Failed to bind AddNewOrder request")
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": err.Error()})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	req.UserID = claims["user_id"].(string)

	order, err := api.validateAndGetOrder(&req)
	if err != nil {
		logrus.WithError(err).Error("Validation failed for AddNewOrder request")
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": err.Error()})
	}

	err = api.app.AddOrder(order)
	if err != nil {
		logrus.WithError(err).Error("Failed to add order")
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func validateGetProductsRequest(limitStr string, pageStr string) (int, int, error) {
	limit := 10
	page := 1
	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, err
		}
	}

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return 0, 0, err
		}
	}

	return limit, page, err
}

func (api *Api) validateAddProductRequest(input *request.ProductDetails) (*model.ProductDetails, error) {
	err := api.validator.Struct(input)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(input.Price, 64)
	if err != nil {
		return nil, errors.New("Invalid price ")
	}

	return &model.ProductDetails{
		Name:          input.Name,
		Price:         price,
		Category:      input.Category,
		AddedQuantity: input.AddedQuantity,
	}, nil
}

func (api *Api) Login(c echo.Context) error {
	req := new(request.UserLogin)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": err.Error()})
	}

	err := api.isValidLoginRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": err.Error()})
	}

	token, err := api.app.GenerateJWTToken(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}

	return c.JSON(http.StatusOK, &request.LoginResponse{
		Status: "success",
		Token:  token,
	})
}

func (api *Api) isValidLoginRequest(input *request.UserLogin) error {
	return api.validator.Struct(input)
}

func (api *Api) validateAndGetOrder(input *request.Order) (*model.Order, error) {
	err := api.validator.Struct(input)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(input.Price, 64)
	if err != nil {
		return nil, errors.New("Invalid price ")
	}

	return &model.Order{
		Quantity:  input.Quantity,
		Price:     price,
		ProductID: input.ProductID,
		UserID:    input.UserID,
	}, nil
}

func (api *Api) UpdateOrderStatus(c echo.Context) error {
	orderId := c.QueryParam("order_id")
	status := c.QueryParam("status")

	// validate the input
	orderDtl, err := validateUpdateOrderRequest(orderId, status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": err.Error()})
	}

	err = api.app.UpdateOrderStatus(orderId, orderDtl.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
	}

	return c.JSON(http.StatusOK, "success")
}

func validateUpdateOrderRequest(orderId string, status string) (*request.OrderDetail, error) {
	if orderId == "" {
		return nil, errors.New("Order Id is required ")
	}

	if status == "" {
		return nil, errors.New("Status is required ")
	}

	var orderStatus util.OrderStatus
	switch status {
	case string(util.OrderStatusConfirmed):
		orderStatus = util.OrderStatusConfirmed
	case string(util.OrderStatusCancelled):
		orderStatus = util.OrderStatusCancelled
	case string(util.OrderStatusShipped):
		orderStatus = util.OrderStatusShipped
	case string(util.OrderStatusDelivered):
		orderStatus = util.OrderStatusDelivered
	default:
		return nil, errors.New("Failed to update order status as status is invalid ")
	}

	return &request.OrderDetail{
		OrderID: orderId,
		Status:  orderStatus,
	}, nil
}
