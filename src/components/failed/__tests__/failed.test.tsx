import React from 'react';
import chalk from 'chalk';
import { render } from 'ink-testing-library';

import Failed from '../failed';

describe('components/failed', () => {
  it('should render text with error text', () => {
    const { lastFrame } = render(<Failed />);
    expect(lastFrame()).toEqual(chalk.bold.red('x Action Failed') + '\n');
  });

  it('should render text with empty error prop', () => {
    const { lastFrame } = render(<Failed error={{}} />);
    expect(lastFrame()).toEqual(chalk.bold.red('x Action Failed') + '\n');
  });

  it('should render details text', () => {
    const { lastFrame } = render(
      <Failed error={{ message: 'Incorrect argument' }} />
    );
    expect(lastFrame()).toEqual(
      expect.stringContaining(chalk.dim('Details') + '   Incorrect argument')
    );
  });
});
