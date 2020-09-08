import React, { useRef, useState, useEffect } from 'react';
import TextInput from 'ink-text-input';
import { Box, Text } from 'ink';
import validUrl from 'valid-url';

interface ArrayFieldProps {
  error?: any | undefined;
  focused: boolean;
  name: string;
  onChange: (key: string, value: string[]) => void;
  required?: boolean;
  type?: 'text' | 'url';
  value?: string[] | null;
}

const ArrayField: React.FC<ArrayFieldProps> = ({
  error,
  focused,
  onChange,
  name,
  required = false,
  type,
  value,
}) => {
  const hasError = Boolean(error);
  const focusedColor = focused ? 'yellow' : '';
  const labelColor = hasError ? 'red' : focusedColor;
  const valueString = value ? value.join(', ') : '';
  const changeHandler = (value: string) => {
    onChange(name, value.split(', '));
  };

  return (
    <Box>
      <Box marginX={1}>
        <Text bold color={labelColor}>
          {name}
          {required && '*'}:
        </Text>
      </Box>
      <Box flexGrow={1} width="50%">
        {focused && <TextInput value={valueString} onChange={changeHandler} />}
        {!focused && <Text>{valueString}</Text>}
      </Box>
      {hasError && (
        <Box>
          <Text color="red">{error}</Text>
        </Box>
      )}
    </Box>
  );
};

export default ArrayField;
