import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { apiService } from '@/services/api'
import { appActions } from '@/stores/appStore'
import { Wand2, Upload, Settings, Download } from 'lucide-react'
import * as Select from '@radix-ui/react-select'
import * as Tabs from '@radix-ui/react-tabs'

export function GeneratePage() {
  const [prompt, setPrompt] = useState('')
  const [count, setCount] = useState(1)
  const [size, setSize] = useState('1024x1024')
  const [quality, setQuality] = useState('standard')
  const [strength, setStrength] = useState(0.8)
  const [sourceImage, setSourceImage] = useState<string>('')
  const [activeTab, setActiveTab] = useState<'text2img' | 'img2img'>('text2img')

  const queryClient = useQueryClient()

  // 文本生成图片
  const text2imgMutation = useMutation({
    mutationFn: (data: { prompt: string; count: number; params: any }) =>
      apiService.generate.text2img(data),
    onMutate: () => {
      appActions.setLoading(true)
      appActions.setError(null)
    },
    onSuccess: (response) => {
      appActions.setLoading(false)
      if (response.success) {
        // 将新生成的图片添加到store
        response.data?.forEach(generation => {
          appActions.addGeneration(generation)
        })
        // 清空输入
        setPrompt('')
        // 刷新相关查询
        queryClient.invalidateQueries({ queryKey: ['generations'] })
      }
    },
    onError: (error: any) => {
      appActions.setLoading(false)
      appActions.setError(error.message || '生成失败')
    }
  })

  // 图生图
  const img2imgMutation = useMutation({
    mutationFn: (data: { prompt: string; source_image: string; count: number; params: any }) =>
      apiService.generate.img2img(data),
    onMutate: () => {
      appActions.setLoading(true)
      appActions.setError(null)
    },
    onSuccess: (response) => {
      appActions.setLoading(false)
      if (response.success) {
        response.data?.forEach(generation => {
          appActions.addGeneration(generation)
        })
        setPrompt('')
        setSourceImage('')
        queryClient.invalidateQueries({ queryKey: ['generations'] })
      }
    },
    onError: (error: any) => {
      appActions.setLoading(false)
      appActions.setError(error.message || '生成失败')
    }
  })

  // 处理文件上传
  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        setSourceImage(e.target?.result as string)
      }
      reader.readAsDataURL(file)
    }
  }

  // 生成图片
  const handleGenerate = () => {
    if (!prompt.trim()) {
      appActions.setError('请输入提示词')
      return
    }

    const params = {
      size,
      quality,
      ...(activeTab === 'img2img' && { strength })
    }

    if (activeTab === 'text2img') {
      text2imgMutation.mutate({ prompt, count, params })
    } else if (activeTab === 'img2img') {
      if (!sourceImage) {
        appActions.setError('请上传源图片')
        return
      }
      img2imgMutation.mutate({ prompt, source_image: sourceImage, count, params })
    }
  }

  const isLoading = text2imgMutation.isPending || img2imgMutation.isPending

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">AI 图片生成</h1>
        <p className="text-gray-600">
          使用先进的AI模型生成高质量图片，支持文本生成图片和图生图功能
        </p>
      </div>

      <Tabs.Root value={activeTab} onValueChange={(value) => setActiveTab(value as any)}>
        <Tabs.List className="flex space-x-1 bg-gray-100 p-1 rounded-lg mb-6">
          <Tabs.Trigger
            value="text2img"
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md data-[state=active]:bg-white data-[state=active]:shadow-sm"
          >
            <Wand2 className="h-4 w-4" />
            文本生成图片
          </Tabs.Trigger>
          <Tabs.Trigger
            value="img2img"
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md data-[state=active]:bg-white data-[state=active]:shadow-sm"
          >
            <Upload className="h-4 w-4" />
            图生图
          </Tabs.Trigger>
        </Tabs.List>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 左侧输入区域 */}
          <div className="lg:col-span-2 space-y-6">
            <Tabs.Content value="text2img">
              <div className="bg-white rounded-lg border p-6">
                <h2 className="text-lg font-semibold mb-4">文本生成图片</h2>
                
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      提示词 *
                    </label>
                    <Textarea
                      placeholder="输入你想要生成的图片描述，例如：一只可爱的小猫坐在窗台上，阳光透过窗户洒在它身上..."
                      value={prompt}
                      onChange={(e) => setPrompt(e.target.value)}
                      className="min-h-[100px]"
                    />
                    <p className="text-xs text-gray-500 mt-1">
                      详细的描述能帮助生成更准确的图片
                    </p>
                  </div>
                  
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        生成数量
                      </label>
                      <Input
                        type="number"
                        min="1"
                        max="4"
                        value={count}
                        onChange={(e) => setCount(parseInt(e.target.value) || 1)}
                      />
                    </div>
                  </div>
                </div>
              </div>
            </Tabs.Content>

            <Tabs.Content value="img2img">
              <div className="bg-white rounded-lg border p-6">
                <h2 className="text-lg font-semibold mb-4">图生图</h2>
                
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      源图片 *
                    </label>
                    <div className="border-2 border-dashed border-gray-300 rounded-lg p-4">
                      {sourceImage ? (
                        <div className="text-center">
                          <img
                            src={sourceImage}
                            alt="Source"
                            className="max-h-40 mx-auto rounded mb-2"
                          />
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setSourceImage('')}
                          >
                            重新选择
                          </Button>
                        </div>
                      ) : (
                        <div className="text-center">
                          <Upload className="h-8 w-8 text-gray-400 mx-auto mb-2" />
                          <input
                            type="file"
                            accept="image/*"
                            onChange={handleFileUpload}
                            className="hidden"
                            id="image-upload"
                          />
                          <label
                            htmlFor="image-upload"
                            className="cursor-pointer text-blue-600 hover:text-blue-500"
                          >
                            点击上传图片
                          </label>
                          <p className="text-xs text-gray-500 mt-1">
                            支持 JPG、PNG 格式
                          </p>
                        </div>
                      )}
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      修改提示词 *
                    </label>
                    <Textarea
                      placeholder="描述你想要对图片进行的修改..."
                      value={prompt}
                      onChange={(e) => setPrompt(e.target.value)}
                      className="min-h-[80px]"
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        变化强度
                      </label>
                      <Input
                        type="number"
                        min="0"
                        max="1"
                        step="0.1"
                        value={strength}
                        onChange={(e) => setStrength(parseFloat(e.target.value) || 0.8)}
                      />
                      <p className="text-xs text-gray-500 mt-1">
                        0.1-0.9，越高变化越大
                      </p>
                    </div>
                    
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        生成数量
                      </label>
                      <Input
                        type="number"
                        min="1"
                        max="4"
                        value={count}
                        onChange={(e) => setCount(parseInt(e.target.value) || 1)}
                      />
                    </div>
                  </div>
                </div>
              </div>
            </Tabs.Content>

            {/* 生成按钮 */}
            <div className="flex justify-center">
              <Button
                onClick={handleGenerate}
                disabled={isLoading || !prompt.trim()}
                size="lg"
                className="px-8"
              >
                {isLoading ? (
                  <>
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                    生成中...
                  </>
                ) : (
                  <>
                    <Wand2 className="h-4 w-4 mr-2" />
                    开始生成
                  </>
                )}
              </Button>
            </div>
          </div>

          {/* 右侧参数设置 */}
          <div className="space-y-6">
            <div className="bg-white rounded-lg border p-6">
              <div className="flex items-center gap-2 mb-4">
                <Settings className="h-5 w-5" />
                <h3 className="font-semibold">生成设置</h3>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    图片尺寸
                  </label>
                  <Select.Root value={size} onValueChange={setSize}>
                    <Select.Trigger className="w-full flex items-center justify-between px-3 py-2 border rounded-md">
                      <Select.Value />
                      <Select.Icon />
                    </Select.Trigger>
                    <Select.Content className="bg-white border rounded-md shadow-lg">
                      <Select.Item value="512x512" className="px-3 py-2 hover:bg-gray-50">
                        512x512
                      </Select.Item>
                      <Select.Item value="1024x1024" className="px-3 py-2 hover:bg-gray-50">
                        1024x1024 (推荐)
                      </Select.Item>
                      <Select.Item value="1024x1792" className="px-3 py-2 hover:bg-gray-50">
                        1024x1792 (竖屏)
                      </Select.Item>
                      <Select.Item value="1792x1024" className="px-3 py-2 hover:bg-gray-50">
                        1792x1024 (横屏)
                      </Select.Item>
                    </Select.Content>
                  </Select.Root>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    图片质量
                  </label>
                  <Select.Root value={quality} onValueChange={setQuality}>
                    <Select.Trigger className="w-full flex items-center justify-between px-3 py-2 border rounded-md">
                      <Select.Value />
                      <Select.Icon />
                    </Select.Trigger>
                    <Select.Content className="bg-white border rounded-md shadow-lg">
                      <Select.Item value="standard" className="px-3 py-2 hover:bg-gray-50">
                        标准质量
                      </Select.Item>
                      <Select.Item value="hd" className="px-3 py-2 hover:bg-gray-50">
                        高清质量
                      </Select.Item>
                    </Select.Content>
                  </Select.Root>
                </div>
              </div>
            </div>

            {/* 使用说明 */}
            <div className="bg-blue-50 rounded-lg p-4">
              <h4 className="font-medium text-blue-900 mb-2">使用提示</h4>
              <ul className="text-sm text-blue-800 space-y-1">
                <li>• 详细描述能生成更准确的图片</li>
                <li>• 可以指定风格：写实、卡通、油画等</li>
                <li>• 描述光线和构图能提升效果</li>
                <li>• 图生图功能适合修改现有图片</li>
              </ul>
            </div>
          </div>
        </div>
      </Tabs.Root>
    </div>
  )
}