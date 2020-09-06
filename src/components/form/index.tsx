import React, {
  Children,
  cloneElement,
  useCallback,
  useEffect,
  useState,
} from 'react';
import { Box, Text, useInput } from 'ink';
import validate from 'validate.js';

import ArrayField from './array-field';
import TextField from './text-field';
import NumberField from './number-field';
import Button from '../button';

type ChangeHandler = (key: string, value: any) => void;
const getElement = ({
  key,
  field,
  value,
  onChange,
}: {
  key: string;
  field: any;
  value: any;
  onChange: ChangeHandler;
}) => {
  switch (field.type) {
    case 'array':
      return (
        <ArrayField
          key={key}
          name={key}
          required={!!field.presence}
          onChange={onChange}
          value={value}
        />
      );
    case 'string':
      return (
        <TextField
          key={key}
          name={key}
          required={!!field.presence}
          onChange={onChange}
          value={value}
        />
      );
    case 'number':
      return (
        <NumberField
          key={key}
          name={key}
          onChange={onChange}
          required={!!field.presence}
          value={value}
        />
      );
    default:
      return <Box key={key} />;
  }
};

interface FormProps {
  constraints: any;
  data: any;
  onSubmit?: (data: any) => void;
}

const Form: React.FC<FormProps> = ({
  constraints,
  data,
  onSubmit = () => false,
}) => {
  const [errors, setErrors] = useState<string[] | null>(null);
  const [formData, setFormData] = useState<any>(data);
  const elements = [];
  const onChange: ChangeHandler = (key, value) => {
    setFormData((state: any) => ({ ...state, [key]: value }));
  };
  const onSubmitClick = useCallback(() => {
    const errors = validate(formData, constraints, { format: 'flat' });

    if (errors) {
      setErrors(errors);
    } else {
      1;
      onSubmit(formData);
    }
  }, [formData]);

  for (const key in constraints) {
    const field = constraints[key];
    const value = formData[key] || data[key];
    elements.push(getElement({ key, field, value, onChange }));
  }

  useEffect(() => {
    if (data !== formData) {
      setFormData(data);
    }
  }, [data]);

  useInput((input, key) => {
    if (input === 's' && key.ctrl) {
      onSubmitClick();
    }
  });

  return (
    <Box flexDirection="column" margin={1}>
      {errors && (
        <Box flexDirection="column" borderColor="redBright" borderStyle="round">
          {errors.map((err) => (
            <Box key={err}>
              <Text color="red">- {err}</Text>
            </Box>
          ))}
        </Box>
      )}
      <Box flexDirection="column">{elements}</Box>
    </Box>
  );
};

export default Form;
