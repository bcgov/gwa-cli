import React, { Children, useState, cloneElement, useCallback } from 'react';
import { Box, Text, useInput } from 'ink';
import { FormValue } from './types';

interface PromptFormProps {
  children: React.ReactElement[];
  onSubmit: () => void;
}

const PromptForm: React.FC<PromptFormProps> = ({ children, onSubmit }) => {
  const [formIndex, setFormIndex] = useState<number>(0);
  const [items, setItems] = useState<FormValue[]>([]);
  const onEnter = useCallback(
    (value: Omit<FormValue, 'id'>) => {
      setItems((s) => [...s, { id: Math.random(), ...value }]);
      setFormIndex((s) => s + 1);
    },
    [setItems, setFormIndex]
  );
  const totalFields = Children.count(children);
  const isFinished = totalFields < formIndex + 1;
  let formElement: React.ReactElement | null = null;

  useInput((input, key) => {
    if (isFinished && input === 'Y') {
      onSubmit();
    }
  });

  Children.forEach(children, (child: React.ReactElement, index: number) => {
    if (index === formIndex) {
      formElement = cloneElement(child, {
        onEnter,
      });
    }
  });

  return (
    <Box flexDirection="column">
      <Box flexDirection="column" marginBottom={1}>
        {items.map((item) => (
          <Box key={item.id}>
            <Box marginRight={1}>
              <Text color="green">
                ✔ <Text bold>{item.label}:</Text>
              </Text>
            </Box>
            <Text underline>{item.value}</Text>
          </Box>
        ))}
        {isFinished && (
          <Box marginY={1}>
            <Text bold>All Finished! Type Y to confirm</Text>
          </Box>
        )}
      </Box>
      {formElement}
    </Box>
  );
};

export default PromptForm;
