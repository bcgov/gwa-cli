import chalk from 'chalk';

import { actionHandler } from '../create';
import { makeConfigFile } from '../create-actions';
import render from '../renderer';

jest.mock('../renderer');
jest.mock('../create-actions');

describe('commands/create', () => {
  it('should render interface if no arguments', async () => {
    await actionHandler(undefined, {});
    expect(render).toHaveBeenCalled();
  });

  it('should call makeConfig', async () => {
    const input = 'file.yaml';
    const options = {
      routeHost: 'http://url.com',
      serviceUrl: 'http://url.com',
    };
    makeConfigFile.mockResolvedValueOnce(input);
    await actionHandler(input, options);

    expect(makeConfigFile).toHaveBeenCalledWith(input, options);
  });

  it('should display an error', async () => {
    const spy = jest.spyOn(console, 'error');
    makeConfigFile.mockRejectedValueOnce(new Error('Unable to make file'));

    await actionHandler('bad input');

    expect(process.exitCode).toEqual(1);
    expect(spy).toHaveBeenCalledWith(
      chalk.bold.red`x Error:` + ' Unable to make file'
    );
  });
});
