import React, { Suspense } from 'react';
import chalk from 'chalk';
import delay from 'delay';
import { render } from 'ink-testing-library';
import { Text } from 'ink';

import { makeConfigFile } from '../../create-actions';
import WriteConfigAction from '../write-config-action';

jest.mock('../../create-actions');

describe('commands/create/views/write-config-action', () => {
  it('should render success text', async () => {
    makeConfigFile.mockResolvedValueOnce('test.yaml');
    const data = {
      url: 'https://swagger.url/json',
      other: 'options',
      outfile: 'test.yaml',
    };
    const { lastFrame } = render(
      <Suspense fallback={<Text>Loading...</Text>}>
        <WriteConfigAction data={data} />
      </Suspense>
    );
    expect(lastFrame()).toEqual('Loading...');

    await delay(100);

    expect(makeConfigFile).toHaveBeenCalledWith('https://swagger.url/json', {
      other: 'options',
      outfile: 'test.yaml',
    });
    expect(lastFrame()).toEqual(
      `${chalk.bold.green('✓')} ${chalk.bold(
        'Config test.yaml successfully generated'
      )}`
    );
  });
});
