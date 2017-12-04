import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import Dropzone from 'react-dropzone';
import { startServer } from '../actions/server';
import { deleteConfig, uploadConfig } from '../actions/config';
import ConfigList from '../components/ConfigList';
import handleDrop from '../ac/dropzone';

class ConfigsContainer extends Component {
  render() {
    const { startServer, deleteConfig, configs } = this.props;

    return (
      <Dropzone
        className="configs__container"
        onDrop={handleDrop((formData) => { this.props.uploadConfig(formData); })}
        disableClick
        multiple
      >
        <ConfigList
          configs={configs}
          startServer={startServer}
          deleteConfig={deleteConfig}
        />
      </Dropzone>
    );
  }
}

ConfigsContainer.propTypes = {
  configs: PropTypes.object.isRequired,
  startServer: PropTypes.func.isRequired,
  deleteConfig: PropTypes.func.isRequired,
  uploadConfig: PropTypes.func.isRequired,
};

const mapDispatchToProps = {
  startServer,
  deleteConfig,
  uploadConfig,
};

const mapStateToProps = (state) => ({
  configs: state.configs,
});

export default connect(mapStateToProps, mapDispatchToProps)(ConfigsContainer);
