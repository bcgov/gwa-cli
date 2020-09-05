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
  const plugin: IPlugin = plugins[match.params.plugin];
  const { focusNext } = useFocusManager();
  const { isFocused } = useFocus();

  useInput((input, key) => {
    if (!isFocused && key.return) {
      focusNext();
    }
  });

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

  return (
    <Box flexDirection="column" width="100%">
      <Box marginBottom={1} justifyContent="space-between">
        <Box>
          <Text inverse color="yellow">
            {plugin.name}
          </Text>
          <Text inverse color="gray">
            {` 1/3 `}
          </Text>
        </Box>
        <Box>
          <Text>Prev [ctrl+p] / Next [ctrl+n]</Text>
        </Box>
      </Box>
      <Box marginBottom={1}>
        <Checkbox
          autoFocus
          checked={plugin.data.enabled}
          name="enabled"
          label="Plugin Enabled"
          onChange={onToggleEnabled}
        />
      </Box>
      <Form
        constraints={plugin.constraints}
        data={plugin.data.config}
        onSubmit={onSubmit}
      />
      {saved && <Text color="greenBright">Settings Saved</Text>}
    </Box>
  );
};

export default PluginEditor;
