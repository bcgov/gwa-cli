import React from 'react';

import useAsync from '../../../hooks/use-async';
import Success from '../../../components/success';
import { makeConfigFile } from '../create-actions';

interface WriteConfigActionProps {
  data: any;
}

const WriteConfigAction: React.FC<WriteConfigActionProps> = ({ data }) => {
  const { url, ...options } = data;
  const result = useAsync(makeConfigFile, url, options);

  return <Success>{`Config ${result} successfully generated`}</Success>;
};

export default WriteConfigAction;
