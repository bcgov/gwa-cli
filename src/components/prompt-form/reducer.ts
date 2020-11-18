import produce, { Draft } from 'immer';
import type { Prompt, PromptAction, PromptFormState } from './types';
import validate from 'validate.js';

export const makeInitialState = (prompts: Prompt[] = []): PromptFormState => ({
  step: 0,
  error: undefined,
  value: '',
  data: {},
  prompts,
});

const reducer = produce(
  (draft: Draft<PromptFormState>, action: PromptAction) => {
    switch (action.type) {
      case 'change':
        draft.value = action.payload;
        draft.error = undefined;
        break;

      case 'next':
        const currentPrompt = draft.prompts[draft.step];
        const errors = validate.single(
          action.payload,
          currentPrompt.constraint,
          {
            format: 'flat',
          }
        );

        if (errors) {
          draft.error = errors;
        } else {
          draft.data[currentPrompt.key] = action.payload;
          draft.step = draft.step + 1;
          draft.value = '';
        }
        break;

      case 'reset':
        draft.value = '';
        draft.error = undefined;
        break;

      default:
        throw new Error();
    }
  },
  makeInitialState()
);

export default reducer;
