import { connect } from 'react-redux';
import { loginUser, logoutUser } from '../actions/user';
import App from '../components/App';

const mapDispatchToProps = {
  onLogin: loginUser,
  onLogout: logoutUser,
};

export default connect(state => ({ user: state.auth }), mapDispatchToProps)(App);
