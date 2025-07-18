import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { MoreHorizontal, Plus, Search, Users, Mail, Calendar, Clock, Loader2 } from 'lucide-react'
import { UserApi } from '@/api'
import type { UserpbUser, UserpbListUsersResponse, TimestamppbTimestamp } from '@/api'

// API function to fetch users using generated client
const fetchUsers = async (page: number = 1, pageSize: number = 10): Promise<UserpbListUsersResponse> => {
  const userApi = new UserApi()
  const response = await userApi.userListGet(page, pageSize)
  return response.data
}

export default function UserManagement() {
  const { t } = useTranslation()
  const [users, setUsers] = useState<UserpbUser[]>([])
  const [totalUsers, setTotalUsers] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState('')
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false)
  const [newUser, setNewUser] = useState({
    name: '',
    email: '',
    role: 0 as number
  })

  // Load users on component mount
  useEffect(() => {
    const loadUsers = async () => {
      try {
        setLoading(true)
        setError(null)
        const data = await fetchUsers(1, 10) // default values
        setUsers(data.users || [])
        setTotalUsers(data.total || 0)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load users')
        console.error('Error loading users:', err)
      } finally {
        setLoading(false)
      }
    }

    loadUsers()
  }, [])

  // Filter users
  const filteredUsers = users.filter(user =>
    (user.name && user.name.toLowerCase().includes(searchTerm.toLowerCase())) ||
    (user.email && user.email.toLowerCase().includes(searchTerm.toLowerCase()))
  )

  // Get role badge style
  const getRoleBadge = (role?: number) => {
    switch (role) {
      case 1: // admin
        return <Badge variant="destructive">{t('common.admin')}</Badge>
      case 0: // user
        return <Badge variant="outline">{t('common.user')}</Badge>
      default:
        return <Badge variant="outline">Unknown</Badge>
    }
  }

  // Format timestamp
  const formatTimestamp = (timestamp?: TimestamppbTimestamp) => {
    if (!timestamp || !timestamp.seconds) return '-'
    const date = new Date(timestamp.seconds * 1000)
    return date.toLocaleDateString()
  }

  // Add user (placeholder - would need backend API)
  const handleAddUser = () => {
    // TODO: Implement add user API call
    console.log('Add user:', newUser)
    setNewUser({ name: '', email: '', role: 0 })
    setIsAddDialogOpen(false)
  }

  // Delete user (placeholder - would need backend API)
  const handleDeleteUser = (userId?: number) => {
    if (!userId) return
    // TODO: Implement delete user API call
    console.log('Delete user:', userId)
  }

  // Toggle user status (placeholder - would need backend API)
  const toggleUserStatus = (userId?: number) => {
    if (!userId) return
    // TODO: Implement toggle user status API call
    console.log('Toggle user status:', userId)
  }

  // Mobile user card component
  const MobileUserCard = ({ user }: { user: UserpbUser }) => (
    <Card className="mb-4">
      <CardContent className="p-4">
        <div className="flex items-start justify-between">
          <div className="flex items-center space-x-3 flex-1">
            <Avatar>
              <AvatarFallback>
                {user.name ? user.name.substring(0, 2).toUpperCase() : 'U'}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <div className="font-medium">{user.name}</div>
              <div className="text-sm text-muted-foreground flex items-center">
                <Mail className="h-3 w-3 mr-1" />
                {user.email || '-'}
              </div>
            </div>
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">{t('userManagement.openMenu')}</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>{t('common.actions')}</DropdownMenuLabel>
              <DropdownMenuItem
                onClick={() => navigator.clipboard.writeText(user.email || '')}
              >
                {t('userManagement.copyEmail')}
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => toggleUserStatus(user.id)}>
                {t('userManagement.toggleUser')}
              </DropdownMenuItem>
              <DropdownMenuItem>{t('userManagement.editUser')}</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem 
                className="text-destructive"
                onClick={() => handleDeleteUser(user.id)}
              >
                {t('userManagement.deleteUser')}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
        <div className="mt-3 flex flex-wrap gap-2">
          {getRoleBadge(user.role)}
        </div>
        <div className="mt-3 grid grid-cols-2 gap-4 text-sm text-muted-foreground">
          <div className="flex items-center">
            <Calendar className="h-3 w-3 mr-1" />
            {t('common.createdAt')}: {formatTimestamp(user.created_at)}
          </div>
          <div className="flex items-center">
            <Clock className="h-3 w-3 mr-1" />
            {t('common.updatedAt')}: {formatTimestamp(user.updated_at)}
          </div>
        </div>
      </CardContent>
    </Card>
  )

  return (
    <div className="container mx-auto p-3 md:p-6 space-y-4 md:space-y-6">
      {/* Header */}
      <div className="flex flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold tracking-tight">{t('userManagement.title')}</h1>
          <p className="text-sm md:text-base text-muted-foreground">
            {t('userManagement.subtitle')}
          </p>
        </div>
        <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
          <DialogTrigger asChild>
            <Button className="w-full md:w-auto">
              <Plus className="mr-2 h-4 w-4" />
              {t('userManagement.addUser')}
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px] mx-4">
            <DialogHeader>
              <DialogTitle>{t('userManagement.addUserTitle')}</DialogTitle>
              <DialogDescription>
                {t('userManagement.addUserSubtitle')}
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-1 md:grid-cols-4 items-center gap-4">
                <Label htmlFor="name" className="md:text-right">
                  {t('common.name')}
                </Label>
                <Input
                  id="name"
                  value={newUser.name}
                  onChange={(e) => setNewUser({ ...newUser, name: e.target.value })}
                  className="md:col-span-3"
                />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-4 items-center gap-4">
                <Label htmlFor="email" className="md:text-right">
                  {t('common.email')}
                </Label>
                <Input
                  id="email"
                  type="email"
                  value={newUser.email}
                  onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
                  className="md:col-span-3"
                />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-4 items-center gap-4">
                <Label htmlFor="role" className="md:text-right">
                  {t('common.role')}
                </Label>
                <select
                  id="role"
                  value={newUser.role}
                  onChange={(e) => setNewUser({ ...newUser, role: parseInt(e.target.value) })}
                  className="md:col-span-3 flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <option value={0}>{t('userManagement.selectUser')}</option>
                  <option value={1}>{t('userManagement.selectAdmin')}</option>
                </select>
              </div>
            </div>
            <DialogFooter>
              <Button type="submit" onClick={handleAddUser}>
                {t('userManagement.addUser')}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {/* Statistics cards */}
      <div className="grid gap-3 grid-cols-2 md:grid-cols-2 lg:grid-cols-4 md:gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-xs md:text-sm font-medium">
              {t('userManagement.totalUsers')}
            </CardTitle>
            <Users className="h-3 w-3 md:h-4 md:w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-xl md:text-2xl font-bold">{totalUsers}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-xs md:text-sm font-medium">
              {t('userManagement.adminUsers')}
            </CardTitle>
            <div className="h-3 w-3 md:h-4 md:w-4 bg-red-500 rounded-full" />
          </CardHeader>
          <CardContent>
            <div className="text-xl md:text-2xl font-bold">
              {users.filter(user => user.role === 1).length}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* User list */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg md:text-xl">{t('userManagement.userList')}</CardTitle>
          <CardDescription className="text-sm">
            {t('userManagement.userListSubtitle')}
          </CardDescription>
          <div className="flex items-center space-x-2">
            <Search className="h-4 w-4 text-muted-foreground" />
            <Input
              placeholder={t('userManagement.searchPlaceholder')}
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full md:max-w-sm"
            />
          </div>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="flex items-center justify-center py-8">
              <Loader2 className="h-6 w-6 animate-spin mr-2" />
              <span>{t('common.loading')}</span>
            </div>
          ) : error ? (
            <div className="text-center py-8 text-red-600">
              <p>{error}</p>
              <Button 
                variant="outline" 
                className="mt-2"
                onClick={() => window.location.reload()}
              >
                {t('common.retry')}
              </Button>
            </div>
          ) : (
            <>
              {/* Mobile card view */}
              <div className="md:hidden">
                {filteredUsers.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    {searchTerm ? t('userManagement.noUsersFound') : t('userManagement.noUsers')}
                  </div>
                ) : (
                  filteredUsers.map((user) => (
                    <MobileUserCard key={user.id} user={user} />
                  ))
                )}
              </div>

              {/* Desktop table view */}
              <div className="hidden md:block">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{t('common.user')}</TableHead>
                      <TableHead>{t('common.role')}</TableHead>
                      <TableHead>{t('common.createdAt')}</TableHead>
                      <TableHead>{t('common.updatedAt')}</TableHead>
                      <TableHead className="text-right">{t('common.actions')}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filteredUsers.map((user) => (
                      <TableRow key={user.id}>
                        <TableCell className="font-medium">
                          <div className="flex items-center space-x-3">
                            <Avatar>
                              <AvatarFallback>
                                {user.name ? user.name.substring(0, 2).toUpperCase() : 'U'}
                              </AvatarFallback>
                            </Avatar>
                            <div>
                              <div className="font-medium">{user.name || 'Unknown'}</div>
                              <div className="text-sm text-muted-foreground">
                                {user.email || '-'}
                              </div>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>{getRoleBadge(user.role)}</TableCell>
                        <TableCell>{formatTimestamp(user.created_at)}</TableCell>
                        <TableCell>{formatTimestamp(user.updated_at)}</TableCell>
                        <TableCell className="text-right">
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" className="h-8 w-8 p-0">
                                <span className="sr-only">{t('userManagement.openMenu')}</span>
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuLabel>{t('common.actions')}</DropdownMenuLabel>
                              <DropdownMenuItem
                                onClick={() => navigator.clipboard.writeText(user.email || '')}
                              >
                                {t('userManagement.copyEmail')}
                              </DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem onClick={() => toggleUserStatus(user.id)}>
                                {t('userManagement.toggleUser')}
                              </DropdownMenuItem>
                              <DropdownMenuItem>{t('userManagement.editUser')}</DropdownMenuItem>
                              <DropdownMenuSeparator />
                              <DropdownMenuItem 
                                className="text-destructive"
                                onClick={() => handleDeleteUser(user.id)}
                              >
                                {t('userManagement.deleteUser')}
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

