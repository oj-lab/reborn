import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '@/hooks/useAuth'
import { UserpbUserRole } from '@/api/api'
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
} from '@/components/ui/sidebar'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { LanguageSwitcher } from '@/components/LanguageSwitcher'
import { ModeToggle } from '@/components/mode-toggle'
import { ColorThemeToggle } from '@/components/color-theme-toggle'
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
  User,
  ChevronUp,
  Palette,
} from 'lucide-react'

interface AdminLayoutProps {
  children: React.ReactNode
  currentRoute?: string
  onRouteChange?: (route: string) => void
}

// Navigation menu items
const useMenuItems = () => {
  const { t } = useTranslation()
  
  return [
    {
      title: t('nav.overview'),
      items: [
        {
          title: t('nav.dashboard'),
          icon: LayoutDashboard,
          url: '/dashboard',
          isActive: true,
        },
      ],
    },
    {
      title: t('nav.userManagement'),
      items: [
        {
          title: t('nav.users'),
          icon: Users,
          url: '/users',
        },
        {
          title: t('nav.permissions'),
          icon: Shield,
          url: '/permissions',
        },
      ],
    },
    {
      title: t('nav.systemManagement'),
      items: [
        {
          title: t('nav.data'),
          icon: Database,
          url: '/data',
        },
        {
          title: t('nav.settings'),
          icon: Settings,
          url: '/settings',
        },
        {
          title: t('nav.analytics'),
          icon: BarChart3,
          url: '/analytics',
        },
        {
          title: t('theme.demo.title', '主题演示'),
          icon: Palette,
          url: '/theme-demo',
        },
      ],
    },
  ]
}

export default function AdminLayout({ children, currentRoute = '/dashboard', onRouteChange }: AdminLayoutProps) {
  const { t } = useTranslation()
  const { user } = useAuth()
  const [activeItem, setActiveItem] = useState(currentRoute)
  const menuItems = useMenuItems()

  const isAdmin = user?.role === UserpbUserRole.UserRole_ADMIN

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
                <span className="truncate text-xs text-muted-foreground">{t('layout.adminSystem')}</span>
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
                          onClick={() => {
                            setActiveItem(item.url)
                            onRouteChange?.(item.url)
                          }}
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
            <SidebarMenu>
              <SidebarMenuItem>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <SidebarMenuButton
                      size="lg"
                      className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                    >
                      <Avatar className="h-8 w-8 rounded-lg">
                        <AvatarImage src="" alt={user?.name || 'User'} />
                        <AvatarFallback className="rounded-lg">
                          {user?.name ? user.name.charAt(0).toUpperCase() : 'U'}
                        </AvatarFallback>
                      </Avatar>
                      <div className="grid flex-1 text-left text-sm leading-tight">
                        <span className="truncate font-semibold">
                          {user?.name || t('layout.adminName')}
                          {user?.id && (
                            <span className="text-muted-foreground font-normal ml-1">
                              #{user.id}
                            </span>
                          )}
                        </span>
                        <span className="truncate text-xs">{user?.email || t('layout.adminEmail')}</span>
                      </div>
                      <ChevronUp className="ml-auto size-4" />
                    </SidebarMenuButton>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent
                    className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                    side="bottom"
                    align="end"
                    sideOffset={4}
                  >
                    <DropdownMenuLabel className="p-0 font-normal">
                      <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                        <Avatar className="h-8 w-8 rounded-lg">
                          <AvatarImage src="" alt={user?.name || 'User'} />
                          <AvatarFallback className="rounded-lg">
                            {user?.name ? user.name.charAt(0).toUpperCase() : 'U'}
                          </AvatarFallback>
                        </Avatar>
                        <div className="grid flex-1 text-left text-sm leading-tight">
                          <div className="flex items-center gap-2 flex-wrap">
                            <span className="truncate font-semibold">
                              {user?.name || t('layout.adminName')}
                              {user?.id && (
                                <span className="text-muted-foreground font-normal ml-1">
                                  #{user.id}
                                </span>
                              )}
                            </span>
                            {isAdmin && (
                              <Badge variant="secondary" className="text-xs">
                                {t('common.admin')}
                              </Badge>
                            )}
                          </div>
                          <span className="truncate text-xs">{user?.email || t('layout.adminEmail')}</span>
                        </div>
                      </div>
                    </DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem>
                      <User className="mr-2 h-4 w-4" />
                      {t('layout.profile')}
                    </DropdownMenuItem>
                    <DropdownMenuItem>
                      <Settings className="mr-2 h-4 w-4" />
                      {t('nav.settings')}
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem>
                      <LogOut className="mr-2 h-4 w-4" />
                      {t('layout.logout')}
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarFooter>
          <SidebarRail />
        </Sidebar>

        {/* Main content area */}
        <div className="flex flex-1 flex-col">
          {/* Top navigation bar */}
          <header className="flex h-16 shrink-0 items-center gap-2 border-b px-3 md:px-4">
            <SidebarTrigger className="-ml-1" />
            <div className="flex flex-1 items-center justify-between">
              <div></div>
              <div className="flex items-center gap-1 md:gap-2">
                {/* Notification button - hidden on small screens */}
                <Button variant="ghost" size="icon" className="h-8 w-8 hidden sm:inline-flex">
                  <Bell className="h-4 w-4" />
                </Button>
                
                {/* Help button - hidden on small screens */}
                <Button variant="ghost" size="icon" className="h-8 w-8 hidden sm:inline-flex">
                  <HelpCircle className="h-4 w-4" />
                </Button>

                {/* Color theme toggle */}
                <ColorThemeToggle />

                {/* Language switcher */}
                <LanguageSwitcher />

                {/* Dark mode toggle */}
                <ModeToggle />
              </div>
            </div>
          </header>

          {/* Main content */}
          <main className="flex-1 overflow-auto p-3 md:p-6">
            {children}
          </main>
        </div>
      </div>
    </SidebarProvider>
  )
}
