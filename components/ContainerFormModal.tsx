"use client";

import { useEffect } from "react";
import { Button, Group, Modal, Select, Stack, TextInput } from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { api, ApiError } from "@/lib/api";
import type { Container, ContainerInput, ContainerType } from "@/lib/types";

const TYPES: { value: ContainerType; label: string }[] = [
  { value: "room", label: "Room" },
  { value: "shelf", label: "Shelf" },
  { value: "box", label: "Box" },
  { value: "other", label: "Other" },
];

interface Props {
  opened: boolean;
  onClose: () => void;
  onSaved: () => void;
  container: Container | null; // null = create
  containers: Container[];
}

interface FormValues {
  name: string;
  type: ContainerType;
  parentId: string | null;
}

export function ContainerFormModal({
  opened,
  onClose,
  onSaved,
  container,
  containers,
}: Props) {
  const form = useForm<FormValues>({
    initialValues: { name: "", type: "box", parentId: null },
    validate: { name: (v) => (v.trim() ? null : "Name is required") },
  });

  useEffect(() => {
    if (!opened) return;
    form.setValues(
      container
        ? {
            name: container.name,
            type: container.type,
            parentId: container.parentId != null ? String(container.parentId) : null,
          }
        : { name: "", type: "box", parentId: null },
    );
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [opened, container]);

  // A container cannot be nested under itself (deeper cycles are rejected by the API).
  const parentOptions = containers
    .filter((c) => c.id !== container?.id)
    .map((c) => ({ value: String(c.id), label: c.name }));

  async function handleSubmit(values: FormValues) {
    const input: ContainerInput = {
      name: values.name.trim(),
      type: values.type,
      parentId: values.parentId ? Number(values.parentId) : null,
    };
    try {
      if (container) await api.updateContainer(container.id, input);
      else await api.createContainer(input);
      notifications.show({
        message: container ? "Location updated" : "Location added",
        color: "green",
      });
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
    <Modal opened={opened} onClose={onClose} title={container ? "Edit location" : "Add location"}>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <TextInput label="Name" withAsterisk {...form.getInputProps("name")} />
          <Select label="Type" data={TYPES} allowDeselect={false} {...form.getInputProps("type")} />
          <Select
            label="Inside (parent location)"
            placeholder="Top level"
            data={parentOptions}
            clearable
            searchable
            {...form.getInputProps("parentId")}
          />
          <Group justify="flex-end" mt="sm">
            <Button variant="default" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit">{container ? "Save" : "Add"}</Button>
          </Group>
        </Stack>
      </form>
    </Modal>
  );
}
