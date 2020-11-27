import ui from '../ui';
import { loadConfig, parseConfig } from '../services/app';
import { initTeamState } from '../state/team';
import { initPluginsState } from '../state/plugins';

export default async function (input: string) {
  try {
    const config = await loadConfig(input);
    const { name, team, host, plugins } = parseConfig(config);

    initTeamState({
      name,
      team,
      host,
    });
    initPluginsState();
    ui();
  } catch (err) {
    console.error(err);
  }
}
