import { ValidateJS } from 'validate.js';

export interface IAppContext {
  dir: string;
}

export interface IPlugin {
  id: string;
  name: string;
  enabled: boolean;
  config: any;
  constraints: any;
}
