import reducer, { makeInitialState } from '../reducer';

describe('components/prompt-form/reducer', () => {
  it('should make an initial state object', () => {
    const options = ['option 1', 'option 2'];
    expect(makeInitialState(options)).toEqual({
      step: 0,
      error: undefined,
      value: '',
      data: {},
      done: false,
      prompts: options,
    });
  });

  it('should throw if no type is set', () => {
    expect(() => reducer(makeInitialState([]), {})).toThrow();
  });
});
