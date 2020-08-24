import * as React from 'react';
import { Box, Newline, Text } from 'ink';

const Review = () => (
  <Box flexDirection="column">
    <Text bold>All Done!</Text>
    <Newline />
    <Text>Your Kong config file has been generated to medications.yml.</Text>
    <Newline />
    <Text>Commit and push this branch and make a PR to complete.</Text>
    <Newline />
    <Text>Press [ q ] to exit</Text>
  </Box>
);

export default Review;
