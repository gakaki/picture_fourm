import { NavLink } from 'react-router-dom'
import { useSnapshot } from 'valtio'
import { appStore } from '@/stores/appStore'
import { 
  Home, 
  Wand2, 
  FileText, 
  History, 
  Images, 
  Settings,
  Queue,
  TrendingUp
} from 'lucide-react'

interface NavItem {
  to: string
  icon: React.ComponentType<{ className?: string }>
  label: string
  badge?: number
}

const navItems: NavItem[] = [
  { to: '/', icon: Home, label: '首页' },
  { to: '/generate', icon: Wand2, label: '生成图片' },
  { to: '/prompts', icon: FileText, label: '提示词管理' },
  { to: '/history', icon: History, label: '生成历史' },
  { to: '/gallery', icon: Images, label: '图片库' },
]

export function Sidebar() {
  const snap = useSnapshot(appStore)

  if (!snap.ui.sidebarOpen) {
    return null
  }

  return (
    <aside className="fixed left-0 top-16 z-40 w-64 h-[calc(100vh-64px)] bg-white shadow-lg border-r">
      <div className="p-4">
        <nav className="space-y-2">
          {navItems.map((item) => {
            const Icon = item.icon
            return (
              <NavLink
                key={item.to}
                to={item.to}
                className={({ isActive }) =>
                  `flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                    isActive
                      ? 'bg-blue-50 text-blue-700 border border-blue-200'
                      : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                  }`
                }
              >
                <Icon className="h-5 w-5" />
                {item.label}
                {item.badge && (
                  <span className="ml-auto bg-red-100 text-red-600 text-xs px-2 py-1 rounded-full">
                    {item.badge}
                  </span>
                )}
              </NavLink>
            )
          })}
        </nav>

        {/* 统计信息 */}
        <div className="mt-8 pt-4 border-t">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">
            统计信息
          </h3>
          <div className="space-y-3">
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-600">生成记录</span>
              <span className="font-medium">{snap.generations.total}</span>
            </div>
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-600">提示词</span>
              <span className="font-medium">{snap.prompts.total}</span>
            </div>
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-600">批量任务</span>
              <span className="font-medium">{snap.batchJobs.total}</span>
            </div>
          </div>
        </div>

        {/* 快速操作 */}
        <div className="mt-8 pt-4 border-t">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">
            快速操作
          </h3>
          <div className="space-y-2">
            <button className="w-full flex items-center gap-3 px-3 py-2 text-left text-sm text-gray-600 hover:bg-gray-50 rounded-lg">
              <Queue className="h-4 w-4" />
              查看队列
            </button>
            <button className="w-full flex items-center gap-3 px-3 py-2 text-left text-sm text-gray-600 hover:bg-gray-50 rounded-lg">
              <TrendingUp className="h-4 w-4" />
              生成统计
            </button>
            <button className="w-full flex items-center gap-3 px-3 py-2 text-left text-sm text-gray-600 hover:bg-gray-50 rounded-lg">
              <Settings className="h-4 w-4" />
              系统设置
            </button>
          </div>
        </div>
      </div>
    </aside>
  )
}