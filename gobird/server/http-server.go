package server

import (
	"fmt"
	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gobird/cmd"
	"gobird/route"
)

var gobird *cmd.GoBird


type CustomValidator struct {
}

func (c *CustomValidator) Validate(i interface{}) error {
	fmt.Printf("validating ==> %t %v\n", i, i)

	v := validate.Struct(i)

	if v.Validate() {
		fmt.Println("No Error")
		return nil
	} else {
		fmt.Printf("Errors => %v\n", v.Errors)
		return v.Errors
	}
}

func StartHTTPD(bird *cmd.GoBird) {
	gobird = bird

	e := echo.New()
	e.Debug = bird.Debug
	e.Validator = &CustomValidator{}
	e.Pre(middleware.RemoveTrailingSlash())

	/*
	  로그 포맷 참고
	 https://echo.labstack.com/middleware/logger
	 */
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} - ${uri} ${status}\n",
	}))

	e.Use(middleware.Recover())
	e.Use(useGoogle)

	e.GET("/drive", func(c echo.Context) error {
		bird := c.(*cmd.GoBird)

		results, err := bird.DriveServ.Drives.List().Fields("kind, nextPageToken, drives").Do()

		if err != nil {
			return c.JSON(500, err)
		}

		return c.JSON(200, results)

	})

	e.GET("/meta", func(c echo.Context) error {
		bird := c.(*cmd.GoBird)

		return c.JSON(200, map[string]interface{}{
			"Version":             cmd.Version,
			"ServiceAccountEmail": bird.ServiceAccountEmail,
		})
	})

	// API 정의
	resourceApi := e.Group("/api/resources")
	resourceApi.GET("", route.GetResources)
	resourceApi.GET("/:id", route.GetResource)
	resourceApi.POST("", route.CreateNewResource)
	resourceApi.DELETE("/:id", route.DeleteResource)


	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", bird.Port)))
}

//custom middleware
//https://echo.labstack.com/guide/context
func useGoogle(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		bird := &cmd.GoBird{}
		bird.Context = c
		bird.ServiceAccountEmail = gobird.ServiceAccountEmail
		bird.SheetServ = gobird.SheetServ
		bird.DriveServ = gobird.DriveServ
		bird.Debug = gobird.Debug

		return next(bird)
	}
}
