// 主题初始化函数，在页面加载时执行
export function initializeTheme() {
  // 恢复颜色主题
  const savedColorTheme = localStorage.getItem("color-theme")
  if (savedColorTheme && savedColorTheme !== "default") {
    document.documentElement.classList.add(`theme-${savedColorTheme}`)
  }
}

// 在 DOM 加载完成时立即执行
if (typeof window !== 'undefined') {
  document.addEventListener('DOMContentLoaded', initializeTheme)
  
  // 如果 DOM 已经加载完成，立即执行
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeTheme)
  } else {
    initializeTheme()
  }
}
