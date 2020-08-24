import React, { useRef, useState, useEffect } from 'react';
import TextInput from 'ink-text-input';
import { Box, Text, useFocus } from 'ink';
import validUrl from 'valid-url';

interface ArrayFieldProps {
  error: any | undefined;
  name: string;
  onChange: (key: string, value: string) => void;
  required?: boolean;
  type?: 'text' | 'url';
  value?: string | null;
}

const ArrayField = ({
  error,
  onChange,
  name,
  required = false,
  type,
  value,
}: ArrayFieldProps) => {
  const { isFocused } = useFocus();
  const hasError = Boolean(error);
  const focusedColor = isFocused ? 'green' : '';
  const labelColor = hasError ? 'red' : focusedColor;
  const valueString = value.join(', ');
  const changeHandler = (value: string) => {
    onChange(name, value.split(', '));
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
          <TextInput value={valueString} onChange={changeHandler} />
        )}
        {!isFocused && <Text>{valueString}</Text>}
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
