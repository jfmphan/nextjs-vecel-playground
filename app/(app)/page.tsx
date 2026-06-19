"use client";

import { useEffect, useState } from "react";
import {
  Card,
  Center,
  Group,
  Loader,
  SimpleGrid,
  Stack,
  Table,
  Text,
  Title,
} from "@mantine/core";
import { api } from "@/lib/api";
import type { Item, Stats } from "@/lib/types";

export default function DashboardPage() {
  const [stats, setStats] = useState<Stats | null>(null);

  useEffect(() => {
    api.stats().then(setStats).catch(() => setStats(null));
  }, []);

  if (!stats) {
    return (
      <Center mih={200}>
        <Loader />
      </Center>
    );
  }

  return (
    <Stack gap="lg">
      <Title order={2}>Dashboard</Title>

      <SimpleGrid cols={{ base: 2, sm: 4 }}>
        <StatCard label="Items" value={stats.totalItems} />
        <StatCard label="Total quantity" value={stats.totalQuantity} />
        <StatCard label="Low stock" value={stats.lowStockCount} color="orange" />
        <StatCard label="Expiring soon" value={stats.expiringCount} color="red" />
      </SimpleGrid>

      <ItemListCard
        title="Low on stock"
        items={stats.lowStock}
        emptyText="Nothing is low on stock."
      />
      <ItemListCard
        title="Expiring within 30 days"
        items={stats.expiring}
        emptyText="Nothing is expiring soon."
      />
    </Stack>
  );
}

function StatCard({
  label,
  value,
  color,
}: {
  label: string;
  value: number;
  color?: string;
}) {
  return (
    <Card withBorder radius="md" padding="md">
      <Text size="sm" c="dimmed">
        {label}
      </Text>
      <Text fz={28} fw={700} c={color}>
        {value}
      </Text>
    </Card>
  );
}

function ItemListCard({
  title,
  items,
  emptyText,
}: {
  title: string;
  items: Item[];
  emptyText: string;
}) {
  return (
    <Card withBorder radius="md" padding="md">
      <Title order={4} mb="sm">
        {title}
      </Title>
      {items.length === 0 ? (
        <Text c="dimmed" size="sm">
          {emptyText}
        </Text>
      ) : (
        <Table highlightOnHover>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>Name</Table.Th>
              <Table.Th>Quantity</Table.Th>
              <Table.Th>Expiry</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {items.map((item) => (
              <Table.Tr key={item.id}>
                <Table.Td>{item.name}</Table.Td>
                <Table.Td>
                  <Group gap="xs">
                    <span>
                      {item.quantity}
                      {item.unit ? ` ${item.unit}` : ""}
                    </span>
                  </Group>
                </Table.Td>
                <Table.Td>{item.expiryDate ?? "—"}</Table.Td>
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      )}
    </Card>
  );
}
