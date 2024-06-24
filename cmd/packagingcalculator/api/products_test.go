package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"re/pkg/db"
	"strings"
	"testing"
	"time"
)

func TestInsertProduct(t *testing.T) {
	type fields struct {
		TestName string
		Product  db.Product
		Request  string
		Status   int
		Response string
		Error    error
	}
	tests := []fields{
		{
			TestName: "test insert product ok",
			Product: db.Product{
				Date:   time.Now(),
				Active: 1,
				PackSizes: db.PackSizes{
					250, 500, 1000, 2000, 5000,
				},
			},
			Request:  `{"pack_sizes":[250,500,1000,2000,5000]}`,
			Status:   http.StatusOK,
			Response: "Product inserted",
		},
		{
			TestName: "test insert product error",
			Product: db.Product{
				Date:   time.Now(),
				Active: 1,
				PackSizes: db.PackSizes{
					250, 500, 1000, 2000, 5000,
				},
			},
			Request: `{"pack_sizes":[250,500,1000,2000,5000]}`,
			Status:  http.StatusInternalServerError,
			Error:   errors.New("db error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/product", strings.NewReader(tt.Request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			d := db.NewMockPersistence(t)
			d.On("InsertProduct", c.Request().Context(), mock.Anything).Return(tt.Error)
			pc := NewProductController(d)
			BuildRoutes(pc)
			assert.NoError(t, pc.NewProduct(c))

			require.Equal(t, tt.Status, rec.Code)
			if tt.Error != nil {
				//assert status 400
				assert.Equal(t, tt.Status, rec.Code)
			} else {
				assert.Equal(
					t,
					tt.Response,
					rec.Body.String(),
				)
			}
		})
	}
}

func TestCalculatePackages(t *testing.T) {
	type fields struct {
		TestName string
		Items    int
		Status   int
		Response string
		Error    error
	}
	tests := []fields{
		{
			TestName: "test calculate packages ok",
			Items:    251,
			Status:   http.StatusOK,
			Response: `{"500":1}`,
		},
		{
			TestName: "test calculate packages error",
			Items:    251,
			Status:   http.StatusInternalServerError,
			Error:    errors.New("db error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			e := echo.New()
			rec := httptest.NewRecorder()
			q := make(url.Values)
			q.Set("items", "251")

			req := httptest.NewRequest(http.MethodGet, "/v1/packagecaclulator?"+q.Encode(), nil)
			c := e.NewContext(req, rec)
			d := db.NewMockPersistence(t)
			d.On("FetchProduct", c.Request().Context()).Return(db.Product{
				Date:   time.Now(),
				Active: 1,
				PackSizes: db.PackSizes{
					250, 500, 1000, 2000, 5000,
				},
			}, tt.Error)
			pc := NewProductController(d)
			BuildRoutes(pc)
			assert.NoError(t, pc.PackageCacl(c))

			require.Equal(t, tt.Status, rec.Code)
			if tt.Error != nil {
				//assert status 400
				assert.Equal(t, tt.Status, rec.Code)
			} else {
				assert.JSONEq(
					t,
					tt.Response,
					rec.Body.String(),
				)
			}
		})
	}
}

func TestCalculatePacks(t *testing.T) {
	type args struct {
		orderQuantity  int
		availablePacks []int
	}
	tests := []struct {
		name string
		args args
		want map[int]int
	}{
		{
			name: "test calculate packs",
			args: args{
				orderQuantity:  1,
				availablePacks: []int{250, 500, 1000, 2000, 5000},
			},
			want: map[int]int{
				250: 1,
			},
		},
		{
			name: "test calculate packs 2",
			args: args{
				orderQuantity:  250,
				availablePacks: []int{250, 500, 1000, 2000, 5000},
			},
			want: map[int]int{
				250: 1,
			},
		},
		{
			name: "test calculate packs 3",
			args: args{
				orderQuantity:  251,
				availablePacks: []int{250, 500, 1000, 2000, 5000},
			},
			want: map[int]int{
				500: 1,
			},
		},
		{
			name: "test calculate packs 4",
			args: args{
				orderQuantity:  501,
				availablePacks: []int{250, 500, 1000, 2000, 5000},
			},
			want: map[int]int{
				500: 1,
				250: 1,
			},
		},
		{
			name: "test calculate packs 5",
			args: args{
				orderQuantity:  12001,
				availablePacks: []int{250, 500, 1000, 2000, 5000},
			},
			want: map[int]int{
				5000: 2,
				2000: 1,
				250:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculatePacks(tt.args.orderQuantity, tt.args.availablePacks)
			assert.Equal(t, tt.want, got)
		})
	}
}
