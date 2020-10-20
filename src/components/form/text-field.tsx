import React, { useRef, useState, useEffect, useCallback } from 'react';
import { Box, Text, useInput } from 'ink';
import TextInput from 'ink-text-input';
import validUrl from 'valid-url';

interface TextFieldProps {
  enabled?: boolean;
  error?: any | undefined;
  focused?: boolean;
  name: string;
  onChange: (key: string, value: string) => void;
  onSubmit: () => void;
  type?: 'text' | 'url';
  value: string;
}

const TextField: React.FC<TextFieldProps> = ({
  enabled,
  error,
  focused,
  onChange,
  onSubmit,
  name,
  type,
  value,
}) => {
  const focusedColor = focused ? 'yellow' : 'cyan';
  const labelColor = error ? 'red' : focusedColor;
  const changeHandler = (v: string) => {
    onChange(name, v);
  };
  const valueString = value ? value : '';

  return (
    <Box>
      <Box marginRight={1}>
        <Text color={labelColor}>{name}</Text>
      </Box>
      <Box flexGrow={1}>
        <TextInput
          focus={focused && enabled}
          value={valueString}
          onChange={changeHandler}
          onSubmit={onSubmit}
        />
      </Box>
    </Box>
  );
};

export default TextField;
