import { proxy } from 'valtio'
import { Prompt, Generation, BatchJob, Image } from '@/services/api'

interface AppState {
  // 当前页面
  currentPage: string
  
  // 加载状态
  loading: boolean
  
  // 错误信息
  error: string | null
  
  // 提示词数据
  prompts: {
    list: Prompt[]
    categories: string[]
    tags: string[]
    currentPrompt: Prompt | null
    total: number
    currentPage: number
    pageSize: number
  }
  
  // 生成记录
  generations: {
    list: Generation[]
    currentGeneration: Generation | null
    total: number
    currentPage: number
    pageSize: number
  }
  
  // 批量任务
  batchJobs: {
    list: BatchJob[]
    currentJob: BatchJob | null
    total: number
    currentPage: number
    pageSize: number
  }
  
  // 图片库
  images: {
    list: Image[]
    currentImage: Image | null
    total: number
    currentPage: number
    pageSize: number
  }
  
  // 生成参数
  generationSettings: {
    size: string
    quality: string
    strength: number
  }
  
  // UI状态
  ui: {
    sidebarOpen: boolean
    theme: 'light' | 'dark'
    showGenerationHistory: boolean
    showBatchQueue: boolean
  }
}

// 创建初始状态
const initialState: AppState = {
  currentPage: '/',
  loading: false,
  error: null,
  
  prompts: {
    list: [],
    categories: [],
    tags: [],
    currentPrompt: null,
    total: 0,
    currentPage: 1,
    pageSize: 20,
  },
  
  generations: {
    list: [],
    currentGeneration: null,
    total: 0,
    currentPage: 1,
    pageSize: 20,
  },
  
  batchJobs: {
    list: [],
    currentJob: null,
    total: 0,
    currentPage: 1,
    pageSize: 20,
  },
  
  images: {
    list: [],
    currentImage: null,
    total: 0,
    currentPage: 1,
    pageSize: 20,
  },
  
  generationSettings: {
    size: '1024x1024',
    quality: 'standard',
    strength: 0.8,
  },
  
  ui: {
    sidebarOpen: true,
    theme: 'light',
    showGenerationHistory: false,
    showBatchQueue: false,
  },
}

// 创建响应式store
export const appStore = proxy(initialState)

// Store操作
export const appActions = {
  // 通用操作
  setLoading: (loading: boolean) => {
    appStore.loading = loading
  },
  
  setError: (error: string | null) => {
    appStore.error = error
  },
  
  setCurrentPage: (page: string) => {
    appStore.currentPage = page
  },
  
  // 提示词操作
  setPrompts: (prompts: Prompt[], total: number, page: number) => {
    appStore.prompts.list = prompts
    appStore.prompts.total = total
    appStore.prompts.currentPage = page
  },
  
  addPrompt: (prompt: Prompt) => {
    appStore.prompts.list.unshift(prompt)
    appStore.prompts.total += 1
  },
  
  updatePrompt: (prompt: Prompt) => {
    const index = appStore.prompts.list.findIndex(p => p.id === prompt.id)
    if (index !== -1) {
      appStore.prompts.list[index] = prompt
    }
  },
  
  removePrompt: (id: string) => {
    appStore.prompts.list = appStore.prompts.list.filter(p => p.id !== id)
    appStore.prompts.total -= 1
  },
  
  setCurrentPrompt: (prompt: Prompt | null) => {
    appStore.prompts.currentPrompt = prompt
  },
  
  setCategories: (categories: string[]) => {
    appStore.prompts.categories = categories
  },
  
  setTags: (tags: string[]) => {
    appStore.prompts.tags = tags
  },
  
  // 生成记录操作
  setGenerations: (generations: Generation[], total: number, page: number) => {
    appStore.generations.list = generations
    appStore.generations.total = total
    appStore.generations.currentPage = page
  },
  
  addGeneration: (generation: Generation) => {
    appStore.generations.list.unshift(generation)
    appStore.generations.total += 1
  },
  
  updateGeneration: (generation: Generation) => {
    const index = appStore.generations.list.findIndex(g => g.id === generation.id)
    if (index !== -1) {
      appStore.generations.list[index] = generation
    }
  },
  
  removeGeneration: (id: string) => {
    appStore.generations.list = appStore.generations.list.filter(g => g.id !== id)
    appStore.generations.total -= 1
  },
  
  setCurrentGeneration: (generation: Generation | null) => {
    appStore.generations.currentGeneration = generation
  },
  
  // 批量任务操作
  setBatchJobs: (jobs: BatchJob[], total: number, page: number) => {
    appStore.batchJobs.list = jobs
    appStore.batchJobs.total = total
    appStore.batchJobs.currentPage = page
  },
  
  addBatchJob: (job: BatchJob) => {
    appStore.batchJobs.list.unshift(job)
    appStore.batchJobs.total += 1
  },
  
  updateBatchJob: (job: BatchJob) => {
    const index = appStore.batchJobs.list.findIndex(j => j.id === job.id)
    if (index !== -1) {
      appStore.batchJobs.list[index] = job
    }
  },
  
  removeBatchJob: (id: string) => {
    appStore.batchJobs.list = appStore.batchJobs.list.filter(j => j.id !== id)
    appStore.batchJobs.total -= 1
  },
  
  setCurrentBatchJob: (job: BatchJob | null) => {
    appStore.batchJobs.currentJob = job
  },
  
  // 图片操作
  setImages: (images: Image[], total: number, page: number) => {
    appStore.images.list = images
    appStore.images.total = total
    appStore.images.currentPage = page
  },
  
  addImage: (image: Image) => {
    appStore.images.list.unshift(image)
    appStore.images.total += 1
  },
  
  removeImage: (id: string) => {
    appStore.images.list = appStore.images.list.filter(i => i.id !== id)
    appStore.images.total -= 1
  },
  
  setCurrentImage: (image: Image | null) => {
    appStore.images.currentImage = image
  },
  
  // 生成设置操作
  updateGenerationSettings: (settings: Partial<typeof initialState.generationSettings>) => {
    Object.assign(appStore.generationSettings, settings)
  },
  
  // UI操作
  toggleSidebar: () => {
    appStore.ui.sidebarOpen = !appStore.ui.sidebarOpen
  },
  
  setSidebarOpen: (open: boolean) => {
    appStore.ui.sidebarOpen = open
  },
  
  setTheme: (theme: 'light' | 'dark') => {
    appStore.ui.theme = theme
  },
  
  toggleGenerationHistory: () => {
    appStore.ui.showGenerationHistory = !appStore.ui.showGenerationHistory
  },
  
  toggleBatchQueue: () => {
    appStore.ui.showBatchQueue = !appStore.ui.showBatchQueue
  },
}

export default appStore