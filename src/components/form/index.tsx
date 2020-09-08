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
import Checkbox from './checkbox';
import Errors from './errors';
import FieldSet from './field-set';
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
          name={key}
          required={!!field.presence}
          onChange={onChange}
          value={value}
        />
      );
    case 'string':
      return (
        <TextField
          name={key}
          required={!!field.presence}
          onChange={onChange}
          value={value}
        />
      );
    case 'number':
      return (
        <NumberField
          name={key}
          onChange={onChange}
          required={!!field.presence}
          value={value}
        />
      );
    case 'boolean':
      return (
        <Checkbox
          required={!!field.precense}
          name={key}
          onChange={onChange}
          checked={value}
        />
      );
    default:
      return <Box />;
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
    const errors = validate(formData, constraints);

    if (errors) {
      setErrors(errors);
    } else {
      setErrors(null);
      onSubmit(formData);
    }
  }, [formData]);

  for (const key in constraints) {
    const field = constraints[key];
    const value = formData[key] || data[key];
    const el = getElement({ key, field, value, onChange });

    if (field) {
      elements.push(
        <FieldSet
          error={errors && errors.hasOwnProperty(key)}
          key={key}
          index={elements.length + 1}
        >
          {el}
        </FieldSet>
      );
    }
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
    <Box flexDirection="column">
      <Box flexDirection="column">{elements}</Box>
      {errors && <Errors errors={errors} />}
    </Box>
  );
};

export default Form;
