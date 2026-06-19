"use client";

import { useEffect } from "react";
import {
  Button,
  Group,
  Modal,
  NumberInput,
  Select,
  Stack,
  TagsInput,
  Textarea,
  TextInput,
} from "@mantine/core";
import { DateInput } from "@mantine/dates";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { api, ApiError } from "@/lib/api";
import type { Container, Item, ItemInput, Named } from "@/lib/types";

interface Props {
  opened: boolean;
  onClose: () => void;
  onSaved: () => void;
  item: Item | null; // null = create
  categories: Named[];
  containers: Container[];
  tagSuggestions: string[];
}

// FormValues uses string-or-number for numeric inputs (Mantine NumberInput) and
// string ids for Selects; toInput converts them to the API's typed shape.
interface FormValues {
  name: string;
  description: string;
  quantity: number | string;
  unit: string;
  categoryId: string | null;
  containerId: string | null;
  lowStockThreshold: number | string;
  purchaseDate: string | null;
  expiryDate: string | null;
  tags: string[];
}

const EMPTY: FormValues = {
  name: "",
  description: "",
  quantity: 1,
  unit: "",
  categoryId: null,
  containerId: null,
  lowStockThreshold: "",
  purchaseDate: null,
  expiryDate: null,
  tags: [],
};

export function ItemFormModal({
  opened,
  onClose,
  onSaved,
  item,
  categories,
  containers,
  tagSuggestions,
}: Props) {
  const form = useForm<FormValues>({
    initialValues: EMPTY,
    validate: { name: (v) => (v.trim() ? null : "Name is required") },
  });

  // Load the selected item (or a blank form) each time the modal opens.
  useEffect(() => {
    if (opened) form.setValues(item ? toForm(item) : EMPTY);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [opened, item]);

  async function handleSubmit(values: FormValues) {
    try {
      const input = toInput(values);
      if (item) await api.updateItem(item.id, input);
      else await api.createItem(input);
      notifications.show({ message: item ? "Item updated" : "Item added", color: "green" });
      onSaved();
      onClose();
    } catch (err) {
      notifications.show({
        message: err instanceof ApiError ? err.message : "Save failed",
        color: "red",
      });
    }
  }

  return (
    <Modal opened={opened} onClose={onClose} title={item ? "Edit item" : "Add item"} size="lg">
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <TextInput label="Name" withAsterisk {...form.getInputProps("name")} />
          <Textarea label="Description" autosize minRows={2} {...form.getInputProps("description")} />

          <Group grow>
            <NumberInput label="Quantity" min={0} {...form.getInputProps("quantity")} />
            <TextInput label="Unit" placeholder="ea, box, kg…" {...form.getInputProps("unit")} />
          </Group>

          <Group grow>
            <Select
              label="Category"
              placeholder="None"
              data={categories.map((c) => ({ value: String(c.id), label: c.name }))}
              clearable
              searchable
              {...form.getInputProps("categoryId")}
            />
            <Select
              label="Location"
              placeholder="None"
              data={containers.map((c) => ({ value: String(c.id), label: c.name }))}
              clearable
              searchable
              {...form.getInputProps("containerId")}
            />
          </Group>

          <Group grow>
            <DateInput
              label="Purchase date"
              clearable
              valueFormat="YYYY-MM-DD"
              {...form.getInputProps("purchaseDate")}
            />
            <DateInput
              label="Expiry date"
              clearable
              valueFormat="YYYY-MM-DD"
              {...form.getInputProps("expiryDate")}
            />
          </Group>

          <NumberInput
            label="Low-stock threshold"
            description="Warn when quantity falls to this or below"
            min={0}
            {...form.getInputProps("lowStockThreshold")}
          />

          <TagsInput label="Tags" data={tagSuggestions} {...form.getInputProps("tags")} />

          <Group justify="flex-end" mt="sm">
            <Button variant="default" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit">{item ? "Save" : "Add"}</Button>
          </Group>
        </Stack>
      </form>
    </Modal>
  );
}

function toForm(item: Item): FormValues {
  return {
    name: item.name,
    description: item.description,
    quantity: item.quantity,
    unit: item.unit,
    categoryId: item.categoryId != null ? String(item.categoryId) : null,
    containerId: item.containerId != null ? String(item.containerId) : null,
    lowStockThreshold: item.lowStockThreshold ?? "",
    purchaseDate: item.purchaseDate,
    expiryDate: item.expiryDate,
    tags: item.tags,
  };
}

function toInput(values: FormValues): ItemInput {
  return {
    name: values.name.trim(),
    description: values.description,
    quantity: Number(values.quantity) || 0,
    unit: values.unit,
    categoryId: values.categoryId ? Number(values.categoryId) : null,
    containerId: values.containerId ? Number(values.containerId) : null,
    lowStockThreshold:
      values.lowStockThreshold === "" ? null : Number(values.lowStockThreshold),
    purchaseDate: values.purchaseDate,
    expiryDate: values.expiryDate,
    photoUrl: "",
    valueCents: null,
    tags: values.tags,
  };
}
