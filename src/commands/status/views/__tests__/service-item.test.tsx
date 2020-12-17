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
      `${chalk.greenBright`▲`} ${chalk.greenBright`My Service`}                                   200 Response   ${chalk.dim`https://httpbin.org`}`
    );
  });

  it('should render an inactive service', () => {
    const { lastFrame, stdout } = render(
      <ServiceItem
        sm={false}
        data={{
          name: 'apidata-abcde-searh-bc-api-service',
          reason: '404 Response',
          envHost: 'apidata-test-url.abc.test.apsgw.xyz',
          upstream: 'https://updstream.url/path/dir/',
          host: '',
          status: 'DOWN',
        }}
      />
    );
    expect(lastFrame()).toEqual(
      `${chalk.redBright`▼`} ${chalk.redBright`apidata-abcde-searh-bc-api-service`}           404 Response   ${chalk.dim`https://updstream.url/path/dir/`}`
    );
  });

  it('it should render responsive items', () => {
    const { lastFrame } = render(
      <ServiceItem
        sm
        data={{
          name: 'apidata-abcde-searh-bc-api-service',
          reason: '404 Response',
          envHost: 'apidata-test-url.abc.test.apsgw.xyz',
          upstream: 'https://updstream.url/path/dir/',
          host: '',
          status: 'UP',
        }}
      />
    );
    expect(lastFrame()).not.toEqual(
      expect.stringContaining('https://updstream.url/path/dir/')
    );
    expect(lastFrame()).toEqual(
      `${chalk.greenBright`▲`} ${chalk.greenBright`apidata-abcde-searh-bc-api-service`}           404 Response`
    );
  });
});
