jest.mock('../renderer');

import { actionHandler } from '../acl';
import render from '../renderer';

describe('commands/acl', () => {
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

  it('should be able to handle just one set of role arguments', () => {
    expect(() => {
      actionHandler({
        managers: ['admin@idir'],
      });
    }).not.toThrow();
  });
});
