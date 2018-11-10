import { connect } from 'react-redux';
import { loginUser, logoutUser } from '../actions/user';
import { withRouter } from 'react-router-dom';
import App from './../components/App';

const mapDispatchToProps = {
  onLogin: loginUser,
  onLogout: logoutUser,
};

export default withRouter(connect(state => ({ user: state.auth }), mapDispatchToProps)(App));
