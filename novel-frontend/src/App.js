import React, { useState, useEffect } from 'react';
import { Routes, Route, Navigate, useLocation, useNavigate } from 'react-router-dom';
import { Layout, Menu, Button, Avatar, Dropdown, message } from 'antd';
import {
  BookOutlined,
  HomeOutlined,
  TrophyOutlined,
  SearchOutlined,
  UserOutlined,
  BookFilled,
  SettingOutlined,
  LogoutOutlined,
  EditOutlined
} from '@ant-design/icons';
import './App.css';
import HomePage from './pages/HomePage';
import NovelDetailPage from './pages/NovelDetailPage';
import NovelListPage from './pages/NovelListPage';
import SearchPage from './pages/SearchPage';
import RankPage from './pages/RankPage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import BookshelfPage from './pages/reader/BookshelfPage';
import ProfilePage from './pages/reader/ProfilePage';
import ChangePasswordPage from './pages/reader/ChangePasswordPage';
import OrdersPage from './pages/reader/OrdersPage';
import WriterNovelListPage from './pages/writer/NovelListPage';
import WriterChapterListPage from './pages/writer/ChapterListPage';
import ChapterReadPage from './pages/reader/ChapterReadPage';
import { getCurrentUser, clearAuth, isLoggedIn, isWriter, isAdmin } from './utils/auth';

const { Header, Content, Footer, Sider } = Layout;

function App() {
  const [currentUser, setCurrentUser] = useState(null);
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    const user = getCurrentUser();
    setCurrentUser(user);
  }, [location]);

  const handleLogout = () => {
    clearAuth();
    setCurrentUser(null);
    message.success('退出登录成功');
    navigate('/');
  };

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/user/profile')
    },
    {
      key: 'bookshelf',
      icon: <BookFilled />,
      label: '我的书架',
      onClick: () => navigate('/reader/bookshelf')
    },
    {
      key: 'orders',
      icon: <BookOutlined />,
      label: '我的订单',
      onClick: () => navigate('/reader/orders')
    },
    {
      key: 'password',
      icon: <SettingOutlined />,
      label: '修改密码',
      onClick: () => navigate('/user/password')
    },
    {
      type: 'divider'
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout
    }
  ];

  // 如果是作家，添加作家菜单
  if (currentUser && isWriter()) {
    userMenuItems.splice(3, 0, {
      key: 'writer',
      icon: <EditOutlined />,
      label: '作家中心',
      onClick: () => navigate('/writer/novels')
    });
  }

  // 如果是管理员，提示前往管理后台
  if (currentUser && isAdmin()) {
    userMenuItems.splice(3, 0, {
      key: 'admin',
      icon: <SettingOutlined />,
      label: '管理后台',
      onClick: () => window.open('http://localhost:3001', '_blank')
    });
  }

  const mainMenuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: '首页',
      onClick: () => navigate('/')
    },
    {
      key: '/novels',
      icon: <BookOutlined />,
      label: '小说分类',
      onClick: () => navigate('/novels')
    },
    {
      key: '/rank',
      icon: <TrophyOutlined />,
      label: '排行榜',
      onClick: () => navigate('/rank')
    },
    {
      key: '/search',
      icon: <SearchOutlined />,
      label: '搜索',
      onClick: () => navigate('/search')
    }
  ];

  const selectedKey = location.pathname === '/' ? '/' : 
    mainMenuItems.find(item => location.pathname.startsWith(item.key))?.key || '/';

  // 登录/注册页面不需要布局
  if (location.pathname === '/login' || location.pathname === '/register') {
    return (
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
      </Routes>
    );
  }

  // 认证路由保护
  const ProtectedRoute = ({ children, requireRole = null }) => {
    if (!isLoggedIn()) {
      return <Navigate to="/login" replace />;
    }
    if (requireRole === 'writer' && !isWriter() && !isAdmin()) {
      message.error('您没有作家权限');
      return <Navigate to="/" replace />;
    }
    return children;
  };

  return (
    <Layout className="app-layout">
      <Header className="app-header">
        <div className="logo">
          <BookOutlined />
          <span>小说阅读网</span>
        </div>
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[selectedKey]}
          items={mainMenuItems}
          className="main-menu"
        />
        <div className="header-right">
          {currentUser ? (
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <div className="user-info">
                <Avatar icon={<UserOutlined />} />
                <span className="username">{currentUser.nickname || currentUser.username}</span>
                {isWriter() && <span className="role-tag writer-tag">作家</span>}
                {isAdmin() && <span className="role-tag admin-tag">管理员</span>}
              </div>
            </Dropdown>
          ) : (
            <div className="auth-buttons">
              <Button type="link" onClick={() => navigate('/login')}>登录</Button>
              <Button type="primary" onClick={() => navigate('/register')}>注册</Button>
            </div>
          )}
        </div>
      </Header>
      <Layout>
        <Content className="app-content">
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/novels" element={<NovelListPage />} />
            <Route path="/novel/:id" element={<NovelDetailPage />} />
            <Route path="/chapter/:id" element={<ChapterReadPage />} />
            <Route path="/rank" element={<RankPage />} />
            <Route path="/search" element={<SearchPage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            
            {/* 需要登录的路由 */}
            <Route path="/user/profile" element={
              <ProtectedRoute><ProfilePage /></ProtectedRoute>
            } />
            <Route path="/user/password" element={
              <ProtectedRoute><ChangePasswordPage /></ProtectedRoute>
            } />
            <Route path="/reader/bookshelf" element={
              <ProtectedRoute><BookshelfPage /></ProtectedRoute>
            } />
            <Route path="/reader/orders" element={
              <ProtectedRoute><OrdersPage /></ProtectedRoute>
            } />
            
            {/* 作家功能路由 */}
            <Route path="/writer/novels" element={
              <ProtectedRoute requireRole="writer"><WriterNovelListPage /></ProtectedRoute>
            } />
            <Route path="/writer/novel/:novelId/chapters" element={
              <ProtectedRoute requireRole="writer"><WriterChapterListPage /></ProtectedRoute>
            } />
          </Routes>
        </Content>
      </Layout>
      <Footer className="app-footer">
        <p>小说阅读网 ©2024 Created by Novel Team</p>
        <p>
          <a href="#">关于我们</a> | 
          <a href="#">联系方式</a> | 
          <a href="#">帮助中心</a> | 
          <a href="#">友情链接</a>
        </p>
      </Footer>
    </Layout>
  );
}

export default App;
