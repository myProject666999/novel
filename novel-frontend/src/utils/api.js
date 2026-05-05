import axios from 'axios';
import { message } from 'antd';
import { getToken, clearAuth } from './auth';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    const { data } = response;
    if (data.code === 200) {
      return data;
    } else {
      message.error(data.message || '请求失败');
      return Promise.reject(new Error(data.message || '请求失败'));
    }
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response;
      switch (status) {
        case 401:
          message.error('登录已过期，请重新登录');
          clearAuth();
          window.location.href = '/login';
          break;
        case 403:
          message.error('权限不足');
          break;
        case 404:
          message.error('请求的资源不存在');
          break;
        case 500:
          message.error('服务器错误');
          break;
        default:
          message.error(data?.message || '请求失败');
      }
    } else {
      message.error('网络错误，请检查网络连接');
    }
    return Promise.reject(error);
  }
);

// 公共API方法
export const novelApi = {
  // 获取小说列表
  getList: (params) => api.get('/public/novels', { params }),
  // 获取小说详情
  getDetail: (id) => api.get(`/public/novels/${id}`),
  // 获取小说排行榜
  getRank: (type, limit) => api.get('/public/novels/rank', { params: { type, limit } }),
  // 获取推荐小说
  getRecommend: (limit) => api.get('/public/novels/recommend', { params: { limit } }),
  // 搜索小说
  search: (keyword, params) => api.get('/public/novels/search', { params: { keyword, ...params } }),
  // 获取章节详情
  getChapter: (id) => api.get(`/public/chapters/${id}`),
  // 获取分类列表
  getCategories: () => api.get('/public/categories'),
  // 获取小说章节列表
  getChapters: (novelId, params) => api.get(`/public/novels/${novelId}/chapters`, { params }),
  // 获取小说评论
  getComments: (novelId, params) => api.get(`/public/comments/${novelId}`, { params }),
};

export const authApi = {
  // 登录
  login: (data) => api.post('/public/login', data),
  // 注册
  register: (data) => api.post('/public/register', data),
};

export const userApi = {
  // 获取用户信息
  getProfile: () => api.get('/user/profile'),
  // 更新用户信息
  updateProfile: (data) => api.put('/user/profile', data),
  // 修改密码
  changePassword: (data) => api.put('/user/password', data),
};

export const readerApi = {
  // 获取书架
  getBookshelf: () => api.get('/reader/bookshelf'),
  // 添加到书架
  addToBookshelf: (data) => api.post('/reader/bookshelf', data),
  // 从书架移除
  removeFromBookshelf: (id) => api.delete(`/reader/bookshelf/${id}`),
  // 创建评论
  createComment: (data) => api.post('/reader/comments', data),
  // 充值
  recharge: (data) => api.post('/reader/recharge', data),
  // 订阅
  subscribe: (data) => api.post('/reader/subscribe', data),
  // 获取订单列表
  getOrders: (params) => api.get('/reader/orders', { params }),
};

export const writerApi = {
  // 获取我的小说列表
  getNovels: (params) => api.get('/writer/novels', { params }),
  // 创建小说
  createNovel: (data) => api.post('/writer/novels', data),
  // 更新小说
  updateNovel: (id, data) => api.put(`/writer/novels/${id}`, data),
  // 删除小说
  deleteNovel: (id) => api.delete(`/writer/novels/${id}`),
  // 获取小说章节列表
  getChapters: (novelId, params) => api.get(`/writer/novels/${novelId}/chapters`, { params }),
  // 创建章节
  createChapter: (novelId, data) => api.post(`/writer/novels/${novelId}/chapters`, data),
  // 更新章节
  updateChapter: (id, data) => api.put(`/writer/chapters/${id}`, data),
  // 删除章节
  deleteChapter: (id) => api.delete(`/writer/chapters/${id}`),
};

export default api;
