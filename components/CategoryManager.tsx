"use client";

import { useState } from "react";
import {
  ActionIcon,
  Button,
  Group,
  Popover,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import { IconCategory, IconTrash } from "@tabler/icons-react";
import { notifications } from "@mantine/notifications";
import { api, ApiError } from "@/lib/api";
import type { Named } from "@/lib/types";

// CategoryManager is a small popover for creating and deleting the managed
// category vocabulary, kept out of the item form to keep that form focused.
export function CategoryManager({
  categories,
  onChanged,
}: {
  categories: Named[];
  onChanged: () => void;
}) {
  const [name, setName] = useState("");

  async function add() {
    if (!name.trim()) return;
    try {
      await api.createCategory(name.trim());
      setName("");
      onChanged();
    } catch (err) {
      notifications.show({
        message: err instanceof ApiError ? err.message : "Failed to add category",
        color: "red",
      });
    }
  }

  async function remove(id: number) {
    try {
      await api.deleteCategory(id);
      onChanged();
    } catch (err) {
      notifications.show({
        message: err instanceof ApiError ? err.message : "Failed to delete category",
        color: "red",
      });
    }
  }

  return (
    <Popover width={280} position="bottom-end" withArrow shadow="md">
      <Popover.Target>
        <Button variant="default" leftSection={<IconCategory size={16} />}>
          Categories
        </Button>
      </Popover.Target>
      <Popover.Dropdown>
        <Stack gap="xs">
          <Group gap="xs">
            <TextInput
              placeholder="New category"
              value={name}
              onChange={(e) => setName(e.currentTarget.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") add();
              }}
              style={{ flex: 1 }}
            />
            <Button onClick={add}>Add</Button>
          </Group>
          {categories.length === 0 ? (
            <Text c="dimmed" size="sm">
              No categories yet.
            </Text>
          ) : (
            categories.map((c) => (
              <Group key={c.id} justify="space-between">
                <Text size="sm">{c.name}</Text>
                <ActionIcon variant="subtle" color="red" onClick={() => remove(c.id)}>
                  <IconTrash size={16} />
                </ActionIcon>
              </Group>
            ))
          )}
        </Stack>
      </Popover.Dropdown>
    </Popover>
  );
}
