import React from 'react';
import { Link } from 'react-router-dom';

const NotFoundPage = () => {
  return (
    <div className="simple-page__container">
      <h1>404 Page Not Found</h1>
      <Link to="/">Go back to homepage</Link>
    </div>
  );
};

export default NotFoundPage;
