import * as React from 'react';
import { Box, Text, useApp } from 'ink';
import has from 'lodash/has';
import Spinner from 'ink-spinner';
import TextInput from 'ink-text-input';
import validate from 'validate.js';
import { uid } from 'react-uid';

import { convertRemote } from '../../services/kong';
import { exportConfig } from '../../services/app';
import { fetchSpec } from '../../services/openapi';
import { generatePluginTemplates } from '../../state/plugins';
import reducer, { initialState } from './reducer';

interface SetupViewProps {}

const SetupView: React.FC<SetupViewProps> = () => {
  const { exit } = useApp();
  const [
    { data, done, error, prompts, specError, status, step, value },
    dispatch,
  ] = React.useReducer(reducer, initialState);
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
          [prompt.key]: value,
        },
      });
    }
  };

  React.useEffect(() => {
    const loadSpec = async () => {
      dispatch({ type: 'spec/loading' });
      try {
        const json = await fetchSpec(data.url);
        const plugins = generatePluginTemplates(
          data.plugins.split(', '),
          data.team
        );
        const result = await convertRemote(json, data.team, plugins);
        dispatch({ type: 'spec/success', payload: result });
        await exportConfig(result, data.outfile);
        dispatch({ type: 'spec/written', payload: result });
      } catch (err) {
        dispatch({ type: 'spec/failed', payload: err.message });
        exit();
      }
    };
    if (!prompt) {
      loadSpec();
    }
  }, [data.team, data.url, fetchSpec, prompt, dispatch]);

  React.useEffect(() => {
    if (done) {
      exit();
    }
  }, [done, exit]);

  return (
    <Box flexDirection="column">
      <Box marginY={1}>
        <Text>Fill in the prompts to build your configuration file</Text>
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
                <Text>{data[d.key]}</Text>
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
      {status === 'success' && (
        <Box flexDirection="column">
          <Box>
            <Box>
              <Text bold color="green">
                ✓
              </Text>
            </Box>
            <Box marginX={1}>
              <Text bold color="green">
                OpenAPI spec imported
              </Text>
            </Box>
          </Box>
          {done && (
            <Box>
              <Box>
                <Text bold color="green">
                  ✓
                </Text>
              </Box>
              <Box marginX={1}>
                <Text bold color="green">
                  {`File ${data.outfile} written`}
                </Text>
              </Box>
            </Box>
          )}
        </Box>
      )}
      {status === 'loading' && (
        <Box>
          <Text>
            <Text color="green">
              <Spinner />
            </Text>{' '}
            Loading Spec...
          </Text>
        </Box>
      )}
      {status === 'failed' && (
        <Box flexDirection="column">
          <Box>
            <Text>
              <Text bold color="red">
                x
              </Text>{' '}
              Failed to load OpenAPI spec
            </Text>
          </Box>
          <Box borderColor="red" borderStyle="round" marginTop={1}>
            <Text>{specError}</Text>
          </Box>
        </Box>
      )}
    </Box>
  );
};

export default SetupView;
