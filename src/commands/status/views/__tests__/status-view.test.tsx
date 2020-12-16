import * as React from 'react';
import chalk from 'chalk';
import delay from 'delay';
import { Text } from 'ink';
import { render } from 'ink-testing-library';

import { cache } from '../../../../hooks/use-async';
import api from '../../../../services/api';
import StatusView from '../status-view';

jest.mock('../../../../services/api');
const CACHED_ENV = process.env;

describe('commands/status/views/service-item', () => {
  beforeEach(() => {
    process.env = {
      ...CACHED_ENV,
      GWA_NAMESPACE: 'sampler',
    };
  });

  afterEach(() => {
    process.env = CACHED_ENV;
    cache.clear();
  });

  it('should handle async request and no results', async () => {
    api.mockResolvedValueOnce([]);
    const { lastFrame } = render(
      <React.Suspense fallback={<Text>Fetching Status...</Text>}>
        <StatusView />
      </React.Suspense>
    );
    expect(lastFrame()).toEqual('Fetching Status...');
    await delay(100);
    expect(api).toHaveBeenCalledWith('/namespaces/:namespace/services', {
      namespace: 'sampler',
    });
    expect(lastFrame()).toEqual(`
sampler Status

You have no services yet.`);
  });

  it('should render a status list', async () => {
    api.mockResolvedValueOnce([
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
    expect(api).toHaveBeenCalledWith('/namespaces/:namespace/services', {
      namespace: 'sampler',
    });
    await delay(100);
    expect(lastFrame()).toEqual(`
sampler Status

${chalk.greenBright`▲`} ${chalk.greenBright`My Service`}        200 Response                  ${chalk.dim`api.host.xyz [https://httpbin.org]`}`);
  });

  it('should exitCode 1 if there is a down service', async () => {
    api.mockResolvedValueOnce([
      {
        name: 'My Service',
        envHost: 'api.host.xyz',
        reason: '200 Response',
        upstream: 'https://httpbin.org',
        status: 'UP',
      },
      {
        name: 'My Other Service',
        envHost: 'api.host.xyz',
        reason: '200 Response',
        upstream: 'https://httpbin.org',
        status: 'DOWN',
      },
    ]);
    const { lastFrame } = render(
      <React.Suspense fallback={<Text>Fetching Status...</Text>}>
        <StatusView />
      </React.Suspense>
    );
    expect(api).toHaveBeenCalledWith('/namespaces/:namespace/services', {
      namespace: 'sampler',
    });
    await delay(100);
    expect(process.exitCode).toEqual(1);
  });
});
