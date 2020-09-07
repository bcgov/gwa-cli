import React, { useState } from 'react';
import {
  Box,
  Text,
  Transform,
  Newline,
  useFocus,
  useFocusManager,
  useInput,
} from 'ink';
import merge from 'deepmerge';
import { RouteComponentProps } from 'react-router';
import TextInput from 'ink-text-input';
import { ValidateJS } from 'validate.js';

import Checkbox from '../../components/form/checkbox';
import Form from '../../components/form';
import { IPlugin } from '../../types';
import { pluginsState } from '../../state/plugins';

interface PluginEditorProps extends RouteComponentProps<{ plugin: string }> {}

const PluginEditor: React.FC<PluginEditorProps> = ({ match }) => {
  const [plugins, setPlugin] = pluginsState.use();
  const [saved, setSaved] = useState<boolean>(false);
  const ids = Object.keys(plugins);
  const plugin: IPlugin = plugins[match.params.plugin];
  const { focusNext } = useFocusManager();
  const { isFocused } = useFocus();
  const index = ids.indexOf(match.params.plugin);
  const total = ids.length;

  const onSubmit = (formData: any) => {
    setPlugin((prev) => ({
      ...prev,
      [plugin.id]: {
        ...prev[plugin.id],
        data: {
          ...prev[plugin.id].data,
          config: merge(prev[plugin.id].data.config, formData),
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

  useInput((input, key) => {
    if (input === 'e' && key.ctrl) {
      onToggleEnabled(!plugin.data.enabled);
    }
    /* if (!isFocused && key.return) {
     *   focusNext();
     * } */
  });

  return (
    <Box flexDirection="column">
      <Box marginY={1} justifyContent="space-between">
        <Text>Form Header</Text>
      </Box>
      <Form
        constraints={plugin.constraints}
        data={plugin.data.config}
        onSubmit={onSubmit}
      />
      {saved && <Text color="greenBright">Settings Saved</Text>}
      <Box justifyContent="flex-end">
        <Text dimColor>Prev [ctrl+p] / Next [ctrl+n]</Text>
      </Box>
    </Box>
  );
};

export default PluginEditor;
