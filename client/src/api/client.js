import Alert from 'react-s-alert';
import { getToken } from './auth';
const appCfg = require('cfg');

export default function fetchApi(endpoint, options, authenticated = true) {
  let token = getToken();
  let config = { ...options };

  if (authenticated) {
    if (!token) {
      Alert.error('Please login to continue');
      return;
    }

    config.headers = { 'Authorization': `Bearer ${token}` };
  }

  return fetch(appCfg.API_URL + endpoint, config);
}
