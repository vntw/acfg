import React from 'react';
import PropTypes from 'prop-types';
import Dropzone from 'react-dropzone';
import ServerItem from './ServerItem';

const ServerList = ({ servers, handleStartTmp, handleStop, handleReconfig }) => {
  const dz = (
    <Dropzone
      onDrop={handleStartTmp}
      className="config__dropzone"
      multiple
    >
      Drag'n'Drop a preset to start a server instantly.
    </Dropzone>
  );

  if (servers.length === 0) {
    return (
      <div className="no-items">
        <h4>No servers</h4>
        {dz}
      </div>
    );
  }

  return (
    <div>
      {servers.map(server => (
        <ServerItem
          key={server.instance.uuid}
          server={server}
          handleStop={handleStop}
          handleReconfig={handleReconfig}
        />
      ))}
      {dz}
    </div>
  );
};

ServerList.propTypes = {
  servers: PropTypes.array.isRequired,
  handleStartTmp: PropTypes.func.isRequired,
  handleStop: PropTypes.func.isRequired,
  handleReconfig: PropTypes.func.isRequired,
};

export default ServerList;
