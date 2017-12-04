import React from 'react';
import PropTypes from 'prop-types';

const AlertTemplate = (props) => (
  <div className={props.classNames} id={props.id} style={props.styles}>
    <div className="s-alert-box-flag">
      <div className="s-alert-box-flag-inner" />
    </div>
    <div className="s-alert-box-inner">
      {props.message}
    </div>
    <span className="s-alert-close" onClick={props.handleClose} />
  </div>
);

AlertTemplate.propTypes = {
  id: PropTypes.string.isRequired,
  classNames: PropTypes.string.isRequired,
  styles: PropTypes.object.isRequired,
  message: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.object,
  ]).isRequired,
  handleClose: PropTypes.func.isRequired,
  customFields: PropTypes.object,
};

export default AlertTemplate;
