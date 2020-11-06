import React from 'react';
import { Text } from 'ink';

import { makeEnvFile } from '../../services/app';
import Success from '../../components/success';
import useAsync from '../../hooks/use-async';
import type { InitOptions } from '../../types';

interface WriteEnvActionProps {
  data: InitOptions;
}

const WriteEnvAction: React.FC<WriteEnvActionProps> = ({ data }) => {
  const result = useAsync(makeEnvFile, data);

  return <Success>{result}</Success>;
};

export default WriteEnvAction;
