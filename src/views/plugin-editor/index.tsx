import * as React from 'react';
import { Box, Text, Newline, useFocus, useFocusManager, useInput } from 'ink';
import { match } from 'react-router';
import TextInput from 'ink-text-input';
import { ValidateJS } from 'validate.js';

import Checkbox from '../../components/form/checkbox';
import Form from '../../components/form';
import { IPlugin } from '../../types';
import { pluginsState } from '../../state/plugins';

interface PluginEditorProps {
  match: match;
}

const PluginEditor: React.FC<PluginEditorProps> = ({ match }) => {
  const [plugins, setPlugin] = pluginsState.use();
  const plugin: IPlugin = plugins[match.params.plugin];
  const { focusNext } = useFocusManager();
  const { isFocused } = useFocus();

  useInput((input, key) => {
    if (!isFocused && key.return) {
      focusNext();
    }
  });

  const onSubmit = (formData: any) => {
    /* setPlugin((prev) => ({
     *   ...prev,
     *   [selected]: {
     *     ...prev[selected],
     *     config: formData,
     *   },
     * })); */
  };
  const onToggleEnabled = (value: boolean) => {
    /* setPlugin((prev) => {
     *   return {
     *     ...prev,
     *     [selected]: {
     *       ...prev[selected],
     *       enabled: value,
     *     },
     *   };
     * }); */
  };

  return (
    <Box flexDirection="column">
      <Box marginBottom={1}>
        <Checkbox
          checked={plugin.enabled}
          name="enabled"
          label="Plugin Enabled"
          onChange={onToggleEnabled}
        />
      </Box>
      <Form
        constraints={plugin.constraints}
        data={plugin.config}
        onSubmit={onSubmit}
      />
    </Box>
  );
};

export default PluginEditor;
