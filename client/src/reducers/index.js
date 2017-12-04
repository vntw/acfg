import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';
import auth from './auth';
import servers from './servers';
import configs from './configs';
import logs from './logs';

const rootReducer = combineReducers({
  auth,
  servers,
  configs,
  logs,
  routing: routerReducer
});

export default rootReducer;
