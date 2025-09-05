import React, { useState } from 'react';
import './styles.css';

interface GeneratedImage {
  id: string;
  prompt_text: string;
  image_url: string;
  thumbnail_url: string;
  status: string;
  generation_time: number;
  created_at: string;
}

function App() {
  const [prompt, setPrompt] = useState('');
  const [loading, setLoading] = useState(false);
  const [images, setImages] = useState<GeneratedImage[]>([]);
  const [error, setError] = useState('');
  const [mode, setMode] = useState<'text2img' | 'img2img'>('text2img');
  const [sourceImage, setSourceImage] = useState<string>('');

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (event) => {
        setSourceImage(event.target?.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!prompt.trim()) return;

    setLoading(true);
    setError('');

    try {
      const endpoint = mode === 'text2img' ? '/api/v1/generate/text2img' : '/api/v1/generate/img2img';
      const requestBody = mode === 'text2img' 
        ? {
            prompt: prompt.trim(),
            count: 1,
            params: {
              size: '512x512',
              quality: 'standard'
            }
          }
        : {
            prompt: prompt.trim(),
            source_image: sourceImage,
            count: 1,
            params: {
              size: '512x512',
              quality: 'standard'
            }
          };

      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      if (data.success && data.data) {
        setImages(prev => [...prev, ...data.data]);
      } else {
        throw new Error(data.message || '生成失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '网络错误');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        <h1 className="text-3xl font-bold text-center mb-8 text-gray-900">
          Nano Banana Qwen - AI图片生成
        </h1>
        
        <div className="bg-white rounded-lg shadow-md p-6 mb-8">
          {/* 模式选择 */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-3">选择生成模式</label>
            <div className="flex space-x-4">
              <button
                type="button"
                onClick={() => setMode('text2img')}
                className={`px-4 py-2 rounded-md font-medium ${
                  mode === 'text2img' 
                    ? 'bg-blue-600 text-white' 
                    : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                }`}
              >
                文本生成图片
              </button>
              <button
                type="button"
                onClick={() => setMode('img2img')}
                className={`px-4 py-2 rounded-md font-medium ${
                  mode === 'img2img' 
                    ? 'bg-blue-600 text-white' 
                    : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                }`}
              >
                图片生成图片
              </button>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            {/* 图片上传 (只在img2img模式显示) */}
            {mode === 'img2img' && (
              <div>
                <label htmlFor="image" className="block text-sm font-medium text-gray-700 mb-2">
                  上传源图片
                </label>
                <input
                  type="file"
                  id="image"
                  accept="image/*"
                  onChange={handleImageUpload}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
                {sourceImage && (
                  <div className="mt-3">
                    <img 
                      src={sourceImage} 
                      alt="Source" 
                      className="w-32 h-32 object-cover rounded-md border"
                    />
                  </div>
                )}
              </div>
            )}

            <div>
              <label htmlFor="prompt" className="block text-sm font-medium text-gray-700 mb-2">
                {mode === 'text2img' ? '输入图片描述 (支持中文)' : '输入修改描述 (支持中文)'}
              </label>
              <textarea
                id="prompt"
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="例如：一只可爱的小猫咪，蓝色的眼睛"
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                disabled={loading}
              />
            </div>
            
            <button
              type="submit"
              disabled={loading || !prompt.trim() || (mode === 'img2img' && !sourceImage)}
              className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading 
                ? (mode === 'text2img' ? '生成中...' : '转换中...') 
                : (mode === 'text2img' ? '生成图片' : '转换图片')
              }
            </button>
          </form>
          
          {error && (
            <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-md text-red-700">
              错误: {error}
            </div>
          )}
        </div>

        {images.length > 0 && (
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold mb-4 text-gray-900">生成的图片</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {images.map((image) => (
                <div key={image.id} className="border border-gray-200 rounded-lg p-4">
                  <img
                    src={image.image_url}
                    alt={image.prompt_text}
                    className="w-full h-48 object-cover rounded-md mb-3"
                    onError={(e) => {
                      console.error('Image load error:', e);
                    }}
                  />
                  <p className="text-sm text-gray-600 mb-2">{image.prompt_text}</p>
                  <p className="text-xs text-gray-400">
                    生成时间: {image.generation_time.toFixed(2)}s
                  </p>
                </div>
              ))}
            </div>
          </div>
        )}

        <div className="mt-8 text-center">
          <p className="text-sm text-gray-500">
            基于 OpenRouter API 和 Google Gemini 2.5 Flash 模型
          </p>
        </div>
      </div>
    </div>
  );
}

export default App;