package service

import (
	"context"
	"errors"
	"testing"

	"homeinventory/internal/domain"
)

func TestLocation_CycleDetection(t *testing.T) {
	svc := NewLocationService(newFakeContainerRepo())
	ctx := context.Background()

	a, err := svc.Create(ctx, ContainerInput{Name: "Garage", Type: domain.ContainerRoom})
	if err != nil {
		t.Fatalf("create A: %v", err)
	}
	b, err := svc.Create(ctx, ContainerInput{Name: "Shelf", Type: domain.ContainerShelf, ParentID: &a.ID})
	if err != nil {
		t.Fatalf("create B: %v", err)
	}

	// Re-parenting A under its own descendant B would form A -> B -> A.
	if _, err := svc.Update(ctx, a.ID, ContainerInput{
		Name: "Garage", Type: domain.ContainerRoom, ParentID: &b.ID,
	}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error for cycle, got %v", err)
	}

	// A container cannot be its own parent.
	if _, err := svc.Update(ctx, a.ID, ContainerInput{
		Name: "Garage", Type: domain.ContainerRoom, ParentID: &a.ID,
	}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error for self-parent, got %v", err)
	}
}

func TestLocation_DefaultsTypeAndValidatesName(t *testing.T) {
	svc := NewLocationService(newFakeContainerRepo())
	ctx := context.Background()

	if _, err := svc.Create(ctx, ContainerInput{Name: "  "}); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error for empty name, got %v", err)
	}

	c, err := svc.Create(ctx, ContainerInput{Name: "Misc"}) // no type → defaults to "other"
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if c.Type != domain.ContainerOther {
		t.Fatalf("expected default type %q, got %q", domain.ContainerOther, c.Type)
	}
}
