import { exportConfig } from '../../../services/app';
import { fetchSpec, importSpec } from '../../../services/openapi';
import { generateConfig } from '../../../services/kong';
import { importInput, makeConfigFile, parseOptions } from '../create-actions';
import {
  initPluginsState,
  generatePluginTemplates,
} from '../../../state/plugins';
jest.mock('../../../services/kong');
jest.mock('../../../services/plugins');
jest.mock('../../../state/plugins');
jest.mock('../../../services/app');
jest.mock('../../../services/openapi');
const CACHED_ENV = process.env;

describe('commands/create/create-actions', () => {
  beforeAll(() => {
    process.env = { ...CACHED_ENV, GWA_NAMESPACE: 'sampler' };
  });

  afterAll(() => {
    process.env = CACHED_ENV;
    jest.resetAllMocks();
  });

  it('should call loadPlugins with dir', () => {
    makeConfigFile('file.json', {});
    expect(initPluginsState).toHaveBeenCalled();
  });

  it('should throw if no arguments are provided', () => {
    expect(async () => await makeConfigFile()).rejects.toThrow();
  });

  it('should throw if errors are encountered', () => {
    exportConfig.mockRejectedValueOnce();
    expect(async () => await makeConfigFile('file.json', {})).rejects.toThrow();
  });

  //it('should return outfile name on success', () => {
  //  initPluginsState.mockResolvedValueOnce(true);
  //  fetchSpec.mockResolvedValueOnce('config');
  //  generateConfig.mockResolvedValue('config');
  //  expect(
  //    makeConfigFile('http://testerrrr.com/file.json', {
  //      outfile: 'file.yaml',
  //      plugins: [],
  //      options: {
  //        routeHost: 'host',
  //        serviceUrl: 'url',
  //      },
  //    })
  //  ).resolves.toBe('file.yaml');
  //  expect(exportConfig).toHaveBeenCalledWith('config', 'file.yaml');
  //});

  describe('#importInput', () => {
    it('should call importInput when local file is detected', () => {
      importInput({
        namespace: 'sampler',
        input: 'file.json',
      });
      expect(importSpec).toHaveBeenCalled();
    });

    it('should call fetchSpec when URL is detected', () => {
      importInput({
        namespace: 'sampler',
        input: 'http://test.com/file.json',
      });
      expect(fetchSpec).toHaveBeenCalled();
    });

    it('should throw on failure', () => {
      fetchSpec.mockRejectedValueOnce();

      expect(
        async () =>
          await importInput({
            namespace: 'sampler',
            input: 'http://test.com/file.json',
          })
      ).rejects.toThrow();
    });
  });

  describe('#parsePlugins', () => {
    const options = {
      namespace: 'sampler',
    };
    it('should parse array of plugins', () => {
      parseOptions({ ...options, plugins: ['acl', 'oidc'] });
      expect(generatePluginTemplates).toHaveBeenCalledWith(
        ['acl', 'oidc'],
        'sampler'
      );
    });

    it('should parse a string of space separated plugins', () => {
      parseOptions({ ...options, plugins: 'acl oidc' });

      expect(generatePluginTemplates).toHaveBeenCalledWith(
        ['acl', 'oidc'],
        'sampler'
      );
    });

    it('should parse a string of comma separated plugins', () => {
      parseOptions({ ...options, plugins: 'acl,oidc' });

      expect(generatePluginTemplates).toHaveBeenCalledWith(
        ['acl', 'oidc'],
        'sampler'
      );
    });

    it('should parse a string of space and comma separated plugins', () => {
      parseOptions({ ...options, plugins: 'acl, oidc,rate-limiting' });

      expect(generatePluginTemplates).toHaveBeenCalledWith(
        ['acl', 'oidc', 'rate-limiting'],
        'sampler'
      );
    });
  });
});
