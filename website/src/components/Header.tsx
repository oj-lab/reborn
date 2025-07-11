import React from 'react'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { LogOut, Settings } from 'lucide-react'
import { GitHubIcon } from '@/components/icons/GitHubIcon'
import { ModeToggle } from '@/components/mode-toggle'
import { LanguageSwitcher } from '@/components/LanguageSwitcher'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import { UserpbUserRole } from '@/api/api'

const Header: React.FC = () => {
  const { user, loading, login, logout, isAuthenticated } = useAuth()
  const { t } = useTranslation()
  const navigate = useNavigate()

  const isAdmin = user?.role === UserpbUserRole.UserRole_ADMIN

  const handleAdminPanel = () => {
    navigate('/admin')
  }

  const handleHome = () => {
    navigate('/')
  }

  return (
    <header className="border-b backdrop-blur-sm bg-background/80 sticky top-0 z-50">
      <div className="container mx-auto px-4 py-4 flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <h1 
            className="text-2xl font-bold text-foreground hover:text-primary transition-colors duration-300 cursor-pointer"
            onClick={handleHome}
          >
            OJ Lab
          </h1>
        </div>
        
        <div className="flex items-center space-x-4">
          <LanguageSwitcher />
          <ModeToggle />
          
          {loading ? (
            <div className="w-8 h-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          ) : isAuthenticated && user ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                  <Avatar className="h-8 w-8">
                    <AvatarImage src="" alt={user.name || 'User'} />
                    <AvatarFallback>
                      {user.name ? user.name.charAt(0).toUpperCase() : 'U'}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-56" align="end" forceMount>
                <div className="flex items-center justify-start gap-2 p-2">
                  <div className="flex flex-col space-y-1 leading-none">
                    {user.name && (
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
                            {t('common.admin')}
                          </Badge>
                        )}
                      </div>
                    )}
                    {user.email && (
                      <p className="w-[200px] truncate text-sm text-muted-foreground">
                        {user.email}
                      </p>
                    )}
                  </div>
                </div>
                <DropdownMenuSeparator />
                {isAdmin && (
                  <DropdownMenuItem onClick={handleAdminPanel}>
                    <Settings className="mr-2 h-4 w-4" />
                    <span>{t('nav.admin', 'Admin Panel')}</span>
                  </DropdownMenuItem>
                )}
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={logout}>
                  <LogOut className="mr-2 h-4 w-4" />
                  <span>{t('auth.logout', 'Log out')}</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <Button onClick={login} className="flex items-center space-x-2">
              <GitHubIcon className="h-4 w-4" />
              <span>Login with GitHub</span>
            </Button>
          )}
        </div>
      </div>
    </header>
  )
}

export default Header
