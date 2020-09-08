import React, { useContext, useEffect } from 'react';
import { Box, Newline, Text } from 'ink';

import AppContext from '../../services/context';
import { buildSpec } from '../../services/kong';

const Review: React.FC = () => {
  const { dir, file } = useContext(AppContext);
  useEffect(() => {
    buildSpec(dir, file);
  }, [dir, file]);

  return (
    <Box
      marginY={3}
      flexDirection="column"
      justifyContent="center"
      alignItems="center"
      width="100%"
    >
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
    </Box>
  );
};

export default Review;
