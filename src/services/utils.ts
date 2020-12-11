import validate from 'validate.js';

export const isLocalInput = (input: string): boolean => {
  const errors = validate.single(input, { url: true });
  return Boolean(errors);
};

export const makeOutputFilename = (
  input: string,
  outfile: string = ''
): string => {
  if (!outfile.trim()) {
    if (isLocalInput(input)) {
      return input.replace(/json$/i, 'yaml');
    } else {
      const urlFileName = input.match(/[^\/]+(?=\/$|$)/);
      if (urlFileName) {
        return urlFileName[0].replace('json', 'yaml');
      } else {
        throw 'No outfile specified';
      }
    }
  }

  return outfile;
};
