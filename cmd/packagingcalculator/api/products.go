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

type Order struct {
	size            int
	packs           map[int]int
	fulfillmentSize int
	totalPacks      int
}

func calculatePacks(orderSize int, packSizes []int) map[int]int {
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))
	packSizesCpy := make([]int, len(packSizes))
	copy(packSizesCpy, packSizes)
	fulfillments := make([]Order, 0)

	for len(packSizes) > 0 {
		o := Order{}
		packs := make(map[int]int)
		remaining := orderSize

		//allocate packs from the largest to the smallest
		for _, pack := range packSizes {
			if remaining == 0 {
				break
			}

			count := remaining / pack
			if count > 0 {
				packs[pack] = count
				o.totalPacks += count
				o.fulfillmentSize += count * pack
				remaining -= count * pack
			}
		}

		// Handle any remaining items that need to be rounded up to the nearest pack
		if remaining > 0 {
			smallestPack := packSizes[len(packSizes)-1]
			packs[smallestPack]++
			o.totalPacks++
			o.fulfillmentSize += smallestPack
		}

		o.size = orderSize
		o.packs = packs

		packSizes = packSizes[1:]
		fulfillments = append(fulfillments, o)
	}

	for _, f := range fulfillments {
		for size, count := range f.packs {
			for _, packSize := range packSizesCpy {
				if packSize == size*count {
					delete(f.packs, size)
					f.packs[packSize] = 1
				}
			}
		}
	}

	// second rule, order by fulfillment size
	sort.Slice(fulfillments, func(i, j int) bool {
		return fulfillments[i].fulfillmentSize < fulfillments[j].fulfillmentSize
	})

	// get the smallest fulfillments of the same size
	smallestFulfillments := make([]Order, 0)
	for _, f := range fulfillments {
		if f.fulfillmentSize == fulfillments[0].fulfillmentSize {
			smallestFulfillments = append(smallestFulfillments, f)
		} else {
			break
		}
	}

	//third rule, order by total packs and return the smallest
	sort.Slice(smallestFulfillments, func(i, j int) bool {
		return smallestFulfillments[i].totalPacks < smallestFulfillments[j].totalPacks
	})

	return smallestFulfillments[0].packs
}
