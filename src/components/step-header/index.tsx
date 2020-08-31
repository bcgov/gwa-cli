import * as React from 'react';
import { Box, Text } from 'ink';

interface StepHeaderProps {
  step: number;
  title: string;
}

const StepHeader: React.FC<StepHeaderProps> = ({ step, title }) => (
  <Box marginBottom={1}>
    <Text>
      <Text bold>{`Step ${(step + 1).toString()}:`}</Text> {title}
    </Text>
  </Box>
);

export default StepHeader;
