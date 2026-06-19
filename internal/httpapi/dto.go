package httpapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"homeinventory/internal/domain"
	"homeinventory/internal/service"
)

// --- Items -----------------------------------------------------------------

type itemRequest struct {
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	CategoryID        *int64       `json:"categoryId"`
	ContainerID       *int64       `json:"containerId"`
	Quantity          int          `json:"quantity"`
	Unit              string       `json:"unit"`
	LowStockThreshold *int         `json:"lowStockThreshold"`
	PurchaseDate      *domain.Date `json:"purchaseDate"`
	ExpiryDate        *domain.Date `json:"expiryDate"`
	PhotoURL          string       `json:"photoUrl"`
	ValueCents        *int64       `json:"valueCents"`
	Tags              []string     `json:"tags"`
}

func (r itemRequest) toInput() service.ItemInput {
	return service.ItemInput{
		Name:              r.Name,
		Description:       r.Description,
		CategoryID:        r.CategoryID,
		ContainerID:       r.ContainerID,
		Quantity:          r.Quantity,
		Unit:              r.Unit,
		LowStockThreshold: r.LowStockThreshold,
		PurchaseDate:      r.PurchaseDate,
		ExpiryDate:        r.ExpiryDate,
		PhotoURL:          r.PhotoURL,
		ValueCents:        r.ValueCents,
		Tags:              r.Tags,
	}
}

type itemResponse struct {
	ID                int64        `json:"id"`
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	CategoryID        *int64       `json:"categoryId"`
	ContainerID       *int64       `json:"containerId"`
	Quantity          int          `json:"quantity"`
	Unit              string       `json:"unit"`
	LowStockThreshold *int         `json:"lowStockThreshold"`
	PurchaseDate      *domain.Date `json:"purchaseDate"`
	ExpiryDate        *domain.Date `json:"expiryDate"`
	PhotoURL          string       `json:"photoUrl"`
	ValueCents        *int64       `json:"valueCents"`
	Tags              []string     `json:"tags"`
	LowStock          bool         `json:"lowStock"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
}

func newItemResponse(it domain.Item) itemResponse {
	return itemResponse{
		ID:                it.ID,
		Name:              it.Name,
		Description:       it.Description,
		CategoryID:        it.CategoryID,
		ContainerID:       it.ContainerID,
		Quantity:          it.Quantity,
		Unit:              it.Unit,
		LowStockThreshold: it.LowStockThreshold,
		PurchaseDate:      it.PurchaseDate,
		ExpiryDate:        it.ExpiryDate,
		PhotoURL:          it.PhotoURL,
		ValueCents:        it.ValueCents,
		Tags:              orEmpty(it.Tags),
		LowStock:          it.IsLowStock(),
		CreatedAt:         it.CreatedAt,
		UpdatedAt:         it.UpdatedAt,
	}
}

func newItemResponses(items []domain.Item) []itemResponse {
	out := make([]itemResponse, 0, len(items))
	for _, it := range items {
		out = append(out, newItemResponse(it))
	}
	return out
}

// --- Containers ------------------------------------------------------------

type containerRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	ParentID *int64 `json:"parentId"`
}

func (r containerRequest) toInput() service.ContainerInput {
	return service.ContainerInput{
		Name:     r.Name,
		Type:     domain.ContainerType(r.Type),
		ParentID: r.ParentID,
	}
}

type containerResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	ParentID  *int64    `json:"parentId"`
	ItemCount int       `json:"itemCount"`
	CreatedAt time.Time `json:"createdAt"`
}

func newContainerResponse(c domain.Container, itemCount int) containerResponse {
	return containerResponse{
		ID:        c.ID,
		Name:      c.Name,
		Type:      string(c.Type),
		ParentID:  c.ParentID,
		ItemCount: itemCount,
		CreatedAt: c.CreatedAt,
	}
}

// --- Categories & tags -----------------------------------------------------

type categoryRequest struct {
	Name string `json:"name"`
}

type namedResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func newCategoryResponses(categories []domain.Category) []namedResponse {
	out := make([]namedResponse, 0, len(categories))
	for _, c := range categories {
		out = append(out, namedResponse{ID: c.ID, Name: c.Name})
	}
	return out
}

func newTagResponses(tags []domain.Tag) []namedResponse {
	out := make([]namedResponse, 0, len(tags))
	for _, t := range tags {
		out = append(out, namedResponse{ID: t.ID, Name: t.Name})
	}
	return out
}

// --- Stats -----------------------------------------------------------------

type statsResponse struct {
	TotalItems    int            `json:"totalItems"`
	TotalQuantity int            `json:"totalQuantity"`
	LowStockCount int            `json:"lowStockCount"`
	ExpiringCount int            `json:"expiringCount"`
	LowStock      []itemResponse `json:"lowStock"`
	Expiring      []itemResponse `json:"expiring"`
}

func newStatsResponse(s service.Stats) statsResponse {
	return statsResponse{
		TotalItems:    s.TotalItems,
		TotalQuantity: s.TotalQuantity,
		LowStockCount: len(s.LowStock),
		ExpiringCount: len(s.Expiring),
		LowStock:      newItemResponses(s.LowStock),
		Expiring:      newItemResponses(s.Expiring),
	}
}

// --- Request parsing helpers ----------------------------------------------

func orEmpty(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

// idParam reads the {id} path parameter as an int64.
func idParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return 0, domain.Invalid("invalid id")
	}
	return id, nil
}

func optionalInt64Query(r *http.Request, key string) (*int64, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, domain.Invalid("invalid " + key)
	}
	return &n, nil
}

func optionalDateQuery(r *http.Request, key string) (*domain.Date, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil, nil
	}
	d, err := domain.ParseDate(v)
	if err != nil {
		return nil, domain.Invalid("invalid " + key)
	}
	return &d, nil
}

// parseItemFilter builds an ItemFilter from the request query string.
func parseItemFilter(r *http.Request) (domain.ItemFilter, error) {
	q := r.URL.Query()
	filter := domain.ItemFilter{
		Query:        q.Get("q"),
		Tag:          q.Get("tag"),
		LowStockOnly: q.Get("lowStock") == "true",
	}
	var err error
	if filter.CategoryID, err = optionalInt64Query(r, "categoryId"); err != nil {
		return filter, err
	}
	if filter.ContainerID, err = optionalInt64Query(r, "containerId"); err != nil {
		return filter, err
	}
	if filter.ExpiringBefore, err = optionalDateQuery(r, "expiringBefore"); err != nil {
		return filter, err
	}
	return filter, nil
}
