import React from 'react';
import { shallow } from 'enzyme';
import LoadIndicator from './LoadIndicator';

describe('<LoadIndicator />', () => {
  it('renders a load indicator', () => {
    const wrapper = shallow(<LoadIndicator />);
    expect(wrapper).toMatchSnapshot();
  });
});
