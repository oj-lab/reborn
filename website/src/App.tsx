import { AuthProvider } from './contexts/AuthContext'
import { ThemeProvider } from './components/theme-provider'
import Header from './components/Header'
import AuthCallback from './components/AuthCallback'
import LandingPage from './components/LandingPage'

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <AuthCallback />
        <div className="min-h-screen flex flex-col">
          <Header />
          <main className="flex-1">
            <LandingPage />
          </main>
        </div>
      </AuthProvider>
    </ThemeProvider>
  )
}

export default App
