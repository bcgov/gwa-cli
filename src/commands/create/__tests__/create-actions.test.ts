import { loadPlugins } from '../../../services/plugins';
import { makeConfigFile, parsePlugins } from '../create-actions';
import {
  initPluginsState,
  generatePluginTemplates,
} from '../../../state/plugins';

jest.mock('../../../services/plugins');
jest.mock('../../../state/plugins');
const CACHED_ENV = process.env;

describe('commands/create/create-actions', () => {
  beforeEach(() => {
    process.env = { ...CACHED_ENV, GWA_NAMESPACE: 'sampler' };
  });

  afterEach(() => {
    process.env = CACHED_ENV;
    jest.resetAllMocks();
  });

  it('should call loadPlugins with dir', () => {
    makeConfigFile('file.json', {});
  });

  it('should throw if no arguments are provided', () => {
    expect(async () => await makeConfigFile()).rejects.toThrow();
  });

  describe('#parsePlugins', () => {
    it('should parse array of plugins', () => {
      parsePlugins(['acl', 'oidc'], 'sampler');
      expect(generatePluginTemplates).toHaveBeenCalledWith(
        ['acl', 'oidc'],
        'sampler'
      );
    });

    it('should parse a string of space separated plugins', () => {
      parsePlugins('acl oidc', 'sampler');

      expect(generatePluginTemplates).toHaveBeenCalledWith(
        ['acl', 'oidc'],
        'sampler'
      );
    });

    it('should parse a string of comma separated plugins', () => {
      parsePlugins('acl,oidc', 'sampler');

      expect(generatePluginTemplates).toHaveBeenCalledWith(
        ['acl', 'oidc'],
        'sampler'
      );
    });

    it('should parse a string of space and comma separated plugins', () => {
      parsePlugins('acl, oidc,rate-limiting', 'sampler');

      expect(generatePluginTemplates).toHaveBeenCalledWith(
        ['acl', 'oidc', 'rate-limiting'],
        'sampler'
      );
    });
  });
});
