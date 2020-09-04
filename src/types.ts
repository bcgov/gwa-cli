export interface IAppContext {
  dir: string;
  file: string | null;
}

export interface IPlugin {
  id: string;
  name: string;
  description: string;
  constraints: any;
  data: {
    name: string;
    enabled: boolean;
    config: any;
  };
}
