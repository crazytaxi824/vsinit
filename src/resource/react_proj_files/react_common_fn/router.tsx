import React from 'react';
import { BrowserRouter, Route, Router, Link } from 'react-router-dom';
import { createBrowserHistory } from 'history';

const history = createBrowserHistory();

export function RouterTest(): JSX.Element {
  return (
    <Router history={history}>
      <div>
        <ul>
          <li>
            <Link to="/">/Root</Link>
          </li>
          <li>
            <Link to="/a">/A</Link>
          </li>
          <li>
            <Link to="/a/b">/A/B</Link>
          </li>
          <li>
            <Link to="/a/b/c">/A/B/C</Link>
          </li>
          <li>
            <Link to="/x">/X</Link>
          </li>
        </ul>
        <hr />
        <Route path="/" component={() => <h2>Root</h2>} />
        <Route exact path="/a" component={() => <h2>A</h2>} />
        <Route path="/a/b" component={() => <h2>AB</h2>} />
      </div>
    </Router>
  );
}
