import React from 'react';

import { api } from '../../services/api';
import Success from '../../components/success';
import useAsync from '../../hooks/use-async';

interface AclRequestViewProps {
  data: any[];
}

const AclRequestView: React.FC<AclRequestViewProps> = ({ data }) => {
  const result = useAsync(api, 'namespaces/:namespace/membership', {
    method: 'PUT',
    body: JSON.stringify(data),
  });

  return <Success>{result}</Success>;
};

export default AclRequestView;
