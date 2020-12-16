import * as React from 'react';
import chalk from 'chalk';
import delay from 'delay';
import { Text } from 'ink';
import { render } from 'ink-testing-library';

import api from '../../../../services/api';
import StatusView from '../status-view';

jest.mock('../../../../services/api');
const CACHED_ENV = process.env;

describe('commands/status/views/service-item', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    process.env = {
      ...process.env,
      GWA_NAMESPACE: 'sampler',
    };
  });

  afterEach(() => {
    process.env = CACHED_ENV;
  });

  it('should handle async request and no results', async () => {
    api.mockResolvedValue([]);
    const { lastFrame } = render(
      <React.Suspense fallback={<Text>Fetching Status...</Text>}>
        <StatusView />
      </React.Suspense>
    );
    expect(lastFrame()).toEqual('Fetching Status...');
    await delay(100);
    expect(lastFrame()).toEqual(`
sampler Status

You have no services yet.`);
  });

  it('should render a status list', async () => {
    api.mockResolvedValue([
      {
        name: 'My Service',
        envHost: 'api.host.xyz',
        reason: '200 Response',
        upstream: 'https://httpbin.org',
        status: 'UP',
      },
    ]);
    const { lastFrame } = render(
      <React.Suspense fallback={<Text>Fetching Status...</Text>}>
        <StatusView />
      </React.Suspense>
    );
    await delay(100);
    //expect(api).toHaveBeenCalledWith('/namespaces/sampler/services');
    expect(lastFrame()).toEqual(`
sampler Status

${chalk.greenBright`▲`} ${chalk.greenBright`My Service`}        200 Response                  ${chalk.dim`api.host.xyz [https://httpbin.org]`}`);
  });
});
