import * as types from '../constants/action-types';
import fetchApi from '../api/client';
import { requestErrorHandler, promiseErrorHandler } from './error-handlers';

function requestServers() {
  return {
    type: types.REQUEST_SERVERS,
  };
}

function serverLoadingStart(uuid) {
  return {
    type: types.SERVER_LOADING_START,
    uuid,
  };
}

function serverLoadingStop(uuid, server) {
  return {
    type: types.SERVER_LOADING_STOP,
    uuid,
    server,
  };
}

function serverStarted(server) {
  return {
    type: types.SERVER_STARTED,
    server,
  };
}

function serverStopped(uuid) {
  return {
    type: types.SERVER_STOPPED,
    uuid,
  };
}

function serverLoadingError(uuid, error) {
  return {
    type: types.SERVER_OP_ERROR,
    uuid,
    error,
  };
}

function receiveServers(servers) {
  return {
    type: types.RECEIVE_SERVERS,
    servers: servers.servers,
    configs: servers.configs,
  };
}

function shouldFetchServers(state) {
  const servers = state.servers;

  if (!servers.items) {
    return true;
  } else if (servers.isFetching) {
    return false;
  } else {
    return servers.didInvalidate;
  }
}

function fetchServers() {
  return dispatch => {
    dispatch(requestServers());

    return fetchApi('/api/servers')
      .then(response =>
        response.json().then(json => {
          if (!response.ok) {
            return Promise.reject(json);
          }

          dispatch(receiveServers(json));
        })
      ).catch(promiseErrorHandler);
  };
}

export function fetchServersIfNeeded() {
  return (dispatch, getState) => {
    if (shouldFetchServers(getState())) {
      return dispatch(fetchServers());
    }
  };
}

export function startServer(configUuid) {
  return dispatch => {
    const formData = new FormData();
    formData.append('serverCfgUuid', configUuid);

    return fetchApi(`/api/servers/start`, { method: 'POST', body: formData })
      .then(response => {
        return response.json().then(json => {
          if (!response.ok) {
            return Promise.reject(json);
          }

          dispatch(serverStarted(json));
        });
      }).catch(promiseErrorHandler);
  };
}

export function startServerTmp(formData) {
  return dispatch => {
    fetchApi(`/api/servers/start/upload`, { method: 'POST', body: formData })
      .then(response =>
        response.json().then(json => {
          if (response.status !== 200) {
            requestErrorHandler(response, json);
            return;
          }

          dispatch(serverStarted(json));
        })
      ).catch(promiseErrorHandler);
  };
}

export function stopServer(uuid) {
  return dispatch => {
    dispatch(serverLoadingStart(uuid));

    return fetchApi(`/api/servers/${uuid}/stop`, { method: 'POST' })
      .then(response => {
        if (response.status === 204) {
          dispatch(serverStopped(uuid));
        } else {
          return response.json().then(json => {
            dispatch(serverLoadingError(uuid, json.message));
          });
        }
      }).catch(promiseErrorHandler);
  };
}

export function reconfigServer(uuid, formData) {
  return dispatch => {
    dispatch(serverLoadingStart(uuid));

    fetchApi(`/api/servers/${uuid}/reconfig`, { method: 'POST', body: formData })
      .then(response =>
        response.json().then(json => {
          if (response.status !== 200) {
            requestErrorHandler(response, json);
            return;
          }

          dispatch(serverLoadingStop(uuid, json));
        })
      ).catch(promiseErrorHandler);
  };
}
