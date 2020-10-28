import React, { useState, useEffect } from 'react';
import { Box, Text, useInput } from 'ink';
import Spinner from 'ink-spinner';
import { useHistory } from 'react-router';

import { useAppState } from '../../state/app';
import PromptForm from '../../components/prompt-form';
// import TextField from '../../components/prompt-form/text-field';
import { useTeamState } from '../../state/team';

type FormData = any;
interface ConfigOrgProps {}

const ConfigOrg: React.FC<ConfigOrgProps> = ({}) => {
  const history = useHistory();
  const team = useTeamState();
  const toggleMode = useAppState((state) => state.toggleMode);
  const [isProcessing, setProcessing] = useState<boolean>(false);
  const [processError, setProcessError] = useState<string | null>(null);
  const [valid, setValid] = useState<boolean>(false);
  const [formData, setFormData] = useState<FormData>({
    name: '',
    specUrl: '',
    file: '',
  });

  const onChange = (name: string, value: string) => {
    setFormData((prev: any) => ({
      ...prev,
      [name]: value,
    }));
  };
  const onSubmit = () => {
    setProcessing(true);
    //try {
    //  setOrg((prev) => ({ ...prev, ...formData }));
    //  const config = await parseYaml(formData.specUrl, formData.name);
    //  setProcessing(false);
    //  setValid(true);
    //} catch (err) {
    //  setValid(true);
    //  setProcessing(false);
    //  setProcessError(err.message);
    //}
  };

  useInput((input, key) => {
    if (valid && key.return) {
      history.push('/editor');
    }
  });

  useEffect(() => {
    toggleMode();
    return () => toggleMode();
  }, []);

  return (
    <Box flexDirection="column" marginTop={2}>
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
            Group config saved! Press [ENTER] to configure plugins.
          </Text>
        </Box>
      )}
    </Box>
  );
};

export default ConfigOrg;
