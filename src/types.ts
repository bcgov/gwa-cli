export type Envs = 'dev' | 'prod' | 'test';

export interface AppContext {
  dir: string;
  file?: string;
  version: string;
}

export interface PluginMeta {
  name: string;
  url: string;
  id: string;
  author: string;
  description: string;
  enabled: boolean;
}

export interface PluginConfig {
  name: string;
  enabled: boolean;
  tags: string[];
  config: any;
}

export interface PluginObject {
  meta: PluginMeta;
  config: PluginConfig;
}

export interface PluginsResult {
  [id: string]: PluginObject;
}

export interface InitOptions {
  namespace: string;
  clientId: string;
  clientSecret: string;
  dataCenter: string;
  env: string;
  dev: boolean;
  test: boolean;
  prod: boolean;
  debug: boolean;
  apiVersion?: string;
  dirApiVersion?: string;
}
