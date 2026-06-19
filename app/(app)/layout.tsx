"use client";

import { useEffect, useState } from "react";
import type { ReactNode } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import {
  AppShell,
  Burger,
  Button,
  Center,
  Group,
  Loader,
  NavLink,
  Title,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  IconBox,
  IconLayoutDashboard,
  IconLogout,
  IconMapPin,
} from "@tabler/icons-react";
import { api } from "@/lib/api";

const NAV = [
  { href: "/", label: "Dashboard", icon: IconLayoutDashboard },
  { href: "/items", label: "Items", icon: IconBox },
  { href: "/locations", label: "Locations", icon: IconMapPin },
];

export default function AuthedLayout({ children }: { children: ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [opened, { toggle, close }] = useDisclosure();
  const [authed, setAuthed] = useState<boolean | null>(null);

  // Guard the authed area: the API enforces auth on every call, but checking
  // up front avoids rendering the shell for signed-out visitors.
  useEffect(() => {
    api
      .session()
      .then((s) => (s.authenticated ? setAuthed(true) : router.replace("/login")))
      .catch(() => router.replace("/login"));
  }, [router]);

  async function handleLogout() {
    await api.logout().catch(() => undefined);
    router.replace("/login");
  }

  if (authed === null) {
    return (
      <Center mih="100vh">
        <Loader />
      </Center>
    );
  }

  return (
    <AppShell
      header={{ height: 56 }}
      navbar={{ width: 240, breakpoint: "sm", collapsed: { mobile: !opened } }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md">
          <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
          <Title order={4}>Home Inventory</Title>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p="md">
        <AppShell.Section grow>
          {NAV.map((item) => (
            <NavLink
              key={item.href}
              component={Link}
              href={item.href}
              label={item.label}
              leftSection={<item.icon size={18} />}
              active={pathname === item.href}
              onClick={close}
            />
          ))}
        </AppShell.Section>
        <Button
          variant="light"
          color="gray"
          leftSection={<IconLogout size={18} />}
          onClick={handleLogout}
        >
          Sign out
        </Button>
      </AppShell.Navbar>

      <AppShell.Main>{children}</AppShell.Main>
    </AppShell>
  );
}
