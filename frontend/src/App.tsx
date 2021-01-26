import * as React from 'react';
import '../node_modules/@blueprintjs/core/lib/css/blueprint.css';
import '../node_modules/@blueprintjs/icons/lib/css/blueprint-icons.css';
import './App.css';
import Sidebar from "./ui/Sidebar";

interface IAppProps {
    toastHandler: JSX.Element,
}

class App extends React.Component<IAppProps> {
  render() {
    return <div className="app">
        {this.props.toastHandler}
        <div className="app__sidebar">
            <Sidebar />
        </div>
        <div className="app__main">
            <div className="app__pages">
              {this.props.children}
            </div>
        </div>
      </div>
    ;
  }
}

export default App;
