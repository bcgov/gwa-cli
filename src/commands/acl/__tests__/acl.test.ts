import { actionHandler } from '../acl';
import render from '../renderer';

jest.mock('../renderer');

describe('commands/acl', () => {
  afterEach(() => {
    render.mockReset();
  });

  it('should adds managers and users', () => {
    actionHandler({
      managers: ['admin1@idir'],
      users: ['user1@idir', 'user2@idir'],
      debug: true,
    });
    expect(render).toHaveBeenCalledWith(
      expect.arrayContaining([
        {
          username: 'user1@idir',
          roles: ['viewer'],
        },
        {
          username: 'user1@idir',
          roles: ['viewer'],
        },
        {
          username: 'admin1@idir',
          roles: ['viewer'],
        },
        {
          username: 'admin1@idir',
          roles: ['admin'],
        },
      ]),
      true
    );
  });

  it('should only send users', () => {
    actionHandler({
      users: ['user1@idir'],
    });
    expect(render).toHaveBeenCalledWith(
      expect.arrayContaining([
        {
          username: 'user1@idir',
          roles: ['viewer'],
        },
      ]),
      undefined
    );
  });

  it('should be able to handle just one set of role arguments', () => {
    expect(() => {
      actionHandler({
        managers: ['admin@idir'],
      });
    }).not.toThrow();
  });
});
