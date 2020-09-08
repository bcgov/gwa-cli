import React, { useContext, useEffect, useState } from 'react';
import { Box, Newline, Text } from 'ink';

import AppContext from '../../services/context';
import { buildSpec } from '../../services/kong';

const Review: React.FC = () => {
  const [error, setError] = useState<boolean>(false);
  const { dir, file } = useContext(AppContext);
  useEffect(() => {
    try {
      buildSpec(dir, file);
    } catch (err) {
      setError(true);
    }
  }, [dir, file]);

  return (
    <Box
      marginY={3}
      flexDirection="column"
      justifyContent="center"
      alignItems="center"
      width="100%"
    >
      {error && (
        <Box flexDirection="column" alignItems="center">
          <Text bold color="red">
            Export Failed
          </Text>
          <Text>[ ctrl + j ] to go back</Text>
        </Box>
      )}
      {!error && (
        <>
          <Text bold>All Done!</Text>
          <Newline />
          <Text>
            Your Kong config file has been generated to{' '}
            <Text inverse color="green">{`${file || 'spec.yaml'}`}</Text>.
          </Text>
          <Newline />
          <Text>Commit and push this branch and make a PR to complete.</Text>
          <Newline />
          <Text>Press [ ctrl + c ] to exit</Text>
        </>
      )}
    </Box>
  );
};

export default Review;
