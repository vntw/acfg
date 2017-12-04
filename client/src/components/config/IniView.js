import React from 'react';
import PropTypes from "prop-types";

const IniView = ({ ini }) => (
  <ul>
    {Object.keys(ini).map(section => (
      <li key={section}>
        <b>{section}</b><br />
        {Object.keys(ini[section]).map(key => (
          <span key={section + key}>{key}: {ini[section][key]}<br /></span>
        ))}
      </li>
    ))}
  </ul>
);

IniView.propTypes = {
  ini: PropTypes.object.isRequired,
};

export default IniView;
