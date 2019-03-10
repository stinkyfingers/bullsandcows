import React from 'react';
import ReactDOM from 'react-dom';
import Body from '../components/body';
import Play from '../components/play';
import { shallow, configure } from 'enzyme';
import { Provider } from 'react-redux';
import { createStore } from 'redux';
import rootReducer from '../reducers/root';
import Adapter from 'enzyme-adapter-react-16';

const store = createStore(rootReducer);
configure({adapter: new Adapter()});


it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<Provider store={store}><Body /></Provider>, div);
  ReactDOM.unmountComponentAtNode(div);
});

it('renders without crashing', () => {
  shallow(<Play />);
});
