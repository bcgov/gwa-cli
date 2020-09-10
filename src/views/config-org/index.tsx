import * as React from 'react';
import { Box, Text, useInput } from 'ink';
import Spinner from 'ink-spinner';
import { useHistory } from 'react-router';

import PromptForm from '../../components/prompt-form';
import TextField from '../../components/prompt-form/text-field';
import StepHeader from '../../components/step-header';
import { orgState, OrgState } from '../../state/org';
import { parseYaml } from '../../services/kong';
import { specState } from '../../state/spec';

type FormData = Omit<OrgState, 'host'>;

interface ConfigOrgProps {}

const ConfigOrg: React.FC<ConfigOrgProps> = ({}) => {
  const history = useHistory();
  const [isProcessing, setProcessing] = React.useState<boolean>(false);
  const [processError, setProcessError] = React.useState<string | null>(null);
  const [valid, setValid] = React.useState<boolean>(false);
  const [formData, setFormData] = React.useState<FormData>({
    name: '',
    specUrl: '',
    file: '',
  });
  const [org, setOrg] = orgState.use();
  const [spec, setSpec] = specState.use();

  const onChange = (name: string, value: string) => {
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };
  const onSubmit = async () => {
    setProcessing(true);
    try {
      setOrg((prev) => ({ ...prev, ...formData }));
      const config = await parseYaml(formData.specUrl, formData.name);
      setSpec(config);
      setProcessing(false);
      setValid(true);
    } catch (err) {
      setValid(true);
      setProcessing(false);
      setProcessError(err.message);
    }
  };

  useInput((input, key) => {
    if (valid && input === 'n' && key.ctrl) {
      history.push('/editor');
    }
  });

  return (
    <Box flexDirection="column">
      <PromptForm complete={valid} onSubmit={onSubmit}>
        <TextField
          required
          label="Organization Name"
          name="name"
          onChange={onChange}
        />
        <TextField
          required
          label="Swagger Docs URL"
          type="url"
          placeholder="URL should end with a .json"
          name="specUrl"
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
            Group config saved! Press ctrl + n to configure plugins.
          </Text>
        </Box>
      )}
    </Box>
  );
};

export default ConfigOrg;
