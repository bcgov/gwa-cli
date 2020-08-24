import { ISpecExpose } from './types';

const constraints = {
  spec_url: {
    type: 'string',
    presence: true,
    url: true,
  },
};

const plugin: IPlugin = {
  id: 'kong-spec-expose',
  name: 'kong-spec-expose',
  enabled: false,
  constraints,
  config: {
    spec_url: '',
  },
};

export default plugin;
