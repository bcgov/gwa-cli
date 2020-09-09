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
  // const { isFocused } = useFocus();
  const index = ids.indexOf(match.params.plugin);
  const total = ids.length;

  const onSubmit = (formData: any) => {
    setPlugin((prev) => ({
      ...prev,
      [plugin.id]: {
        ...prev[plugin.id],
        data: {
          ...prev[plugin.id].data,
          config: merge(prev[plugin.id].data.config, formData, {
            arrayMerge: (destArr, srcArr) => srcArr,
          }),
        },
      },
    }));
    setSaved(true);
    setTimeout(() => setSaved(false), 1500);
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
  const onEncrypt = (key: string) => {
    setPlugin((prev) => ({
      ...prev,
      [plugin.id]: {
        ...prev[plugin.id],
        encrypted: prev[plugin.id].encrypted.includes(key)
          ? prev[plugin.id].encrypted.filter((k) => k !== key)
          : [...prev[plugin.id].encrypted, key],
      },
    }));
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
      <Box padding={1} marginY={1} justifyContent="space-between">
        <Text>{plugin.description}</Text>
      </Box>
      <Form
        encryptedFields={plugin.encrypted}
        constraints={plugin.constraints}
        data={plugin.data.config}
        onEncrypt={onEncrypt}
        onSubmit={onSubmit}
      />
      <Box justifyContent={saved ? 'space-between' : 'flex-end'} marginTop={2}>
        {saved && <Text color="greenBright">Settings Saved</Text>}
        <Box justifyContent="flex-end">
          <Text dimColor>Prev [ctrl+p] / Next [ctrl+n]</Text>
        </Box>
      </Box>
    </Box>
  );
};

export default PluginEditor;
