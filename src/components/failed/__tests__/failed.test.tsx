import React from 'react';
import chalk from 'chalk';
import { render } from 'ink-testing-library';

import Failed from '../failed';

describe('components/failed', () => {
  it('should render text with error text', () => {
    const { lastFrame } = render(<Failed />);
    expect(lastFrame()).toEqual(chalk.bold.red('x Error'));
  });

  it('should render text with empty error prop', () => {
    const { lastFrame } = render(<Failed error={{}} />);
    expect(lastFrame()).toEqual(chalk.bold.red('x Error'));
  });

  it('should render details text', () => {
    const { lastFrame } = render(
      <Failed error={{ message: 'Incorrect argument' }} />
    );
    expect(lastFrame()).toEqual(
      chalk.bold.red('x Error') + ' Incorrect argument'
    );
  });

  it('should render stack when --verbose', () => {
    const { lastFrame } = render(
      <Failed
        verbose
        error={{ message: 'Incorrect argument', stack: 'Error code details' }}
      />
    );
    expect(lastFrame()).toEqual(`${chalk.bold.red('x Error')} Incorrect argument

${chalk.dim('Details')}   Error code details`);
  });
});
