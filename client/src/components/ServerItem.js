import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Dropzone from 'react-dropzone';
import Modal from 'react-modal';
import { getShortSessionName } from '../ac/session';
import LoadIndicator from './LoadIndicator';
import Timeleft from './server/Timeleft';
import IniView from './config/IniView';

class ServerItem extends Component {
  static calculateTime(timeOfDay) {
    const totalHours = ((timeOfDay / 8) * 30 + 780) / 60;
    const hrs = totalHours.toString();

    return Math.floor(totalHours).toString() + ':' + (hrs[hrs.length - 1] === '5' ? '30' : '00');
  }

  constructor(props, context) {
    super(props, context);

    this.state = { tmpConfigModalOpen: false };
  }

  renderModal() {
    const { server } = this.props;

    return (
      <Modal
        isOpen={this.state.tmpConfigModalOpen}
        onRequestClose={() => this.setState({ tmpConfigModalOpen: false })}
        ariaHideApp={false}
        shouldReturnFocusAfterClose={false}
      >
        <h3>Server Details</h3>

        <p>{server.status.name}</p>

        <p>
          <small>Instance: {server.instance.uuid}</small><br />
          <small>Config: {server.instance.spec.preset.uuid}</small><br />
          <small>TmpConfig: {server.instance.spec.tmpConfig.uuid}</small><br />
        </p>

        <h3>Configs</h3>
        <h4>server_cfg.ini</h4>
        <IniView ini={server.instance.spec.tmpConfig.files.server_cfg.ini} />

        <h4>entry_list.ini</h4>
        <IniView ini={server.instance.spec.tmpConfig.files.entry_list.ini} />
      </Modal>
    );
  }

  render() {
    const { server } = this.props;

    if (server.isLoading) {
      return <LoadIndicator />;
    }

    return (
      <Dropzone
        onDrop={(acceptedFiles, rejectedFiles, e) => this.props.handleReconfig(server.instance.uuid)(acceptedFiles, rejectedFiles, e)}
        className="server__item"
        disableClick
        multiple
      >
        <h4 className="name">{server.status.pass ? '🔒 ' : ''}{server.status.name}</h4>
        <div className="sub clients">{server.status.clients}/{server.status.maxclients} drivers</div>
        <div className="sub track">{server.status.track}</div>
        <div className="sub time">{ServerItem.calculateTime(server.status.timeofday)}h</div>
        <div className="sub sessions">
          {server.status.sessiontypes.map((s, i) => {
            const isActive = i === server.status.session;

            return (
              <span key={s} className={`session ${isActive ? 'active' : ''}`}>
                {getShortSessionName(s)}{' '}
                <small>({isActive ? <Timeleft startSeconds={server.status.timeleft} /> : server.status.durations[i]})</small>
              </span>
            );
          })}
        </div>
        <div className="sub cars">{server.status.cars.join(', ')}</div>
        <div className="sub net">
          TCP: {server.instance.spec.ports.tcpPort}{' '}|{' '}
          UDP: {server.instance.spec.ports.udpPort}{' '}|{' '}
          HTTP: <a href={`http://${server.instance.ip}:${server.instance.spec.ports.httpPort}/INFO`}>{server.instance.spec.ports.httpPort}</a>
          {' '}|{' '}
          {server.instance.spec.plugins.map(plugin => (
            <span key={plugin.name}>
              <b>{plugin.name}:</b>{' '}
              TCP: {plugin.tcpPort}{' '}
              HTTP: <a href={`http://${server.instance.ip}:${plugin.httpPort}`}>{plugin.httpPort}</a>
            </span>
          ))}
        </div>
        <div className="controls">
          <button
            className="btn"
            disabled={server.isLoading}
            onClick={() => this.props.handleStop(server.instance.uuid)}
          >Stop</button>

          <button
            className="btn"
            disabled={server.isLoading}
            onClick={(e) => {
              e.preventDefault();
              this.setState({ tmpConfigModalOpen: true });
            }}
          >Details</button>
        </div>
        {this.renderModal()}
      </Dropzone>
    );
  }
}

ServerItem.propTypes = {
  server: PropTypes.object.isRequired,
  handleStop: PropTypes.func.isRequired,
  handleReconfig: PropTypes.func.isRequired,
};

export default ServerItem;
