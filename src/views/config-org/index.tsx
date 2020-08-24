import * as React from 'react';
import { Box, Text } from 'ink';

import PromptForm from '../../components/prompt-form';
import TextField from '../../components/prompt-form/text-field';
import StepHeader from '../../components/step-header';
import { orgState } from '../../state/org';

interface ConfigOrgProps {
  onComplete: () => void;
  step: number;
}
const ConfigOrg = ({ onComplete, step }: ConfigOrgProps) => {
  const [formData, setFormData] = React.useState({});
  const [org, setOrg] = orgState.use();

  const onChange = (name: string, value: any) => {
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };
  const onSubmit = () => {
    console.log('Done!');
    setOrg(formData);
    onComplete();
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
    </Box>
  );
};

export default ConfigOrg;
