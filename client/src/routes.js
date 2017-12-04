import React from 'react';
import { Route, IndexRoute } from 'react-router';

import AppContainer from './containers/AppContainer';
import LogsContainer from './containers/LogsContainer';
import NotFoundPage from './components/NotFoundPage';
import ServersConfigsContainer from './containers/ServersConfigsContainer';

export default (
  <Route path="/" component={AppContainer}>
    <IndexRoute component={ServersConfigsContainer}/>
    <Route path="logs" component={LogsContainer}/>
    <Route path="*" component={NotFoundPage}/>
  </Route>
);
