import React from 'react';
import PropTypes from 'prop-types';
import { Link, IndexLink } from 'react-router';

const Header = ({ onLogout }) => (
  <header>
    <nav>
      <IndexLink to="/" className="header__logo">
        <img src="/static/img/ac-logo.png" />
      </IndexLink>
      <ul>
        <li><IndexLink className="btn btn--primary" to="/">Servers</IndexLink></li>
        <li><Link className="btn btn--primary" to="/logs">Logs</Link></li>
        <li><button onClick={onLogout} className="btn btn--secondary">Logout</button></li>
      </ul>
    </nav>
  </header>
);

Header.propTypes = {
  onLogout: PropTypes.func.isRequired,
};

export default Header;
