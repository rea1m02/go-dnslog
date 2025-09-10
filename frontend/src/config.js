// 前端配置文件
// 可以根据需要修改 API_BASE_URL 来指定后端地址

// 默认使用相对路径（前后端同域部署时）
// 或者指定完整 URL，例如: 'http://localhost:8080'
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

// 获取完整的 API 地址
export const getApiUrl = (path) => {
  if (API_BASE_URL) {
    return `${API_BASE_URL}${path}`
  }
  return path
}
