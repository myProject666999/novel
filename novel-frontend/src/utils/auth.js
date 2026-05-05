const TOKEN_KEY = 'novel_token';
const USER_KEY = 'novel_user';

export function setAuth(token, user) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export function getCurrentUser() {
  const userStr = localStorage.getItem(USER_KEY);
  if (userStr) {
    try {
      return JSON.parse(userStr);
    } catch (e) {
      return null;
    }
  }
  return null;
}

export function clearAuth() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
}

export function isLoggedIn() {
  return !!getToken();
}

export function isWriter() {
  const user = getCurrentUser();
  // 假设角色ID 2是作家
  return user && (user.role_id === 2 || user.role_name === 'writer' || user.role_name === '作家');
}

export function isAdmin() {
  const user = getCurrentUser();
  // 假设角色ID 1是管理员
  return user && (user.role_id === 1 || user.role_name === 'admin' || user.role_name === '管理员');
}
