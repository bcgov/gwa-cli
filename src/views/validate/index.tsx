import React, { useEffect, useState } from 'react';
import { Box, Text, render } from 'ink';
import Spinner from 'ink-spinner';

interface ValidateProps {
  file: string;
}

const Validate: React.FC<ValidateProps> = ({ file }) => {
  const [validating, setValidating] = useState<boolean>(true);
  useEffect(() => {
    setTimeout(() => {
      setValidating(false);
    }, 5000);
  }, []);

  return (
    <Box>
      {validating && (
        <Box>
          <Box marginRight={2}>
            <Spinner />
          </Box>
          <Text>{`Validating ${file}`}</Text>
        </Box>
      )}
      {!validating && (
        <Box>
          <Text bold color="green">
            {'✓ '}
          </Text>
          <Text>File is valid</Text>
        </Box>
      )}
    </Box>
  );
};

export default (file: string) => render(<Validate file={file} />);
