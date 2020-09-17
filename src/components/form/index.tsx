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

type ChangeHandler = (key: string, value: any) => void;
const getElement = ({
  key,
  field,
  value,
  onChange,
  onSubmit,
}: {
  key: string;
  field: any;
  value: any;
  onChange: ChangeHandler;
  onSubmit: () => void;
}) => {
  switch (field.type) {
    case 'array':
      return (
        <ArrayField
          name={key}
          onChange={onChange}
          onSubmit={onSubmit}
          value={value}
        />
      );
    case 'string':
      return (
        <TextField
          name={key}
          onChange={onChange}
          onSubmit={onSubmit}
          value={value}
        />
      );
    case 'number':
      return (
        <NumberField
          name={key}
          onChange={onChange}
          onSubmit={onSubmit}
          value={value}
        />
      );
    case 'boolean':
      return <Checkbox name={key} onChange={onChange} checked={value} />;
    default:
      return <Box />;
  }
};

interface FormProps {
  constraints: any;
  data: any;
  enabled: boolean;
  encryptedFields: string[];
  onEncrypt: (key: string, encrypted: boolean) => void;
  onSubmit?: (data: any) => void;
}

const Form: React.FC<FormProps> = ({
  constraints,
  data,
  enabled,
  encryptedFields,
  onEncrypt,
  onSubmit = () => false,
}) => {
  const [tabIndex, setTabIndex] = useState<number>(0);
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
  const focusNext = () => {
    setTabIndex(Math.min(tabIndex + 1, elements.length));
  };

  for (const key in constraints) {
    const field = constraints[key];
    const value = formData[key] || data[key];
    const el = getElement({ key, field, value, onChange, onSubmit: focusNext });

    if (field) {
      const fieldIndex = elements.length + 1;
      elements.push(
        <FieldSet
          key={key}
          focused={tabIndex === elements.length}
          enabled={enabled}
          error={errors ? errors.hasOwnProperty(key) : false}
          encrypted={encryptedFields.includes(key)}
          index={fieldIndex}
          editing={enabled}
          name={key}
          onEncrypt={onEncrypt}
          required={!!field.presence}
        >
          {el}
        </FieldSet>
      );
    }
  }

  useEffect(() => {
    if (data !== formData) {
      setFormData(data);
      setTabIndex(0);
      setErrors(null);
    }
  }, [data]);

  useEffect(() => {
    return () => setTabIndex(0);
  }, []);

  useInput((input, key) => {
    if (enabled) {
      if (key.escape) {
        onSubmitClick();
      }
    }

    if (key.tab && elements.length > 1) {
      if (key.shift) {
        setTabIndex(Math.max(tabIndex - 1, 0));
      } else {
        setTabIndex(Math.min(tabIndex + 1, elements.length));
      }
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
