import React from 'react';
import chalk from 'chalk';
import { Text } from 'ink';
import { render } from 'ink-testing-library';

import Success from '../success';

describe('components/failed', () => {
  it('should render text with error text', () => {
    const { lastFrame } = render(<Success>Complete</Success>);
    expect(lastFrame()).toEqual(
      chalk.bold.green('✓') + ' ' + chalk.bold('Complete')
    );
  });

  it('should render Text children', () => {
    const { lastFrame } = render(
      <Success>
        <Text underline>Complete</Text>
      </Success>
    );
    expect(lastFrame()).toEqual(
      chalk.bold.green('✓') + ' ' + chalk.bold.underline('Complete')
    );
  });
});
