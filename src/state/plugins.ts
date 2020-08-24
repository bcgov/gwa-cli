import { newRidgeState } from 'react-ridge-state';

import data from '../data';

export const pluginsState = newRidgeState<{ [string]: IPlugin }>(data.plugins);

export const activePluginState = newRidgeState<string>('bcgov-gwa-endpoint');
