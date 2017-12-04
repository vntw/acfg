import React from 'react';
import PropTypes from 'prop-types';

class Timeleft extends React.Component {
  static padZero(num) {
    return num < 10 ? '0' + num.toString() : num;
  }

  constructor(props) {
    super(props);

    this.state = { startSeconds: Number(props.startSeconds) };
  }

  componentDidMount() {
    this.interval = setInterval(::this.tick, 1000);
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  tick() {
    this.setState({ startSeconds: this.state.startSeconds - 1 });
    if (this.state.startSeconds <= 0) {
      clearInterval(this.interval);
    }
  }

  humanTimeLeft() {
    const { startSeconds } = this.state;

    const h = Math.floor(startSeconds / 3600);
    const m = Math.floor(startSeconds % 3600 / 60);
    const s = Math.floor(startSeconds % 3600 % 60);

    if (h + m + s === 0) {
      return '-';
    }

    const hs = h > 0 ? h + 'h ' : '';
    const ms = h > 0 || m > 0 ? Timeleft.padZero(m) + 'm ' : '';
    const ss = h > 0 || m > 0 || s > 0 ? Timeleft.padZero(s) + 's' : '';

    return hs + ms + ss;
  }

  render() {
    return (
      <span data-seconds={this.state.startSeconds}>
        {this.humanTimeLeft()}
      </span>
    );
  }
}

Timeleft.propTypes = {
  startSeconds: PropTypes.number.isRequired,
};

export default Timeleft;
