import * as React from 'react';
import { Box, Text } from 'ink';
import has from 'lodash/has';
import TextInput from 'ink-text-input';
import validate from 'validate.js';
import { uid } from 'react-uid';

import reducer, { makeInitialState } from './reducer';
import type { Prompt } from './types';

interface SetupViewProps {
  children: ({ data }: { data: any }) => React.Node;
  options: Prompt[];
  title: string;
}

const SetupView: React.FC<SetupViewProps> = ({ children, options, title }) => {
  const [
    { data, done, error, prompts, step, value },
    dispatch,
  ] = React.useReducer(reducer, makeInitialState(options));
  const prompt = options[step];
  const onSubmit = (value: string) => {
    const errors = validate.single(value, prompt.constraint, {
      format: 'flat',
    });

    if (errors) {
      dispatch({
        type: 'error',
        payload: errors,
      });
    } else {
      dispatch({
        type: 'next',
        payload: {
          [prompt.key]: value,
        },
      });
    }
  };

  React.useEffect(() => {
    if (step === options.length) {
      dispatch({ type: 'done' });
    }
  }, [options, step]);

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
                onSubmit={onSubmit}
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
      {done && children && children({ data })}
    </Box>
  );
};

export default SetupView;
