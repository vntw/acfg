import React from 'react';
import PropTypes from 'prop-types';
import Log from './Log';

const Logs = ({ logs, showLog }) => {
  if (logs.length === 0) {
    return (
      <div className="no-items">
        <h2>¯\_(ツ)_/¯</h2>
        <h4>No logs</h4>
      </div>
    );
  }

  return (
    <ul className="logs__container">
      {logs.map(log => (
        <Log
          key={log.instanceUuid}
          log={log}
          showLog={showLog}
        />
      ))}
    </ul>
  );
};

Logs.propTypes = {
  logs: PropTypes.array.isRequired,
  showLog: PropTypes.func.isRequired,
};

export default Logs;
