import React from 'react';
import PropTypes from 'prop-types';
import { humanDuration } from '../../ac/time';

const Duration = (props) => <span>{humanDuration(props.seconds)}</span>;

Duration.propTypes = {
  seconds: PropTypes.number.isRequired,
};

export default Duration;
