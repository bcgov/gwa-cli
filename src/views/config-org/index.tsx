import * as React from 'react';
import { Box, Text } from 'ink';
import Spinner from 'ink-spinner';

import PromptForm from '../../components/prompt-form';
import TextField from '../../components/prompt-form/text-field';
import StepHeader from '../../components/step-header';
import { orgState } from '../../state/org';
import { parseYaml } from '../../services/kong';
import { specState } from '../../state/spec';

type FormData = {
  name: string;
  specUrl: string;
  maintainers: string[];
};

interface ConfigOrgProps {
  onComplete: () => void;
  step: number;
}

const ConfigOrg: React.FC<ConfigOrgProps> = ({ onComplete, step }) => {
  const [isProcessing, setProcessing] = React.useState<boolean>(false);
  const [processError, setProcessError] = React.useState<string | null>(null);
  const [formData, setFormData] = React.useState<FormData>({
    name: '',
    specUrl: '',
    maintainers: [],
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
      setOrg(formData);
      const config = await parseYaml(formData.specUrl, formData.name);
      setSpec(config);
      onComplete();
    } catch (err) {
      setProcessing(false);
      setProcessError(err.message);
    }
  };

  return (
    <Box flexDirection="column">
      <StepHeader step={step} title="Configure Your Organization" />
      <PromptForm onSubmit={onSubmit}>
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
          label="Maintainers (comma separated)"
          name="maintainers"
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
    </Box>
  );
};

export default ConfigOrg;
