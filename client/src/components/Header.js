import React from 'react';
import PropTypes from 'prop-types';
import { NavLink } from 'react-router-dom';

const Header = ({ onLogout }) => (
  <header>
    <nav>
      <NavLink to="/" className="header__logo">
        <img src="/static/img/ac-logo.png" />
      </NavLink>
      <ul>
        <li><NavLink className="btn btn--primary" to="/">Servers</NavLink></li>
        <li><NavLink className="btn btn--primary" to="/logs">Logs</NavLink></li>
        <li><button onClick={onLogout} className="btn btn--secondary">Logout</button></li>
      </ul>
    </nav>
  </header>
);

Header.propTypes = {
  onLogout: PropTypes.func.isRequired,
};

export default Header;
