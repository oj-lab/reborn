"use client"

import { Palette } from "lucide-react"
import { useTranslation } from "react-i18next"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

const themes = [
  { name: "default", color: "hsl(222.2, 84%, 4.9%)" },
  { name: "slate", color: "hsl(215.4, 16.3%, 46.9%)" },
  { name: "gray", color: "hsl(220.9, 39.3%, 11%)" },
  { name: "zinc", color: "hsl(240, 10%, 3.9%)" },
  { name: "neutral", color: "hsl(0, 0%, 3.9%)" },
  { name: "stone", color: "hsl(20, 14.3%, 4.1%)" },
  { name: "red", color: "hsl(0, 72.2%, 50.6%)" },
  { name: "rose", color: "hsl(346.8, 77.2%, 49.8%)" },
  { name: "orange", color: "hsl(24.6, 95%, 53.1%)" },
  { name: "green", color: "hsl(142.1, 76.2%, 36.3%)" },
  { name: "blue", color: "hsl(221.2, 83.2%, 53.3%)" },
  { name: "yellow", color: "hsl(47.9, 95.8%, 53.1%)" },
  { name: "violet", color: "hsl(262.1, 83.3%, 57.8%)" },
]

export function ColorThemeToggle() {
  const { t } = useTranslation()

  const applyTheme = (themeName: string) => {
    // 移除所有主题类
    const root = document.documentElement
    themes.forEach(theme => {
      root.classList.remove(`theme-${theme.name}`)
    })
    
    // 添加新主题类
    if (themeName !== "default") {
      root.classList.add(`theme-${themeName}`)
    }
    
    // 保存到本地存储
    localStorage.setItem("color-theme", themeName)
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="h-8 w-8">
          <Palette className="h-4 w-4" />
          <span className="sr-only">{t('theme.colorTheme')}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-48">
        {themes.map((theme) => (
          <DropdownMenuItem
            key={theme.name}
            onClick={() => applyTheme(theme.name)}
            className="flex items-center gap-2"
          >
            <div
              className="h-4 w-4 rounded-full border border-border"
              style={{ backgroundColor: theme.color }}
            />
            {t(`theme.colors.${theme.name}`, theme.name)}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
