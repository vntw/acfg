import * as types from '../constants/action-types';

const initialState = {
  isFetching: true,
  logs: null,
};

export default function logsReducer(state = initialState, action) {
  switch (action.type) {
    case types.REQUEST_LOGS:
      return {
        ...state,
        isFetching: true,
        didInvalidate: false
      };

    case types.RECEIVE_LOGS:
      return {
        ...state,
        isFetching: false,
        // TODO: Remove?
        didInvalidate: false,
        logs: action.logs,
      };

    case types.REQUEST_LOG:
      return state;

    case types.RECEIVE_LOG: {
      const index = state.logs.findIndex(log => log.instanceUuid === action.log.instanceUuid);
      const log = { ...state.logs[index], ...action.log, isLoaded: true };

      return {
        ...state,
        logs: state.logs.slice(0, index)
          .concat([log])
          .concat(state.logs.slice(index + 1))
        ,
      };
    }

    default:
      return state;
  }
}
