package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var app = fiber.New()

func TestRoutingHelloWorld(t *testing.T) {
	app.Get("/", func (ctx * fiber.Ctx) error {
		return ctx.SendString("Hello World!");
	})

	request := httptest.NewRequest("GET", "/", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello World!", string(bytes))
}

func TestCtx(t *testing.T) {
	app.Get("/hello", func (ctx *fiber.Ctx) error {
		name := ctx.Query("name", "Guest")

		return ctx.SendString("Hello " + name)
	})
	
	// Send query param
	request := httptest.NewRequest("GET", "/hello?name=Arza", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Hello Arza", string(bytes))

	// Without send query param
	request = httptest.NewRequest("GET", "/hello", nil)
	response, err = app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, _ = io.ReadAll(response.Body)
	assert.Equal(t, "Hello Guest", string(bytes))
}

func TestRouteParameter(t *testing.T) {
	app.Get("/customers/:customerId/addresses/:addressId", func (ctx *fiber.Ctx) error {
		return ctx.SendString("Get address " + ctx.Params("addressId") + " from customer " + ctx.Params("customerId"))
	})

	request := httptest.NewRequest("GET", "/customers/1/addresses/2", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	assert.Equal(t, "Get address 2 from customer 1", string(bytes))
}

func TestFormRequest(t *testing.T) {
	app.Post("/hello", func(ctx *fiber.Ctx) error {
		name := ctx.FormValue("name")
		return ctx.SendString("Hello " + name)
	})

	body := strings.NewReader("name=Arza")
	request := httptest.NewRequest("POST", "/hello", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, _ := app.Test(request)

	assert.Equal(t, 200, response.StatusCode)
	bytes, _ := io.ReadAll(response.Body)

	assert.Equal(t, "Hello Arza", string(bytes))
}