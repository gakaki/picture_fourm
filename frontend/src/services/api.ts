import axios from 'axios'

// API基础配置
const API_BASE_URL = 'http://localhost:8080/api/v1'

// 创建axios实例
export const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 可以在这里添加认证token等
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    // 统一错误处理
    console.error('API Error:', error)
    throw error
  }
)

// 类型定义
export interface APIResponse<T = any> {
  success: boolean
  message: string
  data?: T
  error?: string
}

export interface GenerationParams {
  model?: string
  size?: string
  quality?: string
  strength?: number
}

export interface Text2ImgRequest {
  prompt: string
  count?: number
  params?: GenerationParams
}

export interface Img2ImgRequest {
  prompt: string
  source_image: string
  count?: number
  params?: GenerationParams
}

export interface Prompt {
  id: string
  title: string
  content: string
  category: string
  tags: string[]
  is_favorite: boolean
  usage_count: number
  created_at: string
  updated_at: string
}

export interface Generation {
  id: string
  prompt_id?: string
  prompt_text: string
  image_url: string
  thumbnail_url: string
  generation_params: GenerationParams
  status: string
  error_message?: string
  generation_time: number
  batch_job_id?: string
  is_img2img: boolean
  source_image_id?: string
  created_at: string
}

export interface BatchJob {
  id: string
  name: string
  prompts: BatchPrompt[]
  total_images: number
  completed_images: number
  failed_images: number
  status: string
  started_at?: string
  completed_at?: string
  created_at: string
}

export interface BatchPrompt {
  prompt_id?: string
  prompt_text: string
  count: number
  completed: number
  failed: number
}

export interface Image {
  id: string
  filename: string
  original_filename: string
  file_path: string
  thumbnail_path: string
  file_size: number
  width: number
  height: number
  format: string
  generation_id?: string
  prompt_text: string
  is_img2img: boolean
  source_image_id?: string
  created_at: string
}

// API函数
export const apiService = {
  // 健康检查
  health: () => api.get('/health'),

  // 提示词管理
  prompts: {
    list: (params?: { page?: number; page_size?: number; keyword?: string; category?: string }) =>
      api.get('/prompts', { params }) as Promise<APIResponse<{ prompts: Prompt[]; total: number; page: number; page_size: number; total_pages: number }>>,
    
    get: (id: string) =>
      api.get(`/prompts/${id}`) as Promise<APIResponse<Prompt>>,
    
    create: (data: Omit<Prompt, 'id' | 'created_at' | 'updated_at' | 'usage_count'>) =>
      api.post('/prompts', data) as Promise<APIResponse<Prompt>>,
    
    update: (id: string, data: Partial<Prompt>) =>
      api.put(`/prompts/${id}`, data) as Promise<APIResponse<Prompt>>,
    
    delete: (id: string) =>
      api.delete(`/prompts/${id}`) as Promise<APIResponse<null>>,
    
    getCategories: () =>
      api.get('/prompts/categories') as Promise<APIResponse<string[]>>,
    
    getTags: () =>
      api.get('/prompts/tags') as Promise<APIResponse<string[]>>,
  },

  // 图片生成
  generate: {
    text2img: (data: Text2ImgRequest) =>
      api.post('/generate/text2img', data) as Promise<APIResponse<Generation[]>>,
    
    img2img: (data: Img2ImgRequest) =>
      api.post('/generate/img2img', data) as Promise<APIResponse<Generation[]>>,
  },

  // 生成记录
  generations: {
    list: (params?: { page?: number; page_size?: number; prompt?: string; status?: string; date_from?: string; date_to?: string; is_img2img?: boolean }) =>
      api.get('/generations', { params }) as Promise<APIResponse<{ generations: Generation[]; total: number; page: number; page_size: number; total_pages: number }>>,
    
    get: (id: string) =>
      api.get(`/generations/${id}`) as Promise<APIResponse<Generation>>,
    
    delete: (id: string) =>
      api.delete(`/generations/${id}`) as Promise<APIResponse<null>>,
  },

  // 批量任务
  batch: {
    create: (data: { name?: string; prompts: BatchPrompt[] }) =>
      api.post('/batch', data) as Promise<APIResponse<BatchJob>>,
    
    list: (params?: { page?: number; page_size?: number; status?: string }) =>
      api.get('/batch', { params }) as Promise<APIResponse<{ jobs: BatchJob[]; total: number; page: number; page_size: number; total_pages: number }>>,
    
    get: (id: string) =>
      api.get(`/batch/${id}`) as Promise<APIResponse<BatchJob>>,
    
    getStatus: (id: string) =>
      api.get(`/batch/${id}/status`) as Promise<APIResponse<{ job_id: string; status: string; total_images: number; completed_images: number; failed_images: number; progress: number; message: string; updated_at: string }>>,
    
    cancel: (id: string) =>
      api.delete(`/batch/${id}/cancel`) as Promise<APIResponse<null>>,
    
    delete: (id: string) =>
      api.delete(`/batch/${id}`) as Promise<APIResponse<null>>,
  },

  // 图片管理
  images: {
    list: (params?: { page?: number; page_size?: number; prompt?: string }) =>
      api.get('/images', { params }) as Promise<APIResponse<{ images: Image[]; total: number; page: number; page_size: number; total_pages: number }>>,
    
    get: (id: string) =>
      api.get(`/images/${id}`) as Promise<APIResponse<Image>>,
    
    delete: (id: string) =>
      api.delete(`/images/${id}`) as Promise<APIResponse<null>>,
    
    download: (id: string) =>
      `${API_BASE_URL}/images/${id}/download`,
    
    getUrl: (path: string) =>
      `${API_BASE_URL}/files${path}`,
    
    getThumbnailUrl: (path: string) =>
      `${API_BASE_URL}/files${path}`,
  },
}

export default apiService