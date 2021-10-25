import { actionHandler } from '../publish-gateway';
import render from '../renderer';

jest.mock('../renderer');

describe('commands/publish', () => {
  it('should pass the options to render correctly', () => {
    actionHandler('input', {
      dryRun: true,
    });

    expect(render).toHaveBeenCalledWith({
      configFile: 'input',
      dryRun: 'true',
    });
  });

  it('should handle empty options object', () => {
    actionHandler('input');

    expect(render).toHaveBeenCalledWith({
      configFile: 'input',
      dryRun: 'false',
    });
  });
});
