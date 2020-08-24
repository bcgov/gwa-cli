import React from 'react';
import { Box, Text, useFocus, useInput } from 'ink';

interface ButtonProps {
  children: string | React.ReactElement | React.ReactElement[];
  color?: string;
  onClick: () => void;
}

const Button = ({ children, color = 'white', onClick }: ButtonProps) => {
  const { isFocused } = useFocus();
  const bgColor = isFocused ? 'green' : color;

  useInput((input, key) => {
    if (isFocused) {
      if (key.return) {
        onClick();
      }
    }
  });

  return (
    <Box>
      <Text inverse color={bgColor}>
        {children}
      </Text>
    </Box>
  );
};

export default Button;
