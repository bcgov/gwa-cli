import * as React from 'react';
import { Box, Text } from 'ink';
import fetch from 'node-fetch';
import fs from 'fs';
const { resolve } = require('path');
import Spinner from 'ink-spinner';
import { UncontrolledTextInput } from 'ink-text-input';

import StepHeader from '../../components/step-header';
import { parseYaml } from '../../services/kong';

interface ConfigSrcProps {
  onComplete: () => void;
  step: number;
}

const ConfigSrc = ({ onComplete, step }: ConfigSrcProps) => {
  const [loading, setLoading] = React.useState<boolean>(false);
  const [error, setError] = React.useState<string | null>(null);
  const onSubmit = async (value: string) => {
    setError(null);
    try {
      const res = await fetch(value);
      const json = await res.json();
      const config = await parseYaml(json);
      fs.writeFileSync(resolve('./test.yml'), config);
      onComplete();
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <Box flexDirection="column">
      <StepHeader step={step} title="Configure Your Organization" />
      <Box>
        <Box marginRight={1}>
          {loading && <Spinner type="dots" />}
          <Text color="white">Spec URL (JSON):</Text>
        </Box>
        <Box>
          <UncontrolledTextInput
            placeholder="URL should end with .json"
            onSubmit={onSubmit}
          />
        </Box>
      </Box>
      {error && (
        <Box>
          <Text color="red">{`Error: ${error}`}</Text>
        </Box>
      )}
    </Box>
  );
};

export default ConfigSrc;
