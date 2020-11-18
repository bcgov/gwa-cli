import * as React from 'react';
import { Box, Text } from 'ink';
import has from 'lodash/has';
import TextInput from 'ink-text-input';
import validate from 'validate.js';
import { uid } from 'react-uid';

import reducer, { makeInitialState } from './reducer';
import type { Prompt } from './types';

interface SetupViewProps {
  // TODO: return generic here
  onSubmit: (data: any) => void;
  options: Prompt[];
  title: string;
}

const SetupView: React.FC<SetupViewProps> = ({ onSubmit, options, title }) => {
  const [state, dispatch] = React.useReducer(
    reducer,
    makeInitialState(options)
  );
  const { data, error, prompts, step, value } = state;
  const done = step >= options.length;
  const prompt = options[step];
  const onInputSubmit = (value: string) => {
    dispatch({
      type: 'next',
      payload: value,
    });
  };

  React.useEffect(() => {
    if (done && options.length > 0) {
      onSubmit(data);
    }
  }, [done, onSubmit, options.length]);

  return (
    <Box flexDirection="column">
      <Box marginY={1}>
        <Text>{title}</Text>
      </Box>
      <Box flexDirection="column">
        {prompts
          .filter((d: any) => has(data, d.key))
          .map((d: any) => (
            <Box key={uid(d)}>
              <Box>
                <Text bold color="green">
                  ✓
                </Text>
              </Box>
              <Box marginX={1}>
                <Text bold>{d.label}</Text>
              </Box>
              <Box>
                <Text>{d.secret ? '**********' : data[d.key]}</Text>
              </Box>
            </Box>
          ))}
      </Box>
      {prompt && (
        <Box flexDirection="column">
          <Box>
            <Box>
              <Text bold color="green">
                ?
              </Text>
            </Box>
            <Box marginX={1}>
              <Text bold>{prompt.label}</Text>
            </Box>
            <Box>
              <TextInput
                value={value}
                mask={prompt.secret ? '*' : ''}
                onChange={(value) =>
                  dispatch({ type: 'change', payload: value })
                }
                onSubmit={onInputSubmit}
              />
              {error && (
                <Box>
                  <Text color="red">{`<-- ${error}`}</Text>
                </Box>
              )}
            </Box>
          </Box>
        </Box>
      )}
    </Box>
  );
};

export default SetupView;
