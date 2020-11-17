import React from 'react';
import { Box, Text } from 'ink';

import useAsync from '../../../hooks/use-async';
import Success from '../../../components/success';

interface WriteConfigActionProps {
  data: any;
  submitHandler: any;
}

const WriteConfigAction: React.FC<WriteConfigActionProps> = ({
  data,
  submitHandler,
}) => {
  const { url, ...options } = data;
  const result = useAsync(submitHandler, url, options);

  return <Success>{`Config ${result} successfully generated`}</Success>;
};

export default WriteConfigAction;
