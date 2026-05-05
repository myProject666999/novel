import React, { useState } from 'react';
import { Form, Input, Button, message, Checkbox } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined, SafetyOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { authApi } from '../utils/api';

function RegisterPage() {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const onFinish = async (values) => {
    if (values.password !== values.confirmPassword) {
      message.error('两次输入的密码不一致');
      return;
    }

    setLoading(true);
    try {
      const registerData = {
        username: values.username,
        password: values.password,
        nickname: values.nickname,
        email: values.email,
        phone: values.phone,
        invite_code: values.inviteCode,
      };

      const res = await authApi.register(registerData);
      if (res.code === 200) {
        message.success('注册成功，请登录');
        navigate('/login');
      }
    } catch (error) {
      console.error('注册失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h2 className="auth-title">注册</h2>
        <Form
          name="register"
          className="auth-form"
          onFinish={onFinish}
          size="large"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名!' }]}
          >
            <Input 
              prefix={<UserOutlined className="site-form-item-icon" />} 
              placeholder="用户名" 
            />
          </Form.Item>

          <Form.Item
            name="nickname"
            rules={[{ required: true, message: '请输入昵称!' }]}
          >
            <Input 
              prefix={<UserOutlined className="site-form-item-icon" />} 
              placeholder="昵称" 
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[
              { required: true, message: '请输入密码!' },
              { min: 6, message: '密码至少6个字符!' }
            ]}
          >
            <Input.Password
              prefix={<LockOutlined className="site-form-item-icon" />}
              placeholder="密码"
            />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            rules={[{ required: true, message: '请确认密码!' }]}
          >
            <Input.Password
              prefix={<LockOutlined className="site-form-item-icon" />}
              placeholder="确认密码"
            />
          </Form.Item>

          <Form.Item
            name="email"
            rules={[
              { required: true, message: '请输入邮箱!' },
              { type: 'email', message: '请输入正确的邮箱格式!' }
            ]}
          >
            <Input 
              prefix={<MailOutlined className="site-form-item-icon" />} 
              placeholder="邮箱" 
            />
          </Form.Item>

          <Form.Item
            name="phone"
            rules={[
              { required: true, message: '请输入手机号!' },
              { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号!' }
            ]}
          >
            <Input 
              prefix={<PhoneOutlined className="site-form-item-icon" />} 
              placeholder="手机号" 
            />
          </Form.Item>

          <Form.Item
            name="inviteCode"
            help="作家注册需要邀请码，普通读者不需要"
          >
            <Input 
              prefix={<SafetyOutlined className="site-form-item-icon" />} 
              placeholder="邀请码（选填）" 
            />
          </Form.Item>

          <Form.Item
            name="agreement"
            valuePropName="checked"
            rules={[
              { validator: (_, value) =>
                value ? Promise.resolve() : Promise.reject(new Error('请阅读并同意用户协议'))
              },
            ]}
          >
            <Checkbox>
              我已阅读并同意 <a href="#">用户协议</a> 和 <a href="#">隐私政策</a>
            </Checkbox>
          </Form.Item>

          <Form.Item>
            <Button 
              type="primary" 
              htmlType="submit" 
              className="auth-form-button"
              loading={loading}
            >
              注册
            </Button>
          </Form.Item>
        </Form>
        <div className="auth-link">
          已有账号? <Link to="/login">立即登录</Link>
        </div>
      </div>
    </div>
  );
}

export default RegisterPage;
