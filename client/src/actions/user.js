import * as types from '../constants/action-types';
import { setToken, unsetToken } from '../api/auth';
import fetchApi from '../api/client';

function requestLogin(creds) {
  return {
    type: types.LOGIN_REQUEST,
    isFetching: true,
    isAuthenticated: false,
    creds,
  };
}

function receiveLogin(user) {
  return {
    type: types.LOGIN_SUCCESS,
    isFetching: false,
    isAuthenticated: true,
    token: user.token,
  };
}

function loginError(message) {
  return {
    type: types.LOGIN_FAILURE,
    isFetching: false,
    isAuthenticated: false,
    message,
  };
}


function requestLogout() {
  return {
    type: types.LOGOUT_REQUEST,
    isFetching: true,
    isAuthenticated: true,
  };
}

function receiveLogout() {
  return {
    type: types.LOGOUT_SUCCESS,
    isFetching: false,
    isAuthenticated: false,
  };
}

export function loginUser(creds) {
  let config = {
    method: 'POST',
    headers: { 'Content-Type':'application/x-www-form-urlencoded' },
    body: `username=${creds.username}&password=${creds.password}`
  };

  return dispatch => {
    dispatch(requestLogin(creds));
    return fetchApi('/api/sessions', config, false)
      .then(response =>
        response.json()
          .then(user => ({ user, response }))
      ).then(({ user, response }) =>  {
        if (!response.ok) {
          dispatch(loginError(user.message));
          return Promise.reject(user);
        } else {
          setToken(user.token);
          dispatch(receiveLogin(user));
        }
      });
  };
}

export function logoutUser() {
  return dispatch => {
    dispatch(requestLogout());
    unsetToken();
    dispatch(receiveLogout());
  };
}
