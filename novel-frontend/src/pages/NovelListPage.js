import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Select, Pagination, Spin, Empty, Tag } from 'antd';
import { EyeOutlined, LikeOutlined } from '@ant-design/icons';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { novelApi } from '../utils/api';

const { Meta } = Card;
const { Option } = Select;

function NovelListPage() {
  const [loading, setLoading] = useState(false);
  const [novels, setNovels] = useState([]);
  const [categories, setCategories] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(12);
  const [total, setTotal] = useState(0);
  const [selectedCategory, setSelectedCategory] = useState(undefined);
  const [selectedStatus, setSelectedStatus] = useState(undefined);
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    fetchCategories();
  }, []);

  useEffect(() => {
    const category = searchParams.get('category');
    if (category) {
      setSelectedCategory(category);
    }
    fetchNovels();
  }, [currentPage, pageSize, selectedCategory, selectedStatus]);

  const fetchCategories = async () => {
    try {
      const res = await novelApi.getCategories();
      if (res.code === 200) {
        setCategories(res.data.list || []);
      }
    } catch (error) {
      console.error('获取分类失败:', error);
    }
  };

  const fetchNovels = async () => {
    setLoading(true);
    try {
      const params = {
        page: currentPage,
        page_size: pageSize,
      };
      if (selectedCategory) {
        params.category_id = selectedCategory;
      }
      if (selectedStatus) {
        params.status = selectedStatus;
      }

      const res = await novelApi.getList(params);
      if (res.code === 200) {
        setNovels(res.data.list || []);
        setTotal(res.data.total || 0);
      }
    } catch (error) {
      console.error('获取小说列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleNovelClick = (novelId) => {
    navigate(`/novel/${novelId}`);
  };

  const handleCategoryChange = (value) => {
    setSelectedCategory(value);
    setCurrentPage(1);
  };

  const handleStatusChange = (value) => {
    setSelectedStatus(value);
    setCurrentPage(1);
  };

  const handlePageChange = (page, size) => {
    setCurrentPage(page);
    if (size) setPageSize(size);
  };

  const statusMap = {
    1: { text: '连载中', color: 'blue' },
    2: { text: '已完结', color: 'green' },
    3: { text: '下架', color: 'red' }
  };

  return (
    <div className="page-container">
      <div className="page-title">小说分类</div>

      {/* 筛选条件 */}
      <div style={{ marginBottom: 24, display: 'flex', gap: 16, alignItems: 'center' }}>
        <span>分类：</span>
        <Select
          style={{ width: 200 }}
          placeholder="全部分类"
          allowClear
          value={selectedCategory}
          onChange={handleCategoryChange}
        >
          {categories.map(cat => (
            <Option key={cat.id} value={cat.id}>{cat.name}</Option>
          ))}
        </Select>

        <span>状态：</span>
        <Select
          style={{ width: 150 }}
          placeholder="全部状态"
          allowClear
          value={selectedStatus}
          onChange={handleStatusChange}
        >
          <Option value={1}>连载中</Option>
          <Option value={2}>已完结</Option>
        </Select>
      </div>

      {/* 小说列表 */}
      {loading ? (
        <div style={{ textAlign: 'center', padding: '50px 0' }}>
          <Spin size="large" />
        </div>
      ) : novels.length > 0 ? (
        <>
          <Row gutter={[16, 16]}>
            {novels.map(novel => (
              <Col xs={12} sm={8} md={6} lg={4} key={novel.id}>
                <Card
                  hoverable
                  className="novel-card"
                  cover={
                    <img
                      alt={novel.title}
                      src={novel.cover || `https://picsum.photos/200/280?random=${novel.id}`}
                      style={{ height: 180, objectFit: 'cover' }}
                    />
                  }
                  onClick={() => handleNovelClick(novel.id)}
                >
                  <Meta
                    title={
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <span className="novel-card-title" style={{ margin: 0 }}>{novel.title}</span>
                        {statusMap[novel.status] && (
                          <Tag color={statusMap[novel.status].color}>
                            {statusMap[novel.status].text}
                          </Tag>
                        )}
                      </div>
                    }
                    description={
                      <div>
                        <div className="novel-card-info">
                          <span>{novel.author?.nickname || novel.author?.username || '未知作者'}</span>
                          <span>{novel.category?.name || '未分类'}</span>
                        </div>
                        <div className="novel-card-info">
                          <span><EyeOutlined /> {novel.click_count || 0}</span>
                          <span><LikeOutlined /> {novel.collect_count || 0}</span>
                        </div>
                        <div className="novel-card-desc">{novel.description || '暂无简介'}</div>
                      </div>
                    }
                  />
                </Card>
              </Col>
            ))}
          </Row>

          {/* 分页 */}
          <div style={{ textAlign: 'center', marginTop: 32 }}>
            <Pagination
              current={currentPage}
              pageSize={pageSize}
              total={total}
              showSizeChanger
              showQuickJumper
              showTotal={(total) => `共 ${total} 条`}
              onChange={handlePageChange}
            />
          </div>
        </>
      ) : (
        <Empty description="暂无小说" style={{ padding: '50px 0' }} />
      )}
    </div>
  );
}

export default NovelListPage;
