package main

import (
	"OnlieStore/internal/api"
	"OnlieStore/internal/app"
	"github.com/labstack/echo/v4"
)

func main() {
	newApp := app.NewApp() // new app

	e := echo.New()
	newApi := api.NewApi(newApp, e) // new api

	err := newApp.LoadData()
	if err != nil {
		panic(err)
	}

	newApi.RegisterFunctions()
	newApi.StartService()
}
