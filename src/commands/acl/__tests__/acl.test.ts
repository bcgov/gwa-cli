jest.mock('../renderer');
import chalk from 'chalk';

import { actionHandler } from '../acl';
import render from '../renderer';

describe('commands/acl', () => {
  it('should adds managers and users', () => {
    actionHandler({
      managers: ['admin1@idir'],
      users: ['user1@idir', 'user2@idir'],
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
          roles: ['admin'],
        },
      ])
    );
  });

  it('should be able to handle just one set of role arguments', () => {
    expect(() => {
      actionHandler({
        managers: ['admin@idir'],
      });
    }).not.toThrow();
  });

  it('should not add the same user to both managers and users', () => {
    expect(() => {
      actionHandler({
        managers: ['admin@idir'],
        users: ['admin@idir'],
      });
    }).toThrow();
  });
});
