import ui from '../ui';
import { loadConfig } from '../services/app';
import { initTeamState } from '../state/team';
import { initPluginsState } from '../state/plugins';

export default async function (input: string) {
  try {
    const config = await loadConfig(input);
    initTeamState({
      name: 'name',
      team: 'test',
      host: 'host',
    });
    initPluginsState();
    ui();
  } catch (err) {
    console.error(err);
  }
}
