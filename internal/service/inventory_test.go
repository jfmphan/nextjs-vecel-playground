package service

import (
	"context"
	"errors"
	"testing"

	"homeinventory/internal/domain"
)

func TestInventoryCreate_Validation(t *testing.T) {
	svc := NewInventoryService(
		newFakeItemRepo(),
		&fakeCategoryRepo{existing: map[int64]bool{}},
		newFakeContainerRepo(),
	)

	cases := map[string]ItemInput{
		"empty name":        {Name: "  ", Quantity: 1},
		"negative quantity": {Name: "Thing", Quantity: -1},
		"unknown category":  {Name: "Thing", Quantity: 1, CategoryID: ptr(int64(99))},
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := svc.Create(context.Background(), in); !errors.Is(err, domain.ErrValidation) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	}
}

func TestInventoryCreate_OK(t *testing.T) {
	svc := NewInventoryService(
		newFakeItemRepo(),
		&fakeCategoryRepo{existing: map[int64]bool{1: true}},
		newFakeContainerRepo(),
	)

	item, err := svc.Create(context.Background(), ItemInput{
		Name:       "Cordless Drill",
		Quantity:   2,
		CategoryID: ptr(int64(1)),
		Tags:       []string{"power"},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if item.ID == 0 || item.Name != "Cordless Drill" {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestInventoryCreate_TrimsName(t *testing.T) {
	svc := NewInventoryService(
		newFakeItemRepo(),
		&fakeCategoryRepo{existing: map[int64]bool{}},
		newFakeContainerRepo(),
	)
	item, err := svc.Create(context.Background(), ItemInput{Name: "  Hammer  ", Quantity: 1})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if item.Name != "Hammer" {
		t.Fatalf("expected trimmed name, got %q", item.Name)
	}
}
