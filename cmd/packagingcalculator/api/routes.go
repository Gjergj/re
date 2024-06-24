package api

import (
	"github.com/labstack/echo/v4"
)

func BuildRoutes(rc *ProductController) *echo.Echo {
	e := echo.New()
	{
		// v1
		v1 := e.Group("/v1")
		{
			// products endpoints
			{
				productRoute := v1.Group("/products")
				rc.RegisterEndpoints(productRoute)
			}
		}
	}
	e.Static("/", "assets")
	return e
}
