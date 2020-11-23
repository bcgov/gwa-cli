import React, { Suspense } from 'react';
import chalk from 'chalk';
import delay from 'delay';
import { Text } from 'ink';
import { render } from 'ink-testing-library';

import api from '../../../services/api';
import RequestView from '../request-view';

jest.mock('../../../services/api');

describe('commands/acl/request-view', () => {
  it('should call a promise', async () => {
    api.mockResolvedValueOnce({
      added: 1,
      removed: 10,
      missing: 0,
    });
    render(
      <Suspense fallback={<Text>Loading...</Text>}>
        <RequestView data={{}} />
      </Suspense>
    );

    await delay(100);

    expect(api).toHaveBeenCalledWith('/namespaces/:namespace/membership', {
      method: 'PUT',
      body: JSON.stringify({}),
    });
  });

  it('should print ACL request results', async () => {
    api.mockResolvedValueOnce({
      added: 1,
      removed: 10,
      missing: 0,
    });
    const { lastFrame } = render(
      <Suspense fallback={<Text>Loading...</Text>}>
        <RequestView data={{}} />
      </Suspense>
    );

    await delay(100);
    expect(lastFrame()).toEqual(`${chalk.bold.green('✓')}${chalk(
      ' '
    )}${chalk.bold('Access Control Updated!')}
${chalk.green('+ Added')}${chalk('   1')}
${chalk.red('- Removed')}${chalk(' ')}10
${chalk.yellow('? Missing')}${chalk(' ')}0`);
  });
});
