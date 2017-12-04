import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { reconfigServer, startServer, startServerTmp, stopServer } from '../actions/server';
import ServerList from '../components/ServerList';
import handleDrop from '../ac/dropzone';

class ServerContainer extends Component {
  onReConfigDrop(uuid) {
    return handleDrop((formData) => { this.props.reconfigServer(uuid, formData); });
  }

  render() {
    const { startServer, stopServer } = this.props;
    const { servers } = this.props.servers;

    return (
      <div className="servers__container">
        <ServerList
          servers={servers}
          handleStart={startServer}
          handleStartTmp={handleDrop((formData) => { this.props.startServerTmp(formData); })}
          handleStop={stopServer}
          handleReconfig={::this.onReConfigDrop}
        />
      </div>
    );
  }
}

ServerContainer.propTypes = {
  servers: PropTypes.object.isRequired,
  startServer: PropTypes.func.isRequired,
  startServerTmp: PropTypes.func.isRequired,
  stopServer: PropTypes.func.isRequired,
  reconfigServer: PropTypes.func.isRequired,
};

const mapDispatchToProps = {
  startServer,
  startServerTmp,
  stopServer,
  reconfigServer,
};

const mapStateToProps = (state) => ({
  servers: state.servers,
});

export default connect(mapStateToProps, mapDispatchToProps)(ServerContainer);
