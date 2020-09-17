import { IPlugin } from '../../types';
import { ISpecExpose } from './types';

const constraints = {
  spec_url: {
    type: 'string',
    presence: {
      allowempty: false,
    },
    url: true,
  },
};

const plugin: IPlugin = {
  id: 'kong-spec-expose',
  name: 'Kong Spec Expose',
  description:
    'This plugin will expose the OpenAPI Spec (OAS), Swagger, or other specification of auth protected API services fronted by the Kong gateway.',
  constraints,
  encrypted: [],
  data: {
    name: 'kong-spec-expose',
    enabled: false,
    config: {
      spec_url: '',
    },
  },
};

export default plugin;
