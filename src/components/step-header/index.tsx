import * as React from 'react';
import { Box, Text } from 'ink';

interface StepHeaderProps {
  title: string;
}

const StepHeader: React.FC<StepHeaderProps> = ({ step, title }) => (
  <Box marginBottom={1}>
    <Text>{title}</Text>
  </Box>
);

export default StepHeader;
