import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import Header from '@/components/Header'
import LandingPage from '@/components/LandingPage'
import AdminApp from '@/components/AdminApp'
import { UserpbUserRole } from '@/api/api'

// Admin Route Guard Component
const AdminRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, loading, isAuthenticated } = useAuth()
  
  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="w-8 h-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
      </div>
    )
  }
  
  if (!isAuthenticated || user?.role !== UserpbUserRole.UserRole_ADMIN) {
    return <Navigate to="/" replace />
  }
  
  return <>{children}</>
}

const AppRouter: React.FC = () => {
  return (
    <Router>
      <Routes>
        {/* Public Routes */}
        <Route 
          path="/" 
          element={
            <div className="min-h-screen flex flex-col">
              <Header />
              <main className="flex-1">
                <LandingPage />
              </main>
            </div>
          } 
        />
        
        {/* Admin Routes */}
        <Route 
          path="/admin/*" 
          element={
            <AdminRoute>
              <AdminApp />
            </AdminRoute>
          } 
        />
        
        {/* Catch all route */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  )
}

export default AppRouter
