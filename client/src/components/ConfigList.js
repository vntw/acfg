import React from 'react';
import PropTypes from 'prop-types';
import LoadIndicator from './LoadIndicator';
import ConfigItem from './ConfigItem';

class ConfigList extends React.Component {
  render() {
    const { configs, isFetching } = this.props.configs;

    if (isFetching) {
      return <LoadIndicator/>;
    }

    if (configs.length === 0) {
      return (
        <div className="no-items">
          <h4>No presets</h4>
          <p>Drag'n'Drop a preset directory or both config files here (server_cfg.ini and entry_list.ini)</p>
        </div>
      );
    }

    return (
      <ul className="configs__list">
        {configs.map(config => (
          <ConfigItem
            key={config.uuid}
            config={config}
            startServer={this.props.startServer}
            deleteConfig={this.props.deleteConfig}
          />
        ))}
      </ul>
    );
  }
}

ConfigList.propTypes = {
  configs: PropTypes.object.isRequired,
  startServer: PropTypes.func.isRequired,
  deleteConfig: PropTypes.func.isRequired,
};

export default ConfigList;
