package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"re/pkg/db"
	"sort"
	"strconv"
	"time"
)

type ProductController struct {
	persistence db.Persistence
}

type ProductRequest struct {
	PackSizes PackSizes `json:"pack_sizes"`
}

type PackSizes []int

func (p *ProductController) RegisterEndpoints(r *echo.Group) {
	r.GET("/packagecaclulator", p.PackageCacl)
	r.POST("/", p.NewProduct)
}

func (p *ProductController) PackageCacl(c echo.Context) error {
	ctx := c.Request().Context()
	items := c.QueryParam("items")
	itemsRequested, err := strconv.Atoi(items)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, db.Product{})
	}
	product, err := p.persistence.FetchProduct(ctx)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, db.Product{})
	}
	resp := calculatePacks(itemsRequested, product.PackSizes)
	return c.JSON(200, resp)
}

func (p *ProductController) NewProduct(c echo.Context) error {
	ctx := c.Request().Context()

	var prod ProductRequest
	if err := c.Bind(&prod); err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, db.Product{})
	}

	// Validate the pack sizes
	if len(prod.PackSizes) == 0 {
		return c.JSON(http.StatusBadRequest, "Pack sizes cannot be empty")
	}
	for _, size := range prod.PackSizes {
		if size <= 0 {
			return c.JSON(http.StatusBadRequest, "Pack sizes must be positive integers")
		}
	}

	newP := db.Product{
		Date:      time.Now(),
		Active:    1,
		PackSizes: db.PackSizes(prod.PackSizes),
	}
	err := p.persistence.InsertProduct(ctx, newP)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, db.Product{})
	}
	return c.String(http.StatusOK, "Product inserted")
}

func NewProductController(p db.Persistence) *ProductController {
	return &ProductController{
		persistence: p,
	}
}

func calculatePacks(orderQuantity int, availablePacks []int) map[int]int {
	// Sort the pack sizes in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(availablePacks)))

	packCounts := make(map[int]int)
	remainingQuantity := orderQuantity

	for i, packSize := range availablePacks {
		if remainingQuantity <= 0 {
			break
		}
		count := remainingQuantity / packSize
		remainder := remainingQuantity % packSize

		// If the current pack size can fit at least once
		if count > 0 {
			// Check if using the next smaller pack size results in using fewer items to fulfill the order
			if i < len(availablePacks)-1 && remainder != 0 {
				nextPackSize := availablePacks[i+1]
				if (packSize - remainder) < nextPackSize {
					continue // Skip to the next pack size if it results in a closer match to the order quantity
				}
			}

			packCounts[packSize] = count
			remainingQuantity -= packSize * count
		}
	}

	// If there's a remainder that doesn't exactly fit any pack size, add the smallest pack to fulfill the order
	if remainingQuantity > 0 {
		smallestPack := availablePacks[len(availablePacks)-1]
		packCounts[smallestPack]++
	}

	for size, count := range packCounts {
		for _, packSize := range availablePacks {
			if packSize == size*count {
				delete(packCounts, size)
				packCounts[packSize] = 1
			}
		}
	}

	return packCounts
}
