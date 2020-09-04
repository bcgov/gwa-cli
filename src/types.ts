export interface IAppContext {
  dir: string;
}

export interface IPlugin {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
  config: any;
  constraints: any;
}
