import React, { Suspense } from 'react';
import { ErrorBoundary } from 'react-error-boundary';

import Failed from '../failed';
import Loading from '../loading';

interface AsyncActionProps {
  children: React.ReactNode;
  loadingText: string;
  verbose?: boolean;
}

const AsyncAction: React.FC<AsyncActionProps> = ({
  children,
  loadingText,
  verbose,
}) => {
  return (
    <ErrorBoundary
      fallbackRender={({ error }) => <Failed error={error} verbose={verbose} />}
    >
      <Suspense fallback={<Loading text={loadingText} />}>{children}</Suspense>
    </ErrorBoundary>
  );
};

export default AsyncAction;
