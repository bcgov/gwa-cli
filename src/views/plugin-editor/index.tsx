import React, { useEffect, useState } from 'react';
import { Box, Text, useFocusManager, useInput } from 'ink';
import merge from 'deepmerge';
import { RouteComponentProps } from 'react-router';
import { clearTimeout, setTimeout } from 'timers';

import { appState } from '../../state/app';
import Form from '../../components/form';
import { IPlugin } from '../../types';
import { encryptedValues, pluginsState } from '../../state/plugins';

interface PluginEditorProps extends RouteComponentProps<{ plugin: string }> {}

const PluginEditor: React.FC<PluginEditorProps> = ({ match }) => {
  const mode = appState.useSelector((state) => state.mode);
  const [encrypted, setEncrypted] = encryptedValues.use();
  const [state, setAppState] = appState.use();
  const [plugins, setPlugin] = pluginsState.use();
  const [saved, setSaved] = useState<boolean>(false);
  const plugin: IPlugin = plugins[match.params.plugin];

  const onSubmit = (formData: any) => {
    setPlugin((prev) => ({
      ...prev,
      [plugin.id]: {
        ...prev[plugin.id],
        data: {
          ...prev[plugin.id].data,
          config: merge(prev[plugin.id].data.config, formData, {
            arrayMerge: (_, srcArr) => srcArr,
          }),
        },
      },
    }));
    setSaved(true);
  };

  const onToggleEnabled = (value: boolean) => {
    setPlugin((prev) => {
      return {
        ...prev,
        [plugin.id]: {
          ...prev[plugin.id],
          data: {
            ...prev[plugin.id].data,
            enabled: value,
          },
        },
      };
    });
  };
  const onEncrypt = (key: string, isEncrypted: boolean) => {
    setEncrypted((prev) => {
      if (isEncrypted) {
        return [...prev, key];
      }
      return prev.filter((v) => v !== key);
    });
  };

  useInput((input, key) => {
    if (state.mode === 'view') {
      if (input === 'e') {
        onToggleEnabled(!plugin.data.enabled);
      }
      if (key.return) {
        setAppState((prev) => ({
          ...prev,
          mode: 'edit',
        }));
      }
    }
  });

  useEffect(() => {
    let timer: NodeJS.Timer;

    if (saved) {
      timer = setTimeout(() => {
        setSaved(false);
      }, 1500);
    }

    return () => {
      clearTimeout(timer);
    };
  }, [saved]);

  return (
    <Box flexDirection="column">
      <Box padding={1} marginY={1} justifyContent="space-between">
        <Text>{plugin.description}</Text>
      </Box>
      <Form
        enabled={state.mode === 'edit'}
        encryptedFields={encrypted}
        constraints={plugin.constraints}
        data={plugin.data.config}
        onEncrypt={onEncrypt}
        onSubmit={onSubmit}
      />
      <Box justifyContent="space-between" marginTop={2}>
        <Box>
          <Text bold inverse>
            {' '}
            {mode.toUpperCase()}{' '}
          </Text>
          <Text inverse color="magenta">
            {` [${mode === 'edit' ? 'ESC to save' : 'ENTER to edit'}] `}
          </Text>
          {saved && <Text color="greenBright">{' Settings Saved '}</Text>}
        </Box>
        <Box justifyContent="flex-end">
          {mode === 'view' && <Text dimColor>[p] Prev / Next [n]</Text>}
        </Box>
      </Box>
    </Box>
  );
};

export default PluginEditor;
