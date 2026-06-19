// Typed client for the Go API. All calls are relative to /api/v1 so they stay
// same-origin (the session cookie is sent automatically).

import type {
  Container,
  ContainerInput,
  Item,
  ItemInput,
  Named,
  Stats,
} from "./types";

const BASE = "/api/v1";

/** ApiError carries the HTTP status alongside the server's error message. */
export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  if (!res.ok) {
    throw new ApiError(res.status, await errorMessage(res));
  }
  if (res.status === 204) {
    return undefined as T;
  }
  return (await res.json()) as T;
}

async function errorMessage(res: Response): Promise<string> {
  try {
    const body = (await res.json()) as { error?: string };
    if (body.error) return body.error;
  } catch {
    // fall through to a generic message
  }
  return `request failed (${res.status})`;
}

/** ItemQuery mirrors the optional filters supported by GET /items. */
export interface ItemQuery {
  q?: string;
  categoryId?: number | null;
  containerId?: number | null;
  tag?: string;
  lowStock?: boolean;
}

function queryString(params: ItemQuery): string {
  const search = new URLSearchParams();
  if (params.q) search.set("q", params.q);
  if (params.categoryId != null) search.set("categoryId", String(params.categoryId));
  if (params.containerId != null) search.set("containerId", String(params.containerId));
  if (params.tag) search.set("tag", params.tag);
  if (params.lowStock) search.set("lowStock", "true");
  const s = search.toString();
  return s ? `?${s}` : "";
}

export const api = {
  // Auth
  session: () => request<{ authenticated: boolean }>("/session"),
  login: (password: string) =>
    request<{ ok: boolean }>("/login", {
      method: "POST",
      body: JSON.stringify({ password }),
    }),
  logout: () => request<{ ok: boolean }>("/logout", { method: "POST" }),

  // Items
  listItems: (params: ItemQuery = {}) => request<Item[]>(`/items${queryString(params)}`),
  getItem: (id: number) => request<Item>(`/items/${id}`),
  createItem: (input: ItemInput) =>
    request<Item>("/items", { method: "POST", body: JSON.stringify(input) }),
  updateItem: (id: number, input: ItemInput) =>
    request<Item>(`/items/${id}`, { method: "PUT", body: JSON.stringify(input) }),
  deleteItem: (id: number) => request<void>(`/items/${id}`, { method: "DELETE" }),

  // Containers
  listContainers: () => request<Container[]>("/containers"),
  createContainer: (input: ContainerInput) =>
    request<Container>("/containers", { method: "POST", body: JSON.stringify(input) }),
  updateContainer: (id: number, input: ContainerInput) =>
    request<Container>(`/containers/${id}`, { method: "PUT", body: JSON.stringify(input) }),
  deleteContainer: (id: number) =>
    request<void>(`/containers/${id}`, { method: "DELETE" }),

  // Categories & tags
  listCategories: () => request<Named[]>("/categories"),
  createCategory: (name: string) =>
    request<Named>("/categories", { method: "POST", body: JSON.stringify({ name }) }),
  deleteCategory: (id: number) =>
    request<void>(`/categories/${id}`, { method: "DELETE" }),
  listTags: () => request<Named[]>("/tags"),

  // Dashboard
  stats: () => request<Stats>("/stats"),
};
