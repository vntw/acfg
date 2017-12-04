import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { fetchServersIfNeeded } from '../actions/server';
import LoadIndicator from '../components/LoadIndicator';
import ServerContainer from './ServerContainer';
import ConfigsContainer from './ConfigsContainer';

class ServersConfigsContainer extends Component {
  componentDidMount() {
    this.props.fetchServersIfNeeded();
  }

  render() {
    const { isFetching } = this.props.servers;

    if (isFetching) {
      return <LoadIndicator/>;
    }

    return (
      <div className="servers__grid">
        <ServerContainer />

        <ConfigsContainer />
      </div>
    );
  }
}

ServersConfigsContainer.propTypes = {
  servers: PropTypes.object.isRequired,
  fetchServersIfNeeded: PropTypes.func.isRequired,
};

const mapStateToProps = (state) => ({
  servers: state.servers,
});

export default connect(mapStateToProps, { fetchServersIfNeeded })(ServersConfigsContainer);
