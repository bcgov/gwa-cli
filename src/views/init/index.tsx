import React, { useContext, useState, useEffect } from 'react';
import { Box, Text, useInput, render } from 'ink';
import Spinner from 'ink-spinner';
import path from 'path';

import AppContext from '../../services/context';
import PromptForm from '../../components/prompt-form';
import TextField from '../../components/prompt-form/text-field';
// import { parseLocalFile, parseYaml } from '../../services/kong';

type FormData = {
  name: string;
  specUrl: string;
  file: string;
};

interface ConfigOrgProps {
  source: string;
}

const ConfigOrg: React.FC<ConfigOrgProps> = ({ source }) => {
  const { dir } = useContext(AppContext);
  const [isProcessing, setProcessing] = useState<boolean>(false);
  const [processError, setProcessError] = useState<string | null>(null);
  const [valid, setValid] = useState<boolean>(false);
  const [formData, setFormData] = useState<FormData>({
    name: '',
    specUrl: '',
    file: '',
  });

  const onChange = (name: string, value: string) => {
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };
  const onSubmit = () => {
    setProcessing(true);
    try {
      //await parseLocalFile(
      //  formData.specUrl || path.resolve(dir, source),
      //  formData.name,
      //  dir + formData.file
      //);
      setProcessing(false);
      setValid(true);
    } catch (err) {
      setValid(false);
      setProcessing(false);
      setProcessError(err.message);
    }
  };

  return (
    <Box flexDirection="column">
      <Box marginBottom={1}>
        <Text>Follow the prompts to generate a config file</Text>
      </Box>
      <PromptForm complete={valid} onSubmit={onSubmit}>
        <TextField
          required
          label="Organization Name"
          name="name"
          onChange={onChange}
        />
        <TextField
          required
          label="Config File Name"
          type="text"
          placeholder="Enter a name for the YAML config file"
          name="file"
          onChange={onChange}
        />
      </PromptForm>
      {isProcessing && (
        <Box>
          <Text>
            <Spinner /> Processing OpenAPI Spec...
          </Text>
        </Box>
      )}
      {processError && (
        <Box borderColor="redBright" borderStyle="round" marginY={1}>
          <Text color="redBright">{processError}</Text>
        </Box>
      )}
      {valid && (
        <Box borderColor="greenBright" borderStyle="round" marginY={1}>
          <Text color="greenBright">
            Group config saved! Press [q] or [ctrl + c] to close.
          </Text>
        </Box>
      )}
    </Box>
  );
};

export default (args: any) =>
  render(<ConfigOrg source={args.f || args.file} />);
