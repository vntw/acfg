import React from 'react';
import PropTypes from 'prop-types';

const Login = ({ onLogin, errorMessage }) => {
  let usernameRef;
  let passwordRef;

  return (
    <div className="login__container">
      <div className="login__box">
        <form
          method="post"
          onSubmit={(e) => {
            e.preventDefault();
            onLogin({
              username: usernameRef.value.trim(),
              password: passwordRef.value.trim(),
            });
          }}
        >
          <input
            type="text"
            name="username"
            placeholder="Username"
            className="login__input"
            ref={ref => { usernameRef = ref; }}
            required
            autoFocus
          />
          <input
            type="password"
            name="password"
            placeholder="Password"
            className="login__input"
            ref={ref => { passwordRef = ref; }}
            required
          />
          <input type="submit" className="btn btn--primary" value="Login" />
        </form>
        {errorMessage && <div className="login__error">{errorMessage}</div>}
      </div>
    </div>
  );
};

Login.propTypes = {
  errorMessage: PropTypes.string,
  onLogin: PropTypes.func.isRequired,
};

export default Login;
