import Alert from 'react-s-alert';

const requestErrorHandler = (response, body) => {
  if (response.status < 400) {
    return;
  }

  if (body && body.type && body.type === 'error' && body.message) {
    Alert.error(body.message);
    return;
  }

  Alert.error(`Got request error: ${response.statusCode}`);
};

const promiseErrorHandler = (reason) => {
  Alert.error(`Request error: ${reason}`);
};

export {
  requestErrorHandler,
  promiseErrorHandler,
};
