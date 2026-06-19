// Shared types mirroring the JSON shapes returned by the Go API (see
// internal/httpapi/dto.go). Dates are ISO "YYYY-MM-DD" strings; timestamps are
// RFC3339 strings.

export type ContainerType = "room" | "shelf" | "box" | "other";

export interface Item {
  id: number;
  name: string;
  description: string;
  categoryId: number | null;
  containerId: number | null;
  quantity: number;
  unit: string;
  lowStockThreshold: number | null;
  purchaseDate: string | null;
  expiryDate: string | null;
  photoUrl: string;
  valueCents: number | null;
  tags: string[];
  lowStock: boolean;
  createdAt: string;
  updatedAt: string;
}

// ItemInput is the request body for create/update. Server-managed fields (id,
// timestamps, lowStock) are omitted.
export interface ItemInput {
  name: string;
  description: string;
  categoryId: number | null;
  containerId: number | null;
  quantity: number;
  unit: string;
  lowStockThreshold: number | null;
  purchaseDate: string | null;
  expiryDate: string | null;
  photoUrl: string;
  valueCents: number | null;
  tags: string[];
}

export interface Container {
  id: number;
  name: string;
  type: ContainerType;
  parentId: number | null;
  itemCount: number;
  createdAt: string;
}

export interface ContainerInput {
  name: string;
  type: ContainerType;
  parentId: number | null;
}

// Named is the shape shared by categories and tags.
export interface Named {
  id: number;
  name: string;
}

export interface Stats {
  totalItems: number;
  totalQuantity: number;
  lowStockCount: number;
  expiringCount: number;
  lowStock: Item[];
  expiring: Item[];
}
