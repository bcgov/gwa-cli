import React, { useCallback, useState } from 'react';
import TextInput from 'ink-text-input';
import { Box, Text, useFocus, Transform } from 'ink';
import validUrl from 'valid-url';

import { FormValue } from './types';

interface TextFieldProps {
  label: string;
  name: string;
  onChange: (key: string, value: string) => void;
  onEnter?: (value: Omit<FormValue, 'id'>) => void;
  placeholder?: string;
  required: boolean;
  type?: 'text' | 'url';
}

const TextField: React.FC<TextFieldProps> = ({
  label,
  name,
  onEnter = () => false,
  onChange,
  placeholder,
  required,
  type,
}) => {
  const { isFocused } = useFocus();
  const [value, setValue] = useState<string>('');
  const [error, setError] = useState<string>('');
  const hasError = Boolean(error);
  const labelColor = hasError ? 'red' : 'yellow';
  const handleChange = (v: string) => {
    setError('');
    setValue(v);
    onChange(name, v);
  };
  const onSubmit = useCallback(
    (value: string) => {
      if (!value && required) {
        setError('This field is required');
        return;
      }

      if (value) {
        if (type === 'url' && !validUrl.isUri(value)) {
          setError('Not a valid URL');
          return;
        }

        if (value && !hasError) {
          setValue('');
          onEnter({ label, value });
        }
      }
    },
    [value, onEnter, setError, setValue, type]
  );

  return (
    <Box flexDirection="column">
      <Box>
        <Box marginRight={1}>
          <Text bold color={labelColor}>
            {label}:
          </Text>
        </Box>
        <Box flexGrow={1}>
          {isFocused && (
            <TextInput
              focus
              placeholder={`${required && '[Required] '}${placeholder || ''}`}
              value={value}
              onChange={handleChange}
              onSubmit={onSubmit}
            />
          )}
          {!isFocused && <Text>{value}</Text>}
        </Box>
      </Box>
      {hasError && (
        <Box>
          <Text bold color="red">
            <Transform
              transform={(output: string) => output.replace(/./gi, '^')}
            >
              <Text>{label}:</Text>
            </Transform>
            {` ${error}`}
          </Text>
        </Box>
      )}
    </Box>
  );
};

export default TextField;
