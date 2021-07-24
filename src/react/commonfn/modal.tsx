import React from 'react';
import ReactDOM from 'react-dom';

const Modal = React.forwardRef(
  (
    props: unknown,
    mRef: React.Ref<{ sendMsg: (m: string) => void }>
  ): JSX.Element | null => {
    const [msg, setMsg] = React.useState('');
    React.useImperativeHandle(mRef, () => ({ sendMsg: setMsg }));

    if (!msg) return null; // don't display while null

    return ReactDOM.createPortal(
      <div className="modal">
        <span>{msg}</span>
        <button
          type="button"
          onClick={() => {
            setMsg('');
          }}
        >
          Close
        </button>
      </div>,
      document.body // append to html <body>
    );
  }
);
