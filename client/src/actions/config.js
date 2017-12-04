import * as types from '../constants/action-types';
import fetchApi from '../api/client';
import { requestErrorHandler, promiseErrorHandler } from './error-handlers';

function receiveUploadedConfig(config) {
  return {
    type: types.RECEIVE_UPLOADED_CONFIG,
    config,
  };
}

function configDeleted(id) {
  return {
    type: types.CONFIG_DELETED,
    id,
  };
}

export function deleteConfig(id) {
  return dispatch => {
    fetchApi(`/api/configs/${id}/delete`, { method: 'DELETE' })
      .catch(promiseErrorHandler);

    return dispatch(configDeleted(id));
  };
}

export function uploadConfig(formData) {
  return dispatch => {
    fetchApi('/api/configs/upload', { method: 'POST', body: formData })
      .then(response =>
        response.json().then(json => {
          if (response.status !== 200) {
            requestErrorHandler(response, json);
            return;
          }

          dispatch(receiveUploadedConfig(json));
        })
      ).catch(promiseErrorHandler);
  };
}
