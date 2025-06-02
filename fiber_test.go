package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var app = fiber.New(fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.SendString("Terjadi Kesalahan: " + err.Error())
	},
})

func TestRoutingHelloWorld(t *testing.T) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World!")
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
	app.Get("/hello", func(ctx *fiber.Ctx) error {
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
	app.Get("/customers/:customerId/addresses/:addressId", func(ctx *fiber.Ctx) error {
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

//go:embed source/contoh.txt
var contohFile []byte

func TestFormUpload(t *testing.T) {
	app.Post("/upload", func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("file")

		if err != nil {
			return err
		}

		// Store file
		err = ctx.SaveFile(file, "./target/"+file.Filename)
		if err != nil {
			return err
		}

		return ctx.SendString("success")
	})

	// Creating file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("file", "contoh.txt")
	file.Write(contohFile)
	writer.Close()

	// Post
	request := httptest.NewRequest("POST", "/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "success", string(bytes))

	assert.FileExists(t, "./target/contoh.txt")
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestRequestBody(t *testing.T) {
	app.Post("/login", func(ctx *fiber.Ctx) error {
		body := ctx.Body()

		request := new(LoginRequest)
		err := json.Unmarshal(body, &request)

		if err != nil {
			return err
		}

		return ctx.SendString("Hello " + request.Username)
	})

	body := strings.NewReader(`{"username":"Arza"}`)
	request := httptest.NewRequest("POST", "/login", body)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Arza", string(bytes))
}

func TestResponseJson(t *testing.T) {
	app.Get("/users/:userId", func (ctx *fiber.Ctx) error {
		userId, _ := strconv.Atoi(ctx.Params("userId"))

		return ctx.JSON(fiber.Map{
			"id": userId,
			"name": "Arza",
			"email": "zaarza03@gmail.com",
		})
	})

	// Send request
	request := httptest.NewRequest("GET", "/users/10", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	// method a
	// assert.Equal(t, `{"email":"zaarza03@gmail.com","id":10,"name":"Arza"}`, string(bytes))

	// method b
	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)

	assert.Equal(t, float64(10), data["id"])
	assert.Equal(t, "Arza", data["name"])
	assert.Equal(t, "zaarza03@gmail.com", data["email"])
}

func TestRoutingGroup(t *testing.T) {
	helloWorld := func(ctx *fiber.Ctx) error {
		return ctx.SendString("Routing Group!");
	}

	api := app.Group("/api")
	api.Get("/hello1", helloWorld)

	web := app.Group("/web")
	web.Get("/hello2", helloWorld)

	// Test /hello1
	request1 := httptest.NewRequest("GET", "/api/hello1", nil)
	response1, _ := app.Test(request1)
	
	bytes1, _ := io.ReadAll(response1.Body)

	assert.Equal(t, "Routing Group!", string(bytes1))

	// Test /hello1
	request2 := httptest.NewRequest("GET", "/web/hello2", nil)
	response2, _ := app.Test(request2)
	
	bytes2, _ := io.ReadAll(response2.Body)

	assert.Equal(t, "Routing Group!", string(bytes2))
}

func TestErrorHandling(t *testing.T) {
	app.Get("/error", func(ctx *fiber.Ctx) error {
		return errors.New("Ups")
	})

	request := httptest.NewRequest("GET", "/error", nil)
	response, err := app.Test(request)

	if err != nil {
		panic(err)
	}

	bytes, _ := io.ReadAll(response.Body)
	assert.Equal(t, 500, response.StatusCode)
	assert.Equal(t, "Terjadi Kesalahan: Ups", string(bytes))
}