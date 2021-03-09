import React from 'react';
import chalk from 'chalk';
import ink from 'ink';

import renderer from '../renderer';

jest.mock('ink', () => ({
  Box: jest.fn(() => 'Box'),
  Text: jest.fn(() => 'Text'),
  render: jest.fn(),
}));

describe('commands/acl/ui', () => {
  it('should render the init UI', () => {
    renderer({});
    expect(ink.render).toHaveBeenCalled();
  });
});
