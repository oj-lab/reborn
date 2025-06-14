import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Users, UserCheck, Shield, Activity, TrendingUp, Database, Settings } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { DailyActiveUsersChart } from '@/components/DailyActiveUsersChart'

export default function Dashboard() {
  const { t } = useTranslation()
  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{t('dashboard.title')}</h1>
        <p className="text-muted-foreground">
          {t('dashboard.subtitle')}
        </p>
      </div>

      {/* Statistics cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t('dashboard.totalUsers')}
            </CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">2,350</div>
            <p className="text-xs text-muted-foreground">
              +180 {t('dashboard.lastMonth')}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t('dashboard.activeUsers')}
            </CardTitle>
            <UserCheck className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">2,103</div>
            <p className="text-xs text-muted-foreground">
              +20.1% {t('dashboard.lastMonth')}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t('dashboard.systemLoad')}
            </CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">12.4%</div>
            <p className="text-xs text-muted-foreground">
              -2% {t('dashboard.yesterday')}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t('dashboard.revenue')}
            </CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">¥573,000</div>
            <p className="text-xs text-muted-foreground">
              +19% {t('dashboard.lastMonth')}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Daily active users chart */}
      <DailyActiveUsersChart />

      {/* Recent activities */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>{t('dashboard.recentActivity')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex items-center gap-4">
              <span className="relative flex h-2 w-2 shrink-0">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
              </span>
              <div className="flex-1 space-y-1">
                <p className="text-sm font-medium leading-none">
                  {t('dashboard.userRegistered')}
                </p>
                <p className="text-sm text-muted-foreground">
                  {t('dashboard.userRegistrationDesc', { name: '张三' })}
                </p>
              </div>
              <div className="text-sm text-muted-foreground shrink-0">
                2 {t('dashboard.minutesAgo')}
              </div>
            </div>
            <div className="flex items-center gap-4">
              <span className="relative flex h-2 w-2 shrink-0">
                <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
              </span>
              <div className="flex-1 space-y-1">
                <p className="text-sm font-medium leading-none">
                  {t('dashboard.systemUpdated')}
                </p>
                <p className="text-sm text-muted-foreground">
                  {t('dashboard.version')} 2.1.0
                </p>
              </div>
              <div className="text-sm text-muted-foreground shrink-0">
                1 {t('dashboard.hoursAgo')}
              </div>
            </div>
            <div className="flex items-center gap-4">
              <span className="relative flex h-2 w-2 shrink-0">
                <span className="relative inline-flex rounded-full h-2 w-2 bg-yellow-500"></span>
              </span>
              <div className="flex-1 space-y-1">
                <p className="text-sm font-medium leading-none">
                  {t('dashboard.backupCompleted')}
                </p>
                <p className="text-sm text-muted-foreground">
                  {t('dashboard.databaseBackup')}
                </p>
              </div>
              <div className="text-sm text-muted-foreground shrink-0">
                3 {t('dashboard.hoursAgo')}
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>{t('dashboard.quickActions')}</CardTitle>
            <CardDescription>
              {t('dashboard.quickActionsDescription')}
            </CardDescription>
          </CardHeader>
          <CardContent className="grid gap-2">
            <div className="grid grid-cols-2 gap-2">
              <div className="flex items-center justify-center h-20 bg-muted rounded-lg cursor-pointer hover:bg-muted/80 transition-colors">
                <div className="text-center">
                  <Users className="h-6 w-6 mx-auto mb-1" />
                  <p className="text-sm font-medium">{t('nav.userManagement')}</p>
                </div>
              </div>
              <div className="flex items-center justify-center h-20 bg-muted rounded-lg cursor-pointer hover:bg-muted/80 transition-colors">
                <div className="text-center">
                  <Shield className="h-6 w-6 mx-auto mb-1" />
                  <p className="text-sm font-medium">{t('nav.permissions')}</p>
                </div>
              </div>
              <div className="flex items-center justify-center h-20 bg-muted rounded-lg cursor-pointer hover:bg-muted/80 transition-colors">
                <div className="text-center">
                  <Database className="h-6 w-6 mx-auto mb-1" />
                  <p className="text-sm font-medium">{t('nav.data')}</p>
                </div>
              </div>
              <div className="flex items-center justify-center h-20 bg-muted rounded-lg cursor-pointer hover:bg-muted/80 transition-colors">
                <div className="text-center">
                  <Settings className="h-6 w-6 mx-auto mb-1" />
                  <p className="text-sm font-medium">{t('nav.settings')}</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* System status */}
      <Card>
        <CardHeader>
          <CardTitle>{t('dashboard.systemStatus')}</CardTitle>
          <CardDescription>
            {t('dashboard.systemStatusDescription')}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-3">
            <div className="flex items-center justify-between p-4 border rounded-lg">
              <div className="flex items-center space-x-3">
                <div className="h-3 w-3 bg-green-500 rounded-full"></div>
                <div>
                  <p className="font-medium">{t('dashboard.webService')}</p>
                  <p className="text-sm text-muted-foreground">{t('dashboard.running')}</p>
                </div>
              </div>
              <div className="text-sm text-muted-foreground">99.9%</div>
            </div>
            <div className="flex items-center justify-between p-4 border rounded-lg">
              <div className="flex items-center space-x-3">
                <div className="h-3 w-3 bg-green-500 rounded-full"></div>
                <div>
                  <p className="font-medium">{t('dashboard.database')}</p>
                  <p className="text-sm text-muted-foreground">{t('dashboard.running')}</p>
                </div>
              </div>
              <div className="text-sm text-muted-foreground">99.8%</div>
            </div>
            <div className="flex items-center justify-between p-4 border rounded-lg">
              <div className="flex items-center space-x-3">
                <div className="h-3 w-3 bg-yellow-500 rounded-full"></div>
                <div>
                  <p className="font-medium">{t('dashboard.cacheService')}</p>
                  <p className="text-sm text-muted-foreground">{t('dashboard.delay')}</p>
                </div>
              </div>
              <div className="text-sm text-muted-foreground">98.2%</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
