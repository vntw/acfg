const TOKEN_NAME = 'token';

export const getToken = () => {
  return localStorage.getItem(TOKEN_NAME);
};

export const setToken = (token) => {
  localStorage.setItem(TOKEN_NAME, token);
};

export const unsetToken = () => {
  localStorage.removeItem(TOKEN_NAME);
};
