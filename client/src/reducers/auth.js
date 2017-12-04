import * as types from '../constants/action-types';

const hasToken = !!localStorage.getItem('token');

const initialState = {
  isFetching: !hasToken,
  isAuthenticated: hasToken,
  user: null,
  errorMessage: null,
};

export default function authReducer(state = initialState, action) {
  switch (action.type) {
    case types.LOGIN_REQUEST:
      return {
        ...state,
        isFetching: true,
        isAuthenticated: false,
        user: null,
        errorMessage: '',
      };
    case types.LOGIN_SUCCESS:
      return {
        ...state,
        isFetching: false,
        isAuthenticated: true,
        user: action.token,
        errorMessage: '',
      };
    case types.LOGIN_FAILURE:
      return {
        ...state,
        isFetching: false,
        isAuthenticated: false,
        user: null,
        errorMessage: action.message,
      };
    case types.LOGOUT_SUCCESS:
      return {
        ...state,
        isFetching: true,
        isAuthenticated: false,
      };
    default:
      return state;
  }
}
