import React from 'react';
import delay from 'delay';
import { render } from 'ink-testing-library';
import * as app from '../../../services/app';

import CreateEnvView from '../create-env';
jest.mock('../../../services/app');

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
    const successText = 'Done';
    app.makeEnvFile = jest.fn().mockResolvedValueOnce(successText);
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
    const errorText = 'Error';
    app.makeEnvFile = jest.fn().mockRejectedValueOnce(errorText);
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
