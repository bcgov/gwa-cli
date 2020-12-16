import * as utils from '../utils';

describe('services/utils', () => {
  describe('#isLocalInput', () => {
    it('should return true for a local path', () => {
      expect(utils.isLocalInput('/test/file.yaml')).toEqual(true);
    });

    it('should return false for URL or any other input', () => {
      expect(utils.isLocalInput('http://test.com/input.yaml')).toEqual(false);
    });
  });

  describe('#makeOutputFilename', () => {
    it('should convert a local file to yaml', () => {
      expect(utils.makeOutputFilename('/path/to/input.json')).toEqual(
        '/path/to/input.yaml'
      );
    });

    it('should convert a url to yaml', () => {
      expect(
        utils.makeOutputFilename('https://test.com/path/to/input.json')
      ).toEqual('input.yaml');
    });

    it('should return outfile', () => {
      expect(
        utils.makeOutputFilename(
          'https://test.com/path/to/input.json',
          'outfile.yaml'
        )
      ).toEqual('outfile.yaml');
    });

    it('should throw an error if the file does not have a JSON extension', () => {
      expect(() => utils.makeOutputFilename('http://test.com/')).toThrow();
    });

    it('should throw errors', () => {
      expect(() => utils.makeOutputFilename()).toThrow();
    });
  });
});
