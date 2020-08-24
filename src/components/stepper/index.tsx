import React, { Children } from 'react';
import { Box } from 'ink';

interface StepperProps {
  children: React.ReactElement[];
  step: number;
}

const Stepper = ({ children, step }: StepperProps) => {
  let element = null;

  Children.forEach(children, (child: React.ReactElement, index: number) => {
    if (index === step) {
      element = child;
    }
  });

  return <Box>{element}</Box>;
};

export default Stepper;
