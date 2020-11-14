import React from 'react';
import chalk from 'chalk';
import { render } from 'ink-testing-library';
import ink from 'ink';

import renderer from '../renderer';
jest.mock('ink', () => ({
  render: jest.fn(),
}));

describe('commands/init/ui', () => {
  it('should render the init UI', () => {
    renderer('dev');
    expect(ink.render).toHaveBeenCalled();
  });
});
