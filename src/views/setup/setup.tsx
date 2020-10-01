import * as React from 'react';
import { Box, Text } from 'ink';
import TextInput from 'ink-text-input';
import validate from 'validate.js';
import { uid } from 'react-uid';

import reducer, { initialState } from './reducer';

interface SetupViewProps {}

const SetupView: React.FC<SetupViewProps> = () => {
  const [{ data, error, prompts, step, value }, dispatch] = React.useReducer(
    reducer,
    initialState
  );
  const prompt = prompts[step];
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
          ...prompt,
          value,
        },
      });
    }
  };

  React.useEffect(() => {
    if (!prompt) {
      console.log('checking...');
    }
  }, [prompt]);

  return (
    <Box flexDirection="column">
      <Box marginY={1}>
        <Text>Fill in the prompts to build your configuration file</Text>
      </Box>
      <Box flexDirection="column">
        {data.map((d) => (
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
              <Text>{d.value}</Text>
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
    </Box>
  );
};

export default SetupView;
