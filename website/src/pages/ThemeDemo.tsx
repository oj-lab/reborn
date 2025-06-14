import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { ModeToggle } from '@/components/mode-toggle'
import { ColorThemeToggle } from '@/components/color-theme-toggle'
import { LanguageSwitcher } from '@/components/LanguageSwitcher'
import { Palette, Sun, Globe } from 'lucide-react'

export default function ThemeDemo() {
  const { t } = useTranslation()

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{t('theme.demo.title', '主题演示')}</h1>
        <p className="text-muted-foreground">
          {t('theme.demo.subtitle', '测试不同的颜色主题和深浅模式效果')}
        </p>
      </div>

      {/* Theme controllers */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Palette className="h-5 w-5" />
            {t('theme.demo.controls', '主题控制')}
          </CardTitle>
          <CardDescription>
            {t('theme.demo.controlsDesc', '切换不同的主题色彩和深浅模式')}
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-4">
          <div className="flex items-center gap-2">
            <Sun className="h-4 w-4" />
            <span className="text-sm font-medium">{t('theme.demo.darkMode', '深浅模式')}:</span>
            <ModeToggle />
          </div>
          <div className="flex items-center gap-2">
            <Palette className="h-4 w-4" />
            <span className="text-sm font-medium">{t('theme.colorTheme')}:</span>
            <ColorThemeToggle />
          </div>
          <div className="flex items-center gap-2">
            <Globe className="h-4 w-4" />
            <span className="text-sm font-medium">{t('theme.demo.language', '语言')}:</span>
            <LanguageSwitcher />
          </div>
        </CardContent>
      </Card>

      {/* Component demos */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {/* Button components */}
        <Card>
          <CardHeader>
            <CardTitle>{t('theme.demo.buttons', '按钮组件')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex flex-wrap gap-2">
              <Button variant="default">Default</Button>
              <Button variant="secondary">Secondary</Button>
              <Button variant="outline">Outline</Button>
              <Button variant="ghost">Ghost</Button>
              <Button variant="destructive">Destructive</Button>
            </div>
          </CardContent>
        </Card>

        {/* Badge components */}
        <Card>
          <CardHeader>
            <CardTitle>{t('theme.demo.badges', '徽章组件')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex flex-wrap gap-2">
              <Badge variant="default">Default</Badge>
              <Badge variant="secondary">Secondary</Badge>
              <Badge variant="outline">Outline</Badge>
              <Badge variant="destructive">Destructive</Badge>
            </div>
          </CardContent>
        </Card>

        {/* Progress bars */}
        <Card>
          <CardHeader>
            <CardTitle>{t('theme.demo.progress', '进度条')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <Progress value={33} className="w-full" />
            <Progress value={66} className="w-full" />
            <Progress value={99} className="w-full" />
          </CardContent>
        </Card>

        {/* Text styles */}
        <Card>
          <CardHeader>
            <CardTitle>{t('theme.demo.typography', '文字样式')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <h1 className="text-2xl font-bold">Heading 1</h1>
            <h2 className="text-xl font-semibold">Heading 2</h2>
            <h3 className="text-lg font-medium">Heading 3</h3>
            <p className="text-base">Regular paragraph text</p>
            <p className="text-sm text-muted-foreground">Muted text</p>
          </CardContent>
        </Card>

        {/* Card variants */}
        <Card>
          <CardHeader>
            <CardTitle>{t('theme.demo.cardVariants', '卡片变体')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="p-3 bg-primary text-primary-foreground rounded">
                Primary background
              </div>
              <div className="p-3 bg-secondary text-secondary-foreground rounded">
                Secondary background
              </div>
              <div className="p-3 bg-muted text-muted-foreground rounded">
                Muted background
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Borders and separators */}
        <Card>
          <CardHeader>
            <CardTitle>{t('theme.demo.borders', '边框和分隔线')}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="border rounded p-2">Border</div>
            <div className="border-2 border-primary rounded p-2">Primary border</div>
            <hr className="border-border" />
            <div className="h-px bg-border"></div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
