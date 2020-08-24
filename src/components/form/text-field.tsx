import React, { useRef, useState, useEffect } from 'react';
import TextInput from 'ink-text-input';
import { Box, Text, useFocus } from 'ink';
import validUrl from 'valid-url';

interface TextFieldProps {
  error: any | undefined;
  name: string;
  onChange: (key: string, value: string) => void;
  required?: boolean;
  type?: 'text' | 'url';
  value?: string | null;
}

const TextField = ({
  error,
  onChange,
  name,
  required = false,
  type,
  value,
}: TextFieldProps) => {
  const { isFocused } = useFocus();
  const hasError = Boolean(error);
  const focusedColor = isFocused ? 'green' : '';
  const labelColor = hasError ? 'red' : focusedColor;
  const changeHandler = (value: string) => {
    onChange(name, value);
  };

  return (
    <Box>
      <Box marginRight={1}>
        <Text bold color={labelColor}>
          {name}
          {required && '*'}:
        </Text>
      </Box>
      <Box flexGrow={1} width="50%">
        {isFocused && (
          <TextInput
            value={value === null ? 'null' : value}
            onChange={changeHandler}
          />
        )}
        {!isFocused && <Text>{value === null ? 'null' : value}</Text>}
      </Box>
      {hasError && (
        <Box>
          <Text color="red">{error}</Text>
        </Box>
      )}
    </Box>
  );
};

export default TextField;
