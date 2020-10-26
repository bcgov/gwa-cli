import React from 'react';
import { render } from 'ink-testing-library';

import Loading from '../loading';

describe('components/loading', () => {
  it('should render with children', () => {
    const { lastFrame } = render(<Loading>Text</Loading>);
    expect(lastFrame()).toEqual('⠋ Text');
  });

  it('should accept a different spinner', () => {
    const { lastFrame } = render(<Loading spinner="arc">Arc</Loading>);
    expect(lastFrame()).toEqual('◜ Arc');
  });
});
