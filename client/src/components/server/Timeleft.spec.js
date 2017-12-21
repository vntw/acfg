import React from 'react';
import { shallow } from 'enzyme';
import Timeleft from './Timeleft';

describe('<Timeleft />', () => {
  const dataProvider = {
    'should pad mins/secs < 10 with zero': {
      sec: 3905,
      expected: '1h 05m 05s',
    },
    'removes zero hour': {
      sec: 65,
      expected: '01m 05s',
    },
    'removes zero minute': {
      sec: 5,
      expected: '05s',
    },
    'removes zero second': {
      sec: 0,
      expected: '-',
    },
  };

  for (const expectation in dataProvider) {
    it(expectation, () => {
      const { sec, expected } = dataProvider[expectation];
      const wrapper = shallow(<Timeleft startSeconds={sec} />);
      const actual = wrapper.find('span').text();
      expect(actual).toEqual(expected);
    });
  }
});
