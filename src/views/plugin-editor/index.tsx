import React, { useEffect, useState } from 'react';
import { Box, Text, useFocusManager, useInput } from 'ink';
import merge from 'deepmerge';
import { RouteComponentProps } from 'react-router';
import { clearTimeout, setTimeout } from 'timers';

import { useAppState } from '../../state/app';
import Form from '../../components/form';
import { usePluginsState, set, toggleEnabled } from '../../state/plugins';
import constraints from '../../validators';

interface PluginEditorProps extends RouteComponentProps<{ plugin: string }> {}

const PluginEditor: React.FC<PluginEditorProps> = ({ match }) => {
  const id = match.params.plugin;
  const mode = useAppState((state) => state.mode);
  const toggleMode = useAppState((state) => state.toggleMode);
  const plugins = usePluginsState();
  const [saved, setSaved] = useState<boolean>(false);
  const plugin: any = plugins[id];

  const onSubmit = (formData: any) => {
    set(id, formData);
    setSaved(true);
  };

  const onToggleEnabled = (value: boolean) => toggleEnabled(id, value);
  const onEncrypt = (key: string, isEncrypted: boolean) => {};

  useInput((input, key) => {
    if (mode === 'view') {
      if (input === 'e') {
        onToggleEnabled(!plugin.meta.enabled);
      }
      if (key.return) {
        toggleMode();
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
        <Text>{plugin.meta.description}</Text>
      </Box>
      <Form
        enabled={mode === 'edit'}
        encryptedFields={[]}
        constraints={constraints[id]}
        data={plugin.config}
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
