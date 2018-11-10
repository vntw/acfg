import React from 'react';
import PropTypes from 'prop-types';
import Alert from 'react-s-alert';
import { hot } from "react-hot-loader";
import { Route, Switch } from 'react-router-dom';
import Header from './Header';
import Login from './Login';
import AlertTemplate from './alert/AlertTemplate';
import ServersConfigsContainer from './../containers/ServersConfigsContainer';
import LogsContainer from './../containers/LogsContainer';
import NotFoundPage from './/NotFoundPage';

// This is a class-based component because the current
// version of hot reloading won't hot reload a stateless
// component at the top-level.
class App extends React.Component {
  render() {
    const { user, onLogin, onLogout } = this.props;

    if (user.isLoggingIn) {
      return <div>Logging In</div>;
    } else if (!user.isAuthenticated) {
      return (
        <Login
          onLogin={onLogin}
          errorMessage={user.errorMessage}
        />
      );
    }

    return (
      <div>
        <Header onLogout={onLogout} />
        <Switch>
          <Route exact path="/" component={ServersConfigsContainer} />
          <Route path="/logs" component={LogsContainer} />
          <Route component={NotFoundPage} />
        </Switch>
        <Alert
          stack
          timeout={5000}
          position="bottom"
          effect="scale"
          contentTemplate={AlertTemplate}
        />
      </div>
    );
  }
}

App.propTypes = {
  user: PropTypes.shape({
    isFetching: PropTypes.bool.isRequired,
    isAuthenticated: PropTypes.bool.isRequired,
    user: PropTypes.string,
    errorMessage: PropTypes.string,
  }),
  onLogin: PropTypes.func.isRequired,
  onLogout: PropTypes.func.isRequired,
};

export default hot(module)(App);
