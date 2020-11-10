jest.mock('../../../services/app');
jest.mock('../renderer');
import chalk from 'chalk';
import { actionHandler } from '../init';
import * as app from '../../../services/app';
import render from '../renderer';

describe('commands/init', () => {
  beforeEach(() => {
    app.checkForEnvFile.mockReturnValue(false);
  });

  it('should throw if a .env file exists', () => {
    app.checkForEnvFile.mockReturnValueOnce(true);
    expect(() => actionHandler({})).toThrow();
  });

  it('should render UI if no arguments are passed', () => {
    actionHandler();
    expect(render).toHaveBeenCalled();
  });

  it('should default to test env', () => {
    app.makeEnvFile.mockImplementationOnce(() => Promise.resolve());
    actionHandler({
      namespace: 'sampler',
    });

    expect(app.makeEnvFile).toHaveBeenCalledWith(
      expect.objectContaining({
        env: 'test',
      })
    );
  });

  it('should process arguments for cli action', () => {
    app.makeEnvFile.mockResolvedValueOnce('success');
    actionHandler({
      namespace: 'sampler',
      clientId: 'client-sampler',
      clientSecret: 'client-sampler',
      dev: true,
      ignore: 'me',
    });

    expect(app.makeEnvFile).toHaveBeenCalledWith(
      expect.objectContaining({
        namespace: 'sampler',
        clientId: 'client-sampler',
        clientSecret: 'client-sampler',
        env: 'dev',
      })
    );
  });

  it('should log success message', async () => {
    const spy = jest.spyOn(console, 'log');
    app.makeEnvFile.mockResolvedValueOnce('.env file generated');
    await actionHandler({
      namespace: 'sampler',
      clientId: 'client-sampler',
      clientSecret: 'client-sampler',
    });

    expect(spy).toHaveBeenCalledWith(
      chalk.green.bold('Success'),
      '.env file generated'
    );
  });

  it('should log failed message', async () => {
    const spy = jest.spyOn(console, 'error');
    app.makeEnvFile.mockRejectedValueOnce(new Error('bad config'));

    await actionHandler({
      namespace: 'bad-sampler',
      clientId: 'bad-client-sampler',
      clientSecret: 'bad-client-sampler',
    });

    expect(process.exitCode).toEqual(1);
    expect(spy).toHaveBeenCalledWith(
      chalk.red.bold('x Error'),
      'Unable to create .env file'
    );
  });

  it('should support debug', async () => {
    const spy = jest.spyOn(console, 'error');
    app.makeEnvFile.mockRejectedValueOnce('bad config');

    await actionHandler({
      namespace: 'sampler',
      clientId: 'client-sampler',
      clientSecret: 'client-sampler',
      debug: true,
    });

    expect(process.exitCode).toEqual(1);
    expect(spy).toHaveBeenCalledWith(
      chalk.red.bold('x Error'),
      'Unable to create .env file'
    );
    expect(spy).toHaveBeenLastCalledWith('bad config');
  });
});
