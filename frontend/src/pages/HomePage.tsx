import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { useSnapshot } from 'valtio'
import { appStore } from '@/stores/appStore'
import { 
  Wand2, 
  FileText, 
  Images, 
  History, 
  TrendingUp, 
  Clock,
  Star,
  Zap
} from 'lucide-react'

export function HomePage() {
  const snap = useSnapshot(appStore)

  const quickActions = [
    {
      title: '生成图片',
      description: '使用AI快速生成高质量图片',
      icon: Wand2,
      to: '/generate',
      color: 'bg-purple-500',
    },
    {
      title: '管理提示词',
      description: '创建和管理你的提示词库',
      icon: FileText,
      to: '/prompts',
      color: 'bg-blue-500',
    },
    {
      title: '图片库',
      description: '浏览和管理生成的图片',
      icon: Images,
      to: '/gallery',
      color: 'bg-green-500',
    },
    {
      title: '生成历史',
      description: '查看历史生成记录',
      icon: History,
      to: '/history',
      color: 'bg-orange-500',
    },
  ]

  const stats = [
    {
      label: '总生成次数',
      value: snap.generations.total,
      icon: Zap,
      color: 'text-purple-600',
    },
    {
      label: '提示词数量',
      value: snap.prompts.total,
      icon: FileText,
      color: 'text-blue-600',
    },
    {
      label: '图片收藏',
      value: snap.images.total,
      icon: Star,
      color: 'text-yellow-600',
    },
    {
      label: '批量任务',
      value: snap.batchJobs.total,
      icon: Clock,
      color: 'text-green-600',
    },
  ]

  return (
    <div className="max-w-6xl mx-auto">
      {/* 欢迎区域 */}
      <div className="text-center mb-12">
        <div className="mb-6">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-purple-500 to-pink-600 rounded-full text-white text-2xl font-bold mb-4">
            NB
          </div>
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            Nano Banana AI 图片生成器
          </h1>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            基于先进AI技术的智能图片生成平台，支持文本生成图片、图生图、批量处理等功能
          </p>
        </div>

        <div className="flex justify-center gap-4">
          <Button asChild size="lg">
            <Link to="/generate">
              <Wand2 className="h-5 w-5 mr-2" />
              开始创作
            </Link>
          </Button>
          <Button asChild variant="outline" size="lg">
            <Link to="/gallery">
              <Images className="h-5 w-5 mr-2" />
              浏览作品
            </Link>
          </Button>
        </div>
      </div>

      {/* 统计数据 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
        {stats.map((stat) => {
          const Icon = stat.icon
          return (
            <div key={stat.label} className="bg-white rounded-lg border p-6 text-center">
              <div className="flex items-center justify-center mb-3">
                <Icon className={`h-8 w-8 ${stat.color}`} />
              </div>
              <div className="text-2xl font-bold text-gray-900 mb-1">
                {stat.value.toLocaleString()}
              </div>
              <div className="text-sm text-gray-600">{stat.label}</div>
            </div>
          )
        })}
      </div>

      {/* 快速操作 */}
      <div className="mb-12">
        <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">快速开始</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {quickActions.map((action) => {
            const Icon = action.icon
            return (
              <Link
                key={action.to}
                to={action.to}
                className="group bg-white rounded-lg border p-6 hover:shadow-md transition-shadow"
              >
                <div className="flex items-center mb-4">
                  <div className={`p-3 rounded-lg ${action.color} text-white mr-4`}>
                    <Icon className="h-6 w-6" />
                  </div>
                </div>
                <h3 className="text-lg font-semibold text-gray-900 mb-2 group-hover:text-blue-600">
                  {action.title}
                </h3>
                <p className="text-gray-600 text-sm">{action.description}</p>
              </Link>
            )
          })}
        </div>
      </div>

      {/* 特色功能 */}
      <div className="bg-white rounded-lg border p-8">
        <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">核心特色</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="text-center">
            <div className="bg-purple-100 rounded-full p-4 w-16 h-16 mx-auto mb-4 flex items-center justify-center">
              <Wand2 className="h-8 w-8 text-purple-600" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">智能生成</h3>
            <p className="text-gray-600 text-sm">
              基于 Google Gemini 2.5 Flash 模型，支持文本描述生成高质量图片
            </p>
          </div>

          <div className="text-center">
            <div className="bg-blue-100 rounded-full p-4 w-16 h-16 mx-auto mb-4 flex items-center justify-center">
              <Images className="h-8 w-8 text-blue-600" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">图生图功能</h3>
            <p className="text-gray-600 text-sm">
              上传现有图片进行AI改造，支持风格转换和内容修改
            </p>
          </div>

          <div className="text-center">
            <div className="bg-green-100 rounded-full p-4 w-16 h-16 mx-auto mb-4 flex items-center justify-center">
              <Clock className="h-8 w-8 text-green-600" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">批量处理</h3>
            <p className="text-gray-600 text-sm">
              支持批量生成任务，队列管理，提升创作效率
            </p>
          </div>
        </div>
      </div>

      {/* 最近活动 */}
      <div className="mt-12 grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg border p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">最近生成</h3>
          {snap.generations.list.slice(0, 3).length > 0 ? (
            <div className="space-y-3">
              {snap.generations.list.slice(0, 3).map((generation) => (
                <div key={generation.id} className="flex items-center gap-3">
                  <div className="w-12 h-12 bg-gray-100 rounded border">
                    {generation.thumbnail_url && (
                      <img
                        src={apiService.images.getThumbnailUrl(generation.thumbnail_url)}
                        alt="Generated"
                        className="w-full h-full object-cover rounded"
                      />
                    )}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 truncate">
                      {generation.prompt_text}
                    </p>
                    <p className="text-xs text-gray-500">
                      {new Date(generation.created_at).toLocaleDateString()}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <Images className="h-12 w-12 text-gray-400 mx-auto mb-3" />
              <p className="text-gray-500">暂无生成记录</p>
              <Button asChild className="mt-3">
                <Link to="/generate">开始生成</Link>
              </Button>
            </div>
          )}
        </div>

        <div className="bg-white rounded-lg border p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">常用提示词</h3>
          {snap.prompts.list.slice(0, 3).length > 0 ? (
            <div className="space-y-3">
              {snap.prompts.list.slice(0, 3).map((prompt) => (
                <div key={prompt.id} className="p-3 bg-gray-50 rounded border">
                  <h4 className="text-sm font-medium text-gray-900 mb-1">
                    {prompt.title}
                  </h4>
                  <p className="text-xs text-gray-600 truncate">
                    {prompt.content}
                  </p>
                  <div className="flex items-center gap-2 mt-2">
                    <span className="text-xs text-gray-500">
                      使用 {prompt.usage_count} 次
                    </span>
                    {prompt.is_favorite && (
                      <Star className="h-3 w-3 text-yellow-500 fill-current" />
                    )}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <FileText className="h-12 w-12 text-gray-400 mx-auto mb-3" />
              <p className="text-gray-500">暂无提示词</p>
              <Button asChild className="mt-3">
                <Link to="/prompts">创建提示词</Link>
              </Button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}