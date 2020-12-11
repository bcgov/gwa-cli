import fs from 'fs';

import * as app from '../app';

jest.mock('fs', () => ({
  existsSync: jest.fn(),
  promises: {
    readFile: jest.fn(),
    writeFile: jest.fn(),
  },
}));

describe('services/app', () => {
  afterEach(() => {
    fs.existsSync.mockReset();
    fs.promises.writeFile.mockReset();
  });

  describe('#checkForEnvFile', () => {
    it('should return false for missing', () => {
      fs.existsSync.mockReturnValue(false);
      expect(app.checkForEnvFile()).toEqual(false);
    });

    it('should return false for missing', () => {
      fs.existsSync.mockReturnValue(true);
      expect(app.checkForEnvFile()).toEqual(true);
    });
  });

  describe('#makeEnvFile', () => {
    it('should write an env file', async () => {
      fs.promises.writeFile.mockResolvedValue(true);
      expect(
        app.makeEnvFile({
          namespace: 'sampler',
          clientId: 'id',
          clientSecret: 'secret',
          env: 'dev',
        })
      ).resolves.toEqual('.env file successfully generated');
    });

    it('should throw', async () => {
      fs.promises.writeFile.mockRejectedValue('err');
      await expect(
        app.makeEnvFile({
          namespace: 'sampler',
          clientId: 'id',
          clientSecret: 'secret',
          env: 'dev',
        })
      ).rejects.toThrow('Unable to write file err');
    });
  });

  describe('#loadConfig', () => {
    it('should load a parsed JSON object from YAML', async () => {
      fs.promises.readFile.mockResolvedValue(`
prop1: string
prop2:
  nested: true
arr:
  - arr`);
      await expect(app.loadConfig('config.yaml')).resolves.toEqual({
        prop1: 'string',
        prop2: {
          nested: true,
        },
        arr: ['arr'],
      });
      expect(fs.promises.readFile).toHaveBeenCalledWith(
        expect.stringContaining('/config.yaml'),
        'utf8'
      );
    });

    it('should throw an error', () => {
      expect(async () => await app.loadConfig()).rejects.toThrow(
        'The "path" argument must be of type string. Received undefined'
      );
    });
  });

  describe('#saveConfig', () => {
    it('should write a config file', async () => {
      fs.promises.writeFile.mockResolvedValue(true);
      await app.saveConfig('test', 'config.yaml');
      expect(fs.promises.writeFile).toHaveBeenCalledWith(
        expect.stringContaining('/config.yaml'),
        'test'
      );
    });

    it('should throw an error', async () => {
      // prettier-ignore
      expect(async () => await app.saveConfig()).rejects.toThrow("The \"path\" argument must be of type string. Received undefined");
    });

    it('should throw a writeFile error', async () => {
      fs.promises.writeFile.mockRejectedValue('err');

      await expect(app.saveConfig('test', 'config.yaml')).rejects.toThrow();
    });
  });
});
