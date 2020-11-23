import React, { Suspense } from 'react';
import chalk from 'chalk';
import delay from 'delay';
import { render } from 'ink-testing-library';

import { makeConfigFile } from '../../create-actions';
import CreateConfigView from '../create-config';

jest.mock('../../create-actions');

describe('commands/create/views/create-config', () => {
  const prompts = [
    {
      label: 'Test Prompt',
      key: 'val',
      constraint: {
        presence: { allowEmpty: false },
      },
    },
  ];

  it('render prompts', () => {
    const { lastFrame } = render(<CreateConfigView prompts={prompts} />);
    expect(lastFrame()).toEqual(`
Create a new configuration file

${chalk.bold.green('?')} ${chalk.bold('Test Prompt')} ${chalk.inverse(' ')}`);
  });

  it('Render async action', async () => {
    const { lastFrame, stdin } = render(<CreateConfigView prompts={prompts} />);
    makeConfigFile.mockResolvedValueOnce('file.yaml');

    await delay(100);
    stdin.write('value');
    await delay(100);
    stdin.write('\r');
    await delay(100);

    expect(makeConfigFile).toHaveBeenCalled();
    expect(lastFrame()).toEqual(`
Create a new configuration file

${chalk.bold.green('✓')} ${chalk.bold('Test Prompt')} value
${chalk.bold.green('✓')} ${chalk.bold(
      'Config file.yaml successfully generated'
    )}`);
  });
});
