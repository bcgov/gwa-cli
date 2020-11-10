jest.mock('../../../services/app');
import React from 'react';
import delay from 'delay';
import { render } from 'ink-testing-library';
import { makeEnvFile } from '../../../services/app';

import CreateEnvView from '../create-env';

describe('views/create-env', () => {
  const prompts = [
    {
      label: 'Namespace',
      key: 'namespace',
      constraint: {},
    },
  ];
  const env = 'test';

  it('Starts with prompt form', () => {
    const { lastFrame } = render(<CreateEnvView env={env} prompts={[]} />);
    expect(lastFrame()).toEqual(
      expect.stringContaining("Configure this folder's environment variables")
    );
  });

  it('should render prompts', () => {
    const { lastFrame } = render(<CreateEnvView env={env} prompts={prompts} />);
    expect(lastFrame()).toEqual(expect.stringContaining('Namespace'));
  });

  it('should render write component after form is complete', async () => {
    const successText = '.env file successfully generated';
    makeEnvFile.mockImplementationOnce(() => Promise.resolve(successText));
    const { lastFrame, stdin } = render(
      <CreateEnvView env={env} prompts={prompts} />
    );

    expect.assertions(2);
    await delay(100);
    stdin.write('value');
    await delay(100);
    stdin.write('\r');
    await delay(1);
    expect(lastFrame()).toEqual(expect.stringContaining('Writing .env file'));
    await delay(100);
    expect(lastFrame()).toEqual(expect.stringContaining(successText));
  });

  it('should render error', async () => {
    const errorText = 'unable to make file';
    makeEnvFile.mockImplementationOnce(() => Promise.reject(errorText));
    const { lastFrame, stdin } = render(
      <CreateEnvView env={env} prompts={prompts} />
    );

    await delay(100);
    stdin.write('will not work');
    await delay(100);
    stdin.write('\r');
    await delay(100);
    expect(lastFrame()).toEqual(expect.stringContaining(errorText));
  });
});
