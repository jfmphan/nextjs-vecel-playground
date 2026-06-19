# Features

A living list of what the Home Inventory app does. Status legend:
✅ done · 🟡 partial · ⏳ planned.

## Core inventory ✅

- **Items** with: name, description, quantity, unit, category, tags,
  **purchase date**, expiry date, low-stock threshold, photo URL, and value.
- **Create / edit / delete** items.
- **Search** by name/description and **filter** by category, location, tag, or
  low-stock status.

## Categories & tags ✅

- **Categories** — a managed list you curate (add/delete); each item can have one.
- **Tags** — free-form labels created on the fly when you type them on an item;
  used for quick filtering.

## Locations & containers ✅

- Organise items by **nested containers**: rooms → shelves → boxes
  (e.g. Garage → Shelf B → Box 3).
- Assign each item to a location; browse locations as a tree with item counts.
- Cyclic nesting is rejected (a container can't be inside itself or a descendant).
- Deleting a location detaches its items and sub-locations rather than deleting them.

## Quantity, low-stock & expiry alerts ✅

- Track quantity per item with an optional **low-stock threshold**.
- Track an **expiry date** per item.
- **Dashboard** summarising totals, items low on stock, and items expiring within
  30 days.

## Access control ✅

- A **single shared password** gates the whole app.
- Successful login sets a signed, stateless session cookie (HMAC-SHA256); no
  session store needed. Reachable from anywhere over HTTPS once deployed.

## Photos 🟡

- The data model and UI carry a `photoUrl` per item, and it renders where present.
- ⏳ Upload UI + storage (Vercel Blob free tier) is the next step to fully wire
  end-to-end.

## Barcode / QR scanning ⏳

- Planned: in-browser barcode scanning (works on a phone camera) to add/find
  items, printable QR labels for containers, and optional product lookup
  (e.g. Open Food Facts) to pre-fill item details.
