import React, { useContext, useEffect } from 'react';
import { Box, Newline, Text } from 'ink';

import AppContext from '../../services/context';
import { buildSpec } from '../../services/kong';

const Review: React.FC = () => {
  const { dir } = useContext(AppContext);
  useEffect(() => {
    buildSpec(dir);
  }, []);

  return (
    <Box flexDirection="column" justifyContent="center">
      <Text bold>All Done!</Text>
      <Newline />
      <Text>Your Kong config file has been generated to spec.yaml.</Text>
      <Newline />
      <Text>Commit and push this branch and make a PR to complete.</Text>
      <Newline />
      <Text>Press [ CTRL + C ] to exit</Text>
    </Box>
  );
};

export default Review;
