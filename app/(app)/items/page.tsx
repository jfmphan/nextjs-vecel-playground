"use client";

import { useCallback, useEffect, useState } from "react";
import {
  ActionIcon,
  Badge,
  Button,
  Center,
  Group,
  Loader,
  Select,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { IconEdit, IconPlus, IconSearch, IconTrash } from "@tabler/icons-react";
import { notifications } from "@mantine/notifications";
import { CategoryManager } from "@/components/CategoryManager";
import { ItemFormModal } from "@/components/ItemFormModal";
import { api, ApiError } from "@/lib/api";
import type { Container, Item, Named } from "@/lib/types";

export default function ItemsPage() {
  const [items, setItems] = useState<Item[]>([]);
  const [categories, setCategories] = useState<Named[]>([]);
  const [containers, setContainers] = useState<Container[]>([]);
  const [tags, setTags] = useState<Named[]>([]);
  const [loading, setLoading] = useState(true);

  const [search, setSearch] = useState("");
  const [categoryId, setCategoryId] = useState<string | null>(null);
  const [containerId, setContainerId] = useState<string | null>(null);

  const [opened, { open, close }] = useDisclosure(false);
  const [editing, setEditing] = useState<Item | null>(null);

  const loadItems = useCallback(async () => {
    setItems(
      await api.listItems({
        q: search || undefined,
        categoryId: categoryId ? Number(categoryId) : undefined,
        containerId: containerId ? Number(containerId) : undefined,
      }),
    );
  }, [search, categoryId, containerId]);

  const loadLookups = useCallback(async () => {
    const [c, ct, tg] = await Promise.all([
      api.listCategories(),
      api.listContainers(),
      api.listTags(),
    ]);
    setCategories(c);
    setContainers(ct);
    setTags(tg);
  }, []);

  // Initial load.
  useEffect(() => {
    Promise.all([loadItems(), loadLookups()]).finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Re-fetch items whenever the filters change (lightly debounced for typing).
  useEffect(() => {
    const handle = setTimeout(loadItems, 200);
    return () => clearTimeout(handle);
  }, [loadItems]);

  const categoryName = (id: number | null) =>
    categories.find((c) => c.id === id)?.name ?? "—";
  const containerName = (id: number | null) =>
    containers.find((c) => c.id === id)?.name ?? "—";

  function openAdd() {
    setEditing(null);
    open();
  }
  function openEdit(item: Item) {
    setEditing(item);
    open();
  }

  async function remove(item: Item) {
    if (!confirm(`Delete "${item.name}"?`)) return;
    try {
      await api.deleteItem(item.id);
      notifications.show({ message: "Item deleted", color: "green" });
      loadItems();
    } catch (err) {
      notifications.show({
        message: err instanceof ApiError ? err.message : "Delete failed",
        color: "red",
      });
    }
  }

  if (loading) {
    return (
      <Center mih={200}>
        <Loader />
      </Center>
    );
  }

  return (
    <Stack gap="md">
      <Group justify="space-between">
        <Title order={2}>Items</Title>
        <Group>
          <CategoryManager categories={categories} onChanged={loadLookups} />
          <Button leftSection={<IconPlus size={16} />} onClick={openAdd}>
            Add item
          </Button>
        </Group>
      </Group>

      <Group>
        <TextInput
          placeholder="Search…"
          leftSection={<IconSearch size={16} />}
          value={search}
          onChange={(e) => setSearch(e.currentTarget.value)}
          style={{ flex: 1, minWidth: 200 }}
        />
        <Select
          placeholder="Category"
          data={categories.map((c) => ({ value: String(c.id), label: c.name }))}
          value={categoryId}
          onChange={setCategoryId}
          clearable
        />
        <Select
          placeholder="Location"
          data={containers.map((c) => ({ value: String(c.id), label: c.name }))}
          value={containerId}
          onChange={setContainerId}
          clearable
        />
      </Group>

      {items.length === 0 ? (
        <Text c="dimmed">No items match. Add one to get started.</Text>
      ) : (
        <Table.ScrollContainer minWidth={760}>
          <Table highlightOnHover verticalSpacing="sm">
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Name</Table.Th>
                <Table.Th>Quantity</Table.Th>
                <Table.Th>Category</Table.Th>
                <Table.Th>Location</Table.Th>
                <Table.Th>Tags</Table.Th>
                <Table.Th>Expiry</Table.Th>
                <Table.Th />
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {items.map((item) => (
                <Table.Tr key={item.id}>
                  <Table.Td>{item.name}</Table.Td>
                  <Table.Td>
                    <Group gap="xs" wrap="nowrap">
                      <span>
                        {item.quantity}
                        {item.unit ? ` ${item.unit}` : ""}
                      </span>
                      {item.lowStock && (
                        <Badge color="orange" size="sm">
                          Low
                        </Badge>
                      )}
                    </Group>
                  </Table.Td>
                  <Table.Td>{categoryName(item.categoryId)}</Table.Td>
                  <Table.Td>{containerName(item.containerId)}</Table.Td>
                  <Table.Td>
                    <Group gap={4}>
                      {item.tags.map((t) => (
                        <Badge key={t} variant="light" size="sm">
                          {t}
                        </Badge>
                      ))}
                    </Group>
                  </Table.Td>
                  <Table.Td>{item.expiryDate ?? "—"}</Table.Td>
                  <Table.Td>
                    <Group gap="xs" justify="flex-end" wrap="nowrap">
                      <ActionIcon variant="subtle" onClick={() => openEdit(item)}>
                        <IconEdit size={16} />
                      </ActionIcon>
                      <ActionIcon variant="subtle" color="red" onClick={() => remove(item)}>
                        <IconTrash size={16} />
                      </ActionIcon>
                    </Group>
                  </Table.Td>
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>
        </Table.ScrollContainer>
      )}

      <ItemFormModal
        opened={opened}
        onClose={close}
        onSaved={() => {
          loadItems();
          loadLookups();
        }}
        item={editing}
        categories={categories}
        containers={containers}
        tagSuggestions={tags.map((t) => t.name)}
      />
    </Stack>
  );
}
