import React from 'react';
import PropTypes from 'prop-types';
import Duration from './Duration';

class Timeleft extends React.Component {
  constructor(props) {
    super(props);

    this.state = { seconds: Number(props.seconds) };
  }

  componentDidMount() {
    this.interval = setInterval(::this.tick, 1000);
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  tick() {
    this.setState({ seconds: this.state.seconds - 1 });
    if (this.state.seconds <= 0) {
      clearInterval(this.interval);
    }
  }

  render() {
    return <Duration seconds={this.state.seconds} />;
  }
}

Timeleft.propTypes = {
  seconds: PropTypes.number.isRequired,
};

export default Timeleft;
