import React from 'react';
import PropTypes from 'prop-types';
import LoadIndicator from './LoadIndicator';

class Log extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.state = { open: false };
  }

  onShow(e) {
    e.preventDefault();
    this.setState({ open: true }, () => {
      this.props.showLog(this.props.log.instanceUuid);
    });
  }

  onHide(e) {
    e.preventDefault();
    this.setState({ open: false });
  }

  render() {
    const { log } = this.props;

    return (
      <li className="log__item">
        <p><b>{new Date(log.time * 1000).toLocaleString()}</b> {log.instanceUuid}</p>
        <p>Available Logs: {log.logFiles.map(logFile => logFile.name).join(', ')}</p>

        {log.isLoaded && this.state.open && log.logFiles.map(logFile => (
          <div key={logFile.name}>
            <p>{logFile.name}</p>

            {
              <p className="log__content">
                {logFile.content.split('\n').map((item, key) => item && <span key={key}>{item}<br/></span>)}
              </p>
            }
          </div>
        ))}

        {!log.isLoaded && this.state.open && <LoadIndicator />}

        {this.state.open && <button className="btn btn__log-item btn__log-item--close" onClick={::this.onHide}>Hide</button>}

        {(!log.isLoaded || !this.state.open) && <button className="btn btn__log-item btn__log-item--open" onClick={::this.onShow}>Show</button>}
      </li>
    );
  }
}

Log.propTypes = {
  log: PropTypes.object.isRequired,
  showLog: PropTypes.func.isRequired,
};

export default Log;
