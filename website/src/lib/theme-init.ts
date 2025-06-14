// Theme initialization function, executed when page loads
export function initializeTheme() {
  // Restore color theme
  const savedColorTheme = localStorage.getItem("color-theme")
  if (savedColorTheme && savedColorTheme !== "default") {
    document.documentElement.classList.add(`theme-${savedColorTheme}`)
  }
}

// Execute immediately when DOM is loaded
if (typeof window !== 'undefined') {
  document.addEventListener('DOMContentLoaded', initializeTheme)
  
  // If DOM is already loaded, execute immediately
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeTheme)
  } else {
    initializeTheme()
  }
}
