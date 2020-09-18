export interface IAppContext {
  dir: string;
  file: string;
  version: string;
}

export interface IPlugin {
  id: string;
  name: string;
  description: string;
  constraints: any;
  encrypted: string[];
  data: {
    name: string;
    enabled: boolean;
    config: any;
  };
}
