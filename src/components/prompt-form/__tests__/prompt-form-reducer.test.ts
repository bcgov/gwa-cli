import reducer, { makeInitialState } from '../reducer';
import type { PromptAction, PromptFormState } from '../types';

const prompts = [
  {
    label: 'Required',
    key: 'required',
    constraint: {
      presence: {
        allowEmpty: false,
      },
    },
  },
  {
    label: 'Optional Prompt',
    key: 'optional',
    constraint: {},
  },
  {
    label: 'Secret Prompt',
    key: 'secret',
    secret: true,
    constraint: {},
  },
];

describe('components/prompt-form/reducer', () => {
  it('should make an initial state object', () => {
    const options = [];
    expect(makeInitialState(options)).toEqual({
      step: 0,
      error: undefined,
      value: '',
      data: {},
      prompts: options,
    });
  });

  it('should throw if no type is set', () => {
    // @ts-ignore
    expect(() => reducer(makeInitialState([]), {})).toThrow();
  });

  it('should handle change action', () => {
    const action: PromptAction = {
      type: 'change',
      payload: 'a',
    };
    expect(reducer(makeInitialState(prompts), action)).toEqual(
      expect.objectContaining({
        value: 'a',
      })
    );
  });

  it('should handle valid next action', () => {
    const action: PromptAction = {
      type: 'next',
      payload: 'required',
    };
    const initialState = {
      ...makeInitialState(prompts),
      step: 0,
      value: 'hello',
      data: {},
    };
    expect(reducer(initialState, action)).toEqual(
      expect.objectContaining({
        value: '',
        step: 1,
        data: {
          required: 'required',
        },
      })
    );
  });

  it('should handle an invalid next action', () => {
    const action: PromptAction = {
      type: 'next',
      payload: '',
    };
    const initialState = {
      ...makeInitialState(prompts),
      step: 0,
      value: '',
      data: {},
    };
    expect(reducer(initialState, action)).toEqual(
      expect.objectContaining({
        value: '',
        step: 0,
        data: {},
        error: ["can't be blank"],
      })
    );
  });

  it('should handle reset action', () => {
    const action: PromptAction = {
      type: 'reset',
    };
    const initialState: PromptFormState = {
      ...makeInitialState(prompts),
      value: 'hello',
      error: ['err'],
    };
    expect(reducer(initialState, action)).toEqual(
      expect.objectContaining({
        value: '',
        error: undefined,
      })
    );
  });
});
