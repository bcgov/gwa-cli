import { render } from 'ink-testing-library';
import ink from 'ink';

import renderer from '../renderer';

jest.mock('../upload-view');
jest.mock('ink', () => ({
  render: jest.fn(),
}));

describe('commands/publish/ui', () => {
  it('should render the init UI', () => {
    const options = {
      input: 'input',
      dryRun: 'true',
    };
    renderer('dev');
    expect(ink.render).toHaveBeenCalled();
  });
});
