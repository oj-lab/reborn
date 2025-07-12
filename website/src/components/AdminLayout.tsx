import { useTranslation } from "react-i18next";
import { useNavigate, useLocation } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { UserpbUserRole } from "@/api/api";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarRail,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { LanguageSwitcher } from "@/components/LanguageSwitcher";
import { ModeToggle } from "@/components/mode-toggle";
import { ColorThemeToggle } from "@/components/color-theme-toggle";
import {
  LayoutDashboard,
  Users,
  Settings,
  Shield,
  Database,
  BarChart3,
  Bell,
  HelpCircle,
  LogOut,
  Palette,
  Home,
} from "lucide-react";

interface AdminLayoutProps {
  children: React.ReactNode;
  currentRoute?: string;
}

// Navigation menu items
const useMenuItems = () => {
  const { t } = useTranslation();

  return [
    {
      title: t("nav.overview"),
      items: [
        {
          title: t("nav.dashboard"),
          icon: LayoutDashboard,
          url: "/admin/dashboard",
          isActive: true,
        },
      ],
    },
    {
      title: t("nav.userManagement"),
      items: [
        {
          title: t("nav.users"),
          icon: Users,
          url: "/admin/users",
        },
        {
          title: t("nav.permissions"),
          icon: Shield,
          url: "/admin/permissions",
        },
      ],
    },
    {
      title: t("nav.systemManagement"),
      items: [
        {
          title: t("nav.data"),
          icon: Database,
          url: "/admin/data",
        },
        {
          title: t("nav.settings"),
          icon: Settings,
          url: "/admin/settings",
        },
        {
          title: t("nav.analytics"),
          icon: BarChart3,
          url: "/admin/analytics",
        },
        {
          title: t("theme.demo.title", "主题演示"),
          icon: Palette,
          url: "/admin/theme-demo",
        },
      ],
    },
  ];
};

export default function AdminLayout({
  children,
  currentRoute,
}: AdminLayoutProps) {
  const { t } = useTranslation();
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const menuItems = useMenuItems();

  const activeItem = currentRoute || location.pathname;
  const isAdmin = user?.role === UserpbUserRole.UserRole_ADMIN;

  const handleNavigation = (url: string) => {
    navigate(url);
  };

  const handleLogout = () => {
    logout();
    navigate("/");
  };

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        {/* Sidebar */}
        <Sidebar variant="inset">
          <SidebarHeader>
            <div className="flex items-center gap-2 px-4 py-2">
              <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                <Shield className="size-4" />
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-semibold">Reborn Admin</span>
                <span className="truncate text-xs text-muted-foreground">
                  {t("layout.adminSystem")}
                </span>
              </div>
            </div>
          </SidebarHeader>

          <SidebarContent>
            {menuItems.map((group) => (
              <SidebarGroup key={group.title}>
                <SidebarGroupLabel>{group.title}</SidebarGroupLabel>
                <SidebarGroupContent>
                  <SidebarMenu>
                    {group.items.map((item) => (
                      <SidebarMenuItem key={item.title}>
                        <SidebarMenuButton
                          asChild
                          isActive={activeItem === item.url}
                          onClick={() => handleNavigation(item.url)}
                        >
                          <div className="flex items-center gap-2 cursor-pointer">
                            <item.icon className="h-4 w-4" />
                            <span>{item.title}</span>
                          </div>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    ))}
                  </SidebarMenu>
                </SidebarGroupContent>
              </SidebarGroup>
            ))}
          </SidebarContent>

          <SidebarFooter>
            {/* User menu moved to header, footer can be used for other content if needed */}
          </SidebarFooter>
          <SidebarRail />
        </Sidebar>

        {/* Main content area */}
        <div className="flex flex-1 flex-col">
          {/* Top navigation bar */}
          <header className="flex h-16 shrink-0 items-center gap-2 border-b px-3 md:px-4">
            <SidebarTrigger className="-ml-1" />
            <div className="flex flex-1 items-center justify-between">
              <div className="flex items-center gap-2">
                {/* Left items */}
              </div>
              <div className="flex items-center gap-1 md:gap-2">
                {/* Notification button - hidden on small screens */}
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-8 w-8 hidden sm:inline-flex"
                >
                  <Bell className="h-4 w-4" />
                </Button>

                {/* Help button - hidden on small screens */}
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-8 w-8 hidden sm:inline-flex"
                >
                  <HelpCircle className="h-4 w-4" />
                </Button>

                {/* Color theme toggle */}
                <ColorThemeToggle />

                {/* Language switcher */}
                <LanguageSwitcher />

                {/* Dark mode toggle */}
                <ModeToggle />

                {/* User menu */}
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      className="relative h-8 w-8 rounded-full"
                    >
                      <Avatar className="h-8 w-8">
                        <AvatarImage src="" alt={user?.name || "User"} />
                        <AvatarFallback>
                          {user?.name ? user.name.charAt(0).toUpperCase() : "U"}
                        </AvatarFallback>
                      </Avatar>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent className="w-56" align="end" forceMount>
                    <div className="flex items-center justify-start gap-2 p-2">
                      <div className="flex flex-col space-y-1 leading-none">
                        {user?.name && (
                          <div className="flex items-center gap-2 flex-wrap">
                            <p className="font-medium">
                              {user.name}
                              {user.id && (
                                <span className="text-muted-foreground font-normal ml-1">
                                  #{user.id}
                                </span>
                              )}
                            </p>
                            {isAdmin && (
                              <Badge variant="secondary" className="text-xs">
                                {t("common.admin")}
                              </Badge>
                            )}
                          </div>
                        )}
                        {user?.email && (
                          <p className="w-[200px] truncate text-sm text-muted-foreground">
                            {user.email}
                          </p>
                        )}
                      </div>
                    </div>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem onClick={() => navigate("/")}>
                      <Home className="mr-2 h-4 w-4" />
                      {t("layout.backToHome")}
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem onClick={handleLogout}>
                      <LogOut className="mr-2 h-4 w-4" />
                      {t("layout.logout")}
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          </header>

          {/* Main content */}
          <main className="flex-1 overflow-auto p-3 md:p-6">{children}</main>
        </div>
      </div>
    </SidebarProvider>
  );
}
