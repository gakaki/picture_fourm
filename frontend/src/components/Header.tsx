import { Menu, Settings, User } from 'lucide-react'
import { Button } from './ui/button'
import { useSnapshot } from 'valtio'
import { appStore, appActions } from '@/stores/appStore'

export function Header() {
  const snap = useSnapshot(appStore)

  return (
    <header className="h-16 bg-white shadow-sm border-b flex items-center justify-between px-6">
      <div className="flex items-center gap-4">
        {/* 侧边栏切换按钮 */}
        <Button
          variant="ghost"
          size="icon"
          onClick={appActions.toggleSidebar}
        >
          <Menu className="h-5 w-5" />
        </Button>
        
        {/* Logo */}
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-gradient-to-br from-purple-500 to-pink-600 rounded-lg flex items-center justify-center text-white font-bold text-sm">
            NB
          </div>
          <h1 className="text-xl font-bold text-gray-800">Nano Banana</h1>
        </div>
      </div>

      <div className="flex items-center gap-2">
        {/* 设置按钮 */}
        <Button variant="ghost" size="icon">
          <Settings className="h-5 w-5" />
        </Button>

        {/* 用户按钮 */}
        <Button variant="ghost" size="icon">
          <User className="h-5 w-5" />
        </Button>

        {/* 当前状态显示 */}
        {snap.loading && (
          <div className="flex items-center gap-2 text-sm text-gray-600">
            <div className="w-4 h-4 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />
            处理中...
          </div>
        )}
      </div>
    </header>
  )
}