import * as types from '../constants/action-types';

const initialState = {
  isFetching: true,
  configs: null,
};

export default function configsReducer(state = initialState, action) {
  switch (action.type) {
    case types.REQUEST_SERVERS:
      return {
        ...state,
        isFetching: true,
      };

    case types.RECEIVE_UPLOADED_CONFIG:
      return {
        ...state,
        configs: [...state.configs, action.config],
      };

    case types.CONFIG_DELETED:
      return {
        ...state,
        configs: state.configs.filter(cfg => cfg.uuid !== action.id)
      };

    case types.RECEIVE_SERVERS:
      return {
        ...state,
        isFetching: false,
        configs: action.configs,
      };

    default:
      return state;
  }
}
