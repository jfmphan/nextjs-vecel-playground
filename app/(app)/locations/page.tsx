"use client";

import { useCallback, useEffect, useState } from "react";
import type { ReactNode } from "react";
import {
  ActionIcon,
  Badge,
  Button,
  Card,
  Center,
  Group,
  Loader,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { IconEdit, IconPlus, IconTrash } from "@tabler/icons-react";
import { notifications } from "@mantine/notifications";
import { ContainerFormModal } from "@/components/ContainerFormModal";
import { api, ApiError } from "@/lib/api";
import type { Container } from "@/lib/types";

export default function LocationsPage() {
  const [containers, setContainers] = useState<Container[]>([]);
  const [loading, setLoading] = useState(true);
  const [opened, { open, close }] = useDisclosure(false);
  const [editing, setEditing] = useState<Container | null>(null);

  const load = useCallback(async () => {
    setContainers(await api.listContainers());
  }, []);

  useEffect(() => {
    load().finally(() => setLoading(false));
  }, [load]);

  function openAdd() {
    setEditing(null);
    open();
  }
  function openEdit(container: Container) {
    setEditing(container);
    open();
  }

  async function remove(container: Container) {
    if (
      !confirm(
        `Delete "${container.name}"? Items and sub-locations will be detached, not deleted.`,
      )
    ) {
      return;
    }
    try {
      await api.deleteContainer(container.id);
      notifications.show({ message: "Location deleted", color: "green" });
      load();
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

  const roots = containers.filter((c) => c.parentId == null);
  const childrenOf = (id: number) => containers.filter((c) => c.parentId === id);

  function renderNode(container: Container, depth: number): ReactNode {
    return (
      <div key={container.id}>
        <Group justify="space-between" py={6} pl={depth * 24} wrap="nowrap">
          <Group gap="xs">
            <Text fw={500}>{container.name}</Text>
            <Badge size="xs" variant="light" color="gray">
              {container.type}
            </Badge>
            {container.itemCount > 0 && (
              <Badge size="xs" variant="light">
                {container.itemCount} items
              </Badge>
            )}
          </Group>
          <Group gap="xs" wrap="nowrap">
            <ActionIcon variant="subtle" onClick={() => openEdit(container)}>
              <IconEdit size={16} />
            </ActionIcon>
            <ActionIcon variant="subtle" color="red" onClick={() => remove(container)}>
              <IconTrash size={16} />
            </ActionIcon>
          </Group>
        </Group>
        {childrenOf(container.id).map((child) => renderNode(child, depth + 1))}
      </div>
    );
  }

  return (
    <Stack gap="md">
      <Group justify="space-between">
        <Title order={2}>Locations</Title>
        <Button leftSection={<IconPlus size={16} />} onClick={openAdd}>
          Add location
        </Button>
      </Group>

      <Card withBorder radius="md" padding="md">
        {containers.length === 0 ? (
          <Text c="dimmed">No locations yet. Add a room to get started.</Text>
        ) : (
          <Stack gap={0}>{roots.map((c) => renderNode(c, 0))}</Stack>
        )}
      </Card>

      <ContainerFormModal
        opened={opened}
        onClose={close}
        onSaved={load}
        container={editing}
        containers={containers}
      />
    </Stack>
  );
}
