import validate from 'validate.js';

import prompts from '../prompts';

describe('commmands/init/prompts', () => {
  const prompt = prompts.find((prompt) => prompt.key === 'namespace');
  const spec = {
    namespace: prompt.constraint,
  };

  it('should report invalid namespaces', () => {
    const invalidOptions = [
      'test',
      'testingalongnamespacename',
      'test!@#$%^&*()_+',
    ];
    expect.assertions(invalidOptions.length);
    invalidOptions.forEach((option) => {
      expect(
        validate(
          {
            namespace: option,
          },
          spec
        )
      ).not.toBeUndefined();
    });
  });

  it('should report valid namespaces', () => {
    const validOptions = ['testing', 'test-1', 'test01', 'name-space01asd'];
    expect.assertions(validOptions.length);
    validOptions.forEach((option) => {
      expect(
        validate(
          {
            namespace: option,
          },
          spec
        )
      ).toBeUndefined();
    });
  });
});
