import * as React from 'react';
import chalk from 'chalk';
import { render } from 'ink-testing-library';

import ServiceItem from '../service-item';

describe('commands/status/views/service-item', () => {
  it('should render an active service', () => {
    const { lastFrame } = render(
      <ServiceItem
        data={{
          name: 'My Service',
          envHost: 'api.host.xyz',
          reason: '200 Response',
          upstream: 'https://httpbin.org',
          status: 'UP',
        }}
      />
    );
    expect(lastFrame()).toEqual(
      `${chalk.greenBright`▲`} ${chalk.greenBright`My Service`}        200 Response                  ${chalk.dim`api.host.xyz [https://httpbin.org]`}`
    );
  });

  it('should render an inactive service', () => {
    const { lastFrame } = render(
      <ServiceItem
        data={{
          name: 'My Service',
          envHost: 'api.host.xyz',
          reason: '500 Response',
          upstream: 'https://httpbin.org',
          status: 'DOWN',
        }}
      />
    );
    expect(lastFrame()).toEqual(
      `${chalk.redBright`▼`} ${chalk.redBright`My Service`}        500 Response                  ${chalk.dim`api.host.xyz [https://httpbin.org]`}`
    );
  });

  it('should just show upstream if it matches host (in prod)', () => {
    const { lastFrame } = render(
      <ServiceItem
        data={{
          name: 'My Service',
          envHost: 'https://httpbin.org',
          reason: '500 Response',
          upstream: 'https://httpbin.org',
          status: 'DOWN',
        }}
      />
    );
    expect(lastFrame()).toEqual(
      `${chalk.redBright`▼`} ${chalk.redBright`My Service`}        500 Response                  ${chalk.dim`https://httpbin.org`}`
    );
  });
});
