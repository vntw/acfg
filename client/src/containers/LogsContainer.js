import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { fetchLogsIfNeeded, fetchLog } from '../actions/logs';
import LoadIndicator from '../components/LoadIndicator';
import Logs from '../components/Logs';

class LogsContainer extends Component {
  componentDidMount() {
    this.props.fetchLogsIfNeeded();
  }

  render() {
    const { logs, isFetching } = this.props.logs;

    if (isFetching) {
      return <LoadIndicator/>;
    }

    return (
      <Logs
        logs={logs}
        showLog={this.props.showLog}
      />
    );
  }
}

LogsContainer.propTypes = {
  logs: PropTypes.object.isRequired,
  showLog: PropTypes.func.isRequired,
  fetchLogsIfNeeded: PropTypes.func.isRequired,
};

const mapDispatchToProps = {
  showLog: fetchLog,
  fetchLogsIfNeeded,
};

const mapStateToProps = (state) => ({
  logs: state.logs,
});

export default connect(mapStateToProps, mapDispatchToProps)(LogsContainer);
