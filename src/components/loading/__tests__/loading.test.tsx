import React from 'react';
import { render } from 'ink-testing-library';

import Loading from '../loading';

describe('components/loading', () => {
  it('should render with children', () => {
    const { lastFrame } = render(<Loading>Text</Loading>);
    expect(lastFrame()).toEqual(expect.stringContaining('Text'));
  });
});
