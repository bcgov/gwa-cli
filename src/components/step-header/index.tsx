import * as React from 'react';
import { Box, Text } from 'ink';

interface StepHeaderProps {
  title: string;
}

const StepHeader: React.FC<StepHeaderProps> = ({ title }) => (
  <Box marginBottom={1}>
    <Text inverse>{` ${title} `}</Text>
  </Box>
);

export default StepHeader;
