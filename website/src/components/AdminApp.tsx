import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import AdminLayout from '@/components/AdminLayout'
import Dashboard from '@/pages/Dashboard'
import UserManagement from '@/pages/UserManagement'
import ThemeDemo from '@/pages/ThemeDemo'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

// Temporary page component
const PermissionsPage = () => {
  const { t } = useTranslation()
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{t('permissions.title')}</h1>
        <p className="text-muted-foreground">{t('permissions.subtitle')}</p>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{t('permissions.permissionSettings')}</CardTitle>
          <CardDescription>{t('permissions.inDevelopment')}</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">{t('permissions.comingSoon')}</p>
        </CardContent>
      </Card>
    </div>
  )
}

const DataManagementPage = () => {
  const { t } = useTranslation()
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{t('dataManagement.title')}</h1>
        <p className="text-muted-foreground">{t('dataManagement.subtitle')}</p>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{t('dataManagement.databaseStatus')}</CardTitle>
          <CardDescription>{t('dataManagement.inDevelopment')}</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">{t('dataManagement.comingSoon')}</p>
        </CardContent>
      </Card>
    </div>
  )
}

const SettingsPage = () => {
  const { t } = useTranslation()
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{t('settings.title')}</h1>
        <p className="text-muted-foreground">{t('settings.subtitle')}</p>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{t('settings.systemConfig')}</CardTitle>
          <CardDescription>{t('settings.inDevelopment')}</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">{t('settings.comingSoon')}</p>
        </CardContent>
      </Card>
    </div>
  )
}

const AnalyticsPage = () => {
  const { t } = useTranslation()
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{t('analytics.title')}</h1>
        <p className="text-muted-foreground">{t('analytics.subtitle')}</p>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{t('analytics.dataStatistics')}</CardTitle>
          <CardDescription>{t('analytics.inDevelopment')}</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">{t('analytics.comingSoon')}</p>
        </CardContent>
      </Card>
    </div>
  )
}

// Route configuration
const routes = {
  '/dashboard': Dashboard,
  '/users': UserManagement,
  '/permissions': PermissionsPage,
  '/data': DataManagementPage,
  '/settings': SettingsPage,
  '/analytics': AnalyticsPage,
  '/theme-demo': ThemeDemo,
}

export default function AdminApp() {
  const [currentRoute, setCurrentRoute] = useState('/dashboard')

  // Get current page component
  const getCurrentPage = () => {
    const PageComponent = routes[currentRoute as keyof typeof routes] || Dashboard
    return <PageComponent />
  }

  return (
    <AdminLayout currentRoute={currentRoute} onRouteChange={setCurrentRoute}>
      {getCurrentPage()}
    </AdminLayout>
  )
}
