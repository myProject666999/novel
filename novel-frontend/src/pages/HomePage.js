import React, { useState, useEffect } from 'react';
import { Row, Col, Card, List, Carousel, Tag, Spin, Empty } from 'antd';
import { EyeOutlined, LikeOutlined, CommentOutlined, FireOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { novelApi } from '../utils/api';

const { Meta } = Card;

function HomePage() {
  const [loading, setLoading] = useState(true);
  const [recommendNovels, setRecommendNovels] = useState([]);
  const [hotNovels, setHotNovels] = useState([]);
  const [newNovels, setNewNovels] = useState([]);
  const [categories, setCategories] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      // 获取推荐小说
      const recommendRes = await novelApi.getRecommend(8);
      if (recommendRes.code === 200) {
        setRecommendNovels(recommendRes.data.list || []);
      }

      // 获取热门小说（点击榜）
      const hotRes = await novelApi.getRank('click', 10);
      if (hotRes.code === 200) {
        setHotNovels(hotRes.data.list || []);
      }

      // 获取最新小说
      const newRes = await novelApi.getList({ page_size: 8 });
      if (newRes.code === 200) {
        setNewNovels(newRes.data.list || []);
      }

      // 获取分类
      const categoryRes = await novelApi.getCategories();
      if (categoryRes.code === 200) {
        setCategories(categoryRes.data.list || []);
      }
    } catch (error) {
      console.error('获取首页数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleNovelClick = (novelId) => {
    navigate(`/novel/${novelId}`);
  };

  const handleCategoryClick = (categoryId) => {
    navigate(`/novels?category=${categoryId}`);
  };

  const carouselItems = [
    {
      id: 1,
      title: '热门推荐',
      desc: '精选优质小说，让阅读更精彩',
      image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fantasy%20book%20cover%20with%20dragons%20and%20magic&image_size=landscape_16_9'
    },
    {
      id: 2,
      title: '新书速递',
      desc: '每日更新，最新章节抢先看',
      image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=modern%20library%20with%20digital%20books%20reading%20app&image_size=landscape_16_9'
    },
    {
      id: 3,
      title: '作家专区',
      desc: '成为作家，分享你的故事',
      image: 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=writer%20typing%20on%20laptop%20creative%20writing&image_size=landscape_16_9'
    }
  ];

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px 0' }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      {/* 轮播图 */}
      <Carousel autoplay style={{ marginBottom: 24, borderRadius: 8, overflow: 'hidden' }}>
        {carouselItems.map(item => (
          <div key={item.id} style={{ position: 'relative', height: 300 }}>
            <img 
              src={item.image} 
              alt={item.title} 
              style={{ width: '100%', height: 300, objectFit: 'cover' }}
            />
            <div style={{
              position: 'absolute',
              bottom: 0,
              left: 0,
              right: 0,
              background: 'linear-gradient(transparent, rgba(0,0,0,0.7))',
              padding: '40px 50px',
              color: 'white'
            }}>
              <h2 style={{ color: 'white', fontSize: 28, marginBottom: 8 }}>{item.title}</h2>
              <p style={{ fontSize: 16, opacity: 0.9 }}>{item.desc}</p>
            </div>
          </div>
        ))}
      </Carousel>

      {/* 分类导航 */}
      <Card style={{ marginBottom: 24 }}>
        <h3 style={{ marginBottom: 16, fontSize: 18, fontWeight: 'bold' }}>小说分类</h3>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 12 }}>
          {categories.map(cat => (
            <Tag 
              key={cat.id}
              color="blue"
              style={{ fontSize: 14, padding: '8px 16px', cursor: 'pointer' }}
              onClick={() => handleCategoryClick(cat.id)}
            >
              {cat.name}
            </Tag>
          ))}
        </div>
      </Card>

      {/* 推荐小说 */}
      {recommendNovels.length > 0 && (
        <div style={{ marginBottom: 24 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
            <h3 style={{ fontSize: 18, fontWeight: 'bold', margin: 0 }}>
              <FireOutlined style={{ color: '#ff4d4f', marginRight: 8 }} />
              推荐小说
            </h3>
          </div>
          <Row gutter={[16, 16]}>
            {recommendNovels.slice(0, 8).map(item => {
              const novel = item.novel || item;
              return (
                <Col xs={12} sm={8} md={6} lg={4} key={novel.id}>
                  <Card
                    hoverable
                    className="novel-card"
                    cover={
                      <img
                        alt={novel.title}
                        src={novel.cover || `https://picsum.photos/200/280?random=${novel.id}`}
                      />
                    }
                    onClick={() => handleNovelClick(novel.id)}
                  >
                    <Meta
                      title={<div className="novel-card-title">{novel.title}</div>}
                      description={
                        <div>
                          <div className="novel-card-info">
                            <span>{novel.author?.nickname || novel.author?.username || '未知作者'}</span>
                            <span><EyeOutlined /> {novel.click_count || 0}</span>
                          </div>
                          <div className="novel-card-desc">{novel.description || '暂无简介'}</div>
                        </div>
                      }
                    />
                  </Card>
                </Col>
              );
            })}
          </Row>
        </div>
      )}

      {/* 热门榜单 + 最新小说 */}
      <Row gutter={[24, 24]}>
        {/* 热门榜单 */}
        <Col xs={24} lg={12}>
          <Card title="热门榜单">
            {hotNovels.length > 0 ? (
              <List
                dataSource={hotNovels}
                renderItem={(novel, index) => (
                  <List.Item 
                    key={novel.id}
                    style={{ cursor: 'pointer', padding: '12px 0' }}
                    onClick={() => handleNovelClick(novel.id)}
                  >
                    <List.Item.Meta
                      avatar={
                        <span style={{
                          display: 'inline-flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          width: 24,
                          height: 24,
                          borderRadius: 4,
                          background: index < 3 ? '#ff4d4f' : '#d9d9d9',
                          color: 'white',
                          fontWeight: 'bold',
                          fontSize: 12
                        }}>
                          {index + 1}
                        </span>
                      }
                      title={
                        <span style={{ fontWeight: index < 3 ? 'bold' : 'normal' }}>
                          {novel.title}
                        </span>
                      }
                      description={
                        <span style={{ color: '#999', fontSize: 12 }}>
                          <EyeOutlined style={{ marginRight: 4 }} />
                          {novel.click_count || 0}
                          <LikeOutlined style={{ marginLeft: 16, marginRight: 4 }} />
                          {novel.collect_count || 0}
                        </span>
                      }
                    />
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无数据" />
            )}
          </Card>
        </Col>

        {/* 最新小说 */}
        <Col xs={24} lg={12}>
          <Card title="最新小说">
            {newNovels.length > 0 ? (
              <List
                dataSource={newNovels}
                renderItem={(novel) => (
                  <List.Item 
                    key={novel.id}
                    style={{ cursor: 'pointer', padding: '12px 0' }}
                    onClick={() => handleNovelClick(novel.id)}
                  >
                    <List.Item.Meta
                      avatar={
                        <img
                          alt={novel.title}
                          src={novel.cover || `https://picsum.photos/60/80?random=${novel.id}`}
                          style={{ width: 50, height: 70, objectFit: 'cover', borderRadius: 4 }}
                        />
                      }
                      title={<span style={{ fontWeight: 'bold' }}>{novel.title}</span>}
                      description={
                        <div style={{ color: '#999', fontSize: 12 }}>
                          <p style={{ margin: 0 }}>{novel.category?.name || '未分类'}</p>
                          <p style={{ margin: 0 }}>
                            {novel.author?.nickname || novel.author?.username || '未知作者'}
                            <span style={{ marginLeft: 16 }}>
                              <CommentOutlined style={{ marginRight: 4 }} />
                              {novel.comment_count || 0}
                            </span>
                          </p>
                        </div>
                      }
                    />
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无数据" />
            )}
          </Card>
        </Col>
      </Row>
    </div>
  );
}

export default HomePage;
