import { newRidgeState } from 'react-ridge-state';

import { IPlugin } from '../types';
import data from '../data';

export const pluginsState = newRidgeState<{ [prop: string]: IPlugin }>(
  data.plugins
);

export const activePluginState = newRidgeState<string>('bcgov-gwa-endpoint');
