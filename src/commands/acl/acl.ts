import { Command } from 'commander';
import compact from 'lodash/compact';
import intersection from 'lodash/intersection';

import render from './renderer';

type AclOptions = {
  debug: boolean;
  managers: string[];
  users: string[];
};

export const actionHandler = ({ managers, users, debug }: AclOptions) => {
  const duplicates = intersection(managers, users);

  if (duplicates.length > 0) {
    process.exitCode = 1;
    throw new Error('Unable to add duplicate accounts to users and managers');
  }

  const usersToAdd = users.map((username: string) => ({
    username,
    roles: ['viewer'],
  }));
  const managersToAdd = managers.map((username: string) => ({
    username,
    roles: ['admin'],
  }));
  const data = compact([...managersToAdd, ...usersToAdd]);

  render(data);
};

const acl = new Command('acl');

acl
  .description(
    'Update the full membership. Note that this command will overwrite the remote list of users, use with caution'
  )
  .option('-u, --users <users...>', 'Users to add')
  .option('-m, --managers <managers...>', 'Managers to add')
  .option('--debug')
  .action(actionHandler);

export default acl;
