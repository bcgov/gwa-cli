import validate from 'validate.js';

import prompts from '../prompts';

describe('commmands/init/prompts', () => {
  const prompt = prompts.find((prompt) => prompt.key === 'outfile');
  const spec = {
    outfile: prompt.constraint,
  };

  it('should invalidate bad outfiles', () => {
    const invalidOptions = ['file', 'file.json'];
    expect.assertions(invalidOptions.length);
    invalidOptions.forEach((option) => {
      expect(
        validate(
          {
            outfile: option,
          },
          spec
        )
      ).not.toBeUndefined();
    });
  });

  it('should validate bad outfiles', () => {
    const validOptions = ['file.yaml', 'fileyaml.yaml'];
    expect.assertions(validOptions.length);
    validOptions.forEach((option) => {
      expect(
        validate(
          {
            outfile: option,
          },
          spec
        )
      ).toBeUndefined();
    });
  });
});
