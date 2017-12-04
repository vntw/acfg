import React from 'react';
import PropTypes from 'prop-types';
import Modal from 'react-modal';
import IniView from './config/IniView';

class ConfigItem extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.state = { detailsOpen: false };
  }

  render() {
    const { config, startServer, deleteConfig } = this.props;

    return (
      <li className="configs__list-item">
        {config.name}<br />
        <small>
          {config.files.server_cfg.ini.SERVER.TRACK},{' '}
          {config.files.server_cfg.ini.SERVER.CARS.split(';').length} Cars
        </small>
        <br />

        <button
          className="btn"
          onClick={(e) => {
            e.preventDefault();
            startServer(config.uuid);
          }}
        >Start</button>
        <button
          className="btn"
          onClick={(e) => {
            e.preventDefault();
            deleteConfig(config.uuid);
          }}
        >Remove</button>
        <button
          className="btn"
          onClick={(e) => {
            e.preventDefault();
            this.setState({ detailsOpen: true });
          }}
        >Details</button>

        <Modal
          isOpen={this.state.detailsOpen}
          onRequestClose={() => this.setState({ detailsOpen: false })}
        >
          <h3>Config Details</h3>

          <p>
            <small>server_cfg.ini: {config.files.server_cfg.checksum}</small><br />
            <small>entry_list.ini: {config.files.entry_list.checksum}</small><br />
          </p>

          <h4>server_cfg.ini</h4>
          <IniView ini={config.files.server_cfg.ini} />

          <h4>entry_list.ini</h4>
          <IniView ini={config.files.entry_list.ini} />
        </Modal>
      </li>
    );
  }
}

ConfigItem.propTypes = {
  config: PropTypes.object.isRequired,
  startServer: PropTypes.func.isRequired,
  deleteConfig: PropTypes.func.isRequired,
};

export default ConfigItem;
