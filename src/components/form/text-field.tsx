import React, { useRef, useState, useEffect, useCallback } from 'react';
import { Box, Text, useInput } from 'ink';
import TextInput from 'ink-text-input';
import validUrl from 'valid-url';

interface TextFieldProps {
  enabled?: boolean;
  editing?: boolean;
  encrypted?: boolean;
  error?: any | undefined;
  focused?: boolean;
  name: string;
  onChange: (key: string, value: string) => void;
  onEncrypt: (key: string) => void;
  required?: boolean;
  type?: 'text' | 'url';
  value?: string | null;
}

const TextField: React.FC<TextFieldProps> = ({
  enabled,
  editing,
  encrypted = false,
  error,
  focused,
  onChange,
  onEncrypt,
  name,
  required = false,
  type,
  value,
}) => {
  const focusedColor = focused ? 'yellow' : 'cyan';
  const labelColor = error ? 'red' : focusedColor;
  const changeHandler = useCallback(
    (value: string) => {
      if (value.includes('~')) {
        onEncrypt(name);
      } else {
        onChange(name, value);
      }
    },
    [focused]
  );

  return (
    <Box>
      <Box marginRight={1}>
        <Text color={labelColor}>{name}</Text>
      </Box>
      <Box flexGrow={1}>
        <TextInput
          focus={focused && enabled}
          value={value ? value : ''}
          onChange={changeHandler}
        />
      </Box>
    </Box>
  );
};

export default TextField;
