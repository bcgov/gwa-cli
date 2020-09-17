import React, { useCallback, useRef, useState, useEffect } from 'react';
import TextInput from 'ink-text-input';
import { Box, Text, useInput } from 'ink';
import compact from 'lodash/compact';
import validUrl from 'valid-url';

interface ArrayFieldProps {
  enabled?: boolean;
  error?: any | undefined;
  focused?: boolean;
  name: string;
  onChange: (key: string, value: string[]) => void;
  onSubmit: () => void;
  required?: boolean;
  type?: 'text' | 'url';
  value?: string[] | null;
}

const ArrayField: React.FC<ArrayFieldProps> = ({
  enabled = false,
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
  const valueString = value ? value.join(', ') : '';
  const changeHandler = (value: string) => {
    const arrValue = value.split(/,\s+/g);
    onChange(name, compact(arrValue));
  };

  return (
    <Box>
      <Box marginX={1}>
        <Text color={labelColor}>{name}</Text>
      </Box>
      <Box flexGrow={1} width="50%">
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

export default ArrayField;
