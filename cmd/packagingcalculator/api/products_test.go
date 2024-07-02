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
		{
			name: "test calculate packs 6",
			args: args{
				orderQuantity:  14,
				availablePacks: []int{3, 7, 13},
			},
			want: map[int]int{
				7: 2,
			},
		},
		{
			name: "test calculate packs 7",
			args: args{
				orderQuantity:  5,
				availablePacks: []int{4, 8, 13},
			},
			want: map[int]int{
				8: 1,
			},
		},
		{
			name: "test calculate packs 8",
			args: args{
				orderQuantity:  32,
				availablePacks: []int{4, 7, 9, 25},
			},
			want: map[int]int{
				7:  1,
				25: 1,
			},
		},
		{
			name: "test calculate packs 9",
			args: args{
				orderQuantity:  17,
				availablePacks: []int{3, 7, 13},
			},
			want: map[int]int{
				3: 1,
				7: 2,
			},
		},
		{
			name: "test calculate packs 10",
			args: args{
				orderQuantity:  14,
				availablePacks: []int{3, 7, 13, 14},
			},
			want: map[int]int{
				14: 1,
			},
		},
		{
			name: "test calculate packs 11",
			args: args{
				orderQuantity:  8,
				availablePacks: []int{3, 7, 13, 14},
			},
			want: map[int]int{
				3: 3,
			},
		},
		{
			name: "test calculate packs 12",
			args: args{
				orderQuantity:  5,
				availablePacks: []int{3, 7, 13, 14},
			},
			want: map[int]int{
				3: 2,
			},
		},
		{
			name: "test calculate packs 13",
			args: args{
				orderQuantity:  11,
				availablePacks: []int{2, 3, 4, 7, 14},
			},
			want: map[int]int{
				4: 1,
				7: 1,
			},
		},
		{
			name: "test calculate packs 14",
			args: args{
				orderQuantity:  13,
				availablePacks: []int{3, 7, 13, 14},
			},
			want: map[int]int{
				13: 1,
			},
		},
		{
			name: "test calculate packs 15",
			args: args{
				orderQuantity:  12,
				availablePacks: []int{3, 7, 13, 14},
			},
			want: map[int]int{
				3: 4,
			},
		},
		{
			name: "test calculate packs 16",
			args: args{
				orderQuantity:  20,
				availablePacks: []int{3, 7, 13, 14},
			},
			want: map[int]int{
				7:  1,
				13: 1,
			},
		},
		{
			name: "test calculate packs 17",
			args: args{
				orderQuantity:  35,
				availablePacks: []int{3, 7, 13, 14},
			},
			want: map[int]int{
				7:  1,
				14: 2,
			},
		},

		{
			name: "test calculate packs 18",
			args: args{
				orderQuantity:  37,
				availablePacks: []int{3, 7, 13, 14},
			},
			want: map[int]int{
				3:  1,
				7:  1,
				14: 2,
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
