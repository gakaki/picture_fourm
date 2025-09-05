import { Sidebar } from './Sidebar'
import { Header } from './Header'
import { useSnapshot } from 'valtio'
import { appStore } from '@/stores/appStore'
import { ToastProvider, ToastViewport } from './ui/toast'

interface LayoutProps {
  children: React.ReactNode
}

export function Layout({ children }: LayoutProps) {
  const snap = useSnapshot(appStore)

  return (
    <ToastProvider>
      <div className="min-h-screen bg-gray-50">
        {/* Header */}
        <Header />
        
        <div className="flex h-[calc(100vh-64px)]">
          {/* Sidebar */}
          <Sidebar />
          
          {/* Main Content */}
          <main
            className={`flex-1 overflow-auto transition-all duration-300 ${
              snap.ui.sidebarOpen ? 'ml-64' : 'ml-0'
            }`}
          >
            <div className="p-6">
              {children}
            </div>
          </main>
        </div>
      </div>
      <ToastViewport />
    </ToastProvider>
  )
}