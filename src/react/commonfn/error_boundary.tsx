// ErrorBoundary only need to define once.

import React from "react";

export class ErrorBoundary extends React.Component<
  { children: JSX.Element; fallBack?: JSX.Element }, // props
  { hasError: boolean } // state
> {
  constructor(props: { children: JSX.Element; fallBack?: JSX.Element }) {
    super(props);
    this.state = { hasError: false };
  }

  // handle state when receiving error from component
  static getDerivedStateFromError(err: Error): { hasError: boolean } {
    return { hasError: true };
  }

  // send errors to logging server
  componentDidCatch(err: Error, errInfo: React.ErrorInfo): void {
    console.log(err);
    console.log(errInfo);
  }

  render(): JSX.Element {
    // render fallback component
    if (this.state.hasError) {
      if (this.props.fallBack) {
        return this.props.fallBack;
      }

      return <div>something went wrong</div>;
    }

    // render children
    return this.props.children;
  }
}
