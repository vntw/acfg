import * as types from '../constants/action-types';
import fetchApi from '../api/client';
import { requestErrorHandler, promiseErrorHandler } from './error-handlers';

function requestLogs() {
  return {
    type: types.REQUEST_LOGS,
  };
}

function receiveLogs(logs) {
  return {
    type: types.RECEIVE_LOGS,
    logs,
  };
}

function requestLog(instanceUuid) {
  return {
    type: types.REQUEST_LOG,
    instanceUuid,
  };
}

export function fetchLog(instanceUuid) {
  return dispatch => {
    dispatch(requestLog(instanceUuid));

    return fetchApi(`/api/logs/${instanceUuid}`)
      .then(response =>
        response.json().then(json => {
          if (response.status !== 200) {
            requestErrorHandler(response, json);
            return;
          }

          dispatch(receiveLog(json));
        })
      );
  };
}

function receiveLog(log) {
  return {
    type: types.RECEIVE_LOG,
    log,
  };
}

function shouldFetchLogs(state) {
  const logs = state.logs;

  if (!logs.items) {
    return true;
  } else if (logs.isFetching) {
    return false;
  } else {
    return logs.didInvalidate;
  }
}

function fetchLogs() {
  return dispatch => {
    dispatch(requestLogs());

    return fetchApi('/api/logs')
      .then(response =>
        response.json().then(json => {
          if (response.status !== 200) {
            requestErrorHandler(response, json);
            return;
          }

          dispatch(receiveLogs(json));
        })
      ).catch(promiseErrorHandler);
  };
}

export function fetchLogsIfNeeded() {
  return (dispatch, getState) => {
    if (shouldFetchLogs(getState())) {
      return dispatch(fetchLogs());
    }
  };
}
