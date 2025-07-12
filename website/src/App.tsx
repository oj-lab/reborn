import { AuthProvider } from './contexts/AuthContext'
import { ThemeProvider } from './components/theme-provider'
import AuthCallback from './components/AuthCallback'
import AppRouter from './routes/AppRouter'

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <AuthCallback />
        <AppRouter />
      </AuthProvider>
    </ThemeProvider>
  )
}

export default App
