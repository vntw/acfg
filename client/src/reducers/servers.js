import Alert from 'react-s-alert';
import * as types from '../constants/action-types';

const initialState = {
  isFetching: true,
  servers: null,
};

export default function serversReducer(state = initialState, action) {
  switch (action.type) {
    case types.REQUEST_SERVERS:
      return {
        ...state,
        isFetching: true,
        didInvalidate: false
      };

    case types.RECEIVE_SERVERS:
      return {
        ...state,
        isFetching: false,
        // TODO: Remove?
        didInvalidate: false,
        servers: action.servers,
      };

    case types.SERVER_LOADING_START: {
      const index = state.servers.findIndex(srv => srv.instance.uuid === action.uuid);
      const srv = { ...state.servers[index], isLoading: true };
      srv.instance = { uuid: action.uuid };

      return {
        ...state,
        servers: state.servers.slice(0, index)
          .concat([srv])
          .concat(state.servers.slice(index + 1))
        ,
      };
    }

    case types.SERVER_LOADING_STOP: {
      const index = state.servers.findIndex(srv => srv.instance.uuid === action.uuid);
      const srv = { ...action.server, isLoading: false };

      return {
        ...state,
        servers: state.servers.slice(0, index)
          .concat([srv])
          .concat(state.servers.slice(index + 1))
        ,
      };
    }

    case types.SERVER_STARTED: {
      const srv = { ...action.server, isLoading: false };

      return {
        ...state,
        servers: state.servers.concat([srv]),
      };
    }

    case types.SERVER_STOPPED:
      return {
        ...state,
        servers: state.servers.filter(srv => srv.instance.uuid !== action.uuid)
      };

    case types.SERVER_OP_ERROR: {
      Alert.error(state.error);
      return state;
    }

    default:
      return state;
  }
}
