export interface AppContext {
  dir: string;
  file?: string;
  version: string;
}

export interface PluginMeta {
  name: string;
  url: string;
  bcgov: boolean;
  description: string;
}

export interface PluginConfig {
  name: string;
  enabled: boolean;
  tags: string[];
  config: any;
}

export interface PluginsResult {
  [id: string]: {
    meta: PluginMeta;
    config: PluginConfig;
  };
}
