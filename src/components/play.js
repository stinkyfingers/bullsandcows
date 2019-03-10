import React, { Component } from 'react';
import { connect } from 'react-redux';
import { setSequence, setError, submit, dragColor } from '../actions/game';

import '../css/play.css';

class Play extends Component {
  constructor() {
    super();
    this.ondragover = this.ondragover.bind(this);
    this.ondrop = this.ondrop.bind(this);
    this.onclick = this.onclick.bind(this);
    this.onmousedown = this.onmousedown.bind(this);
  }

  ondragover(e) {
    e.preventDefault();
  }

  ondrop(e) {
    this.props.seterror(null);
    let key = '';
    switch(e.target.id) {
      case '0':
        key = 'zero'; break;
      case '1':
        key = 'one'; break;
      case '2':
        key = 'two'; break;
      case '3':
        key = 'three'; break;
      default:
    }
    this.setState({ [key]: this.props.color });
    this.props.setsequence(this.props.color, parseInt(e.target.id, 10));
    this.onmousedown = this.onmousedown.bind(this);
  }

  onclick() {
    if (!this.props.sequence) {
      this.props.seterror('You must populate the entire sequence');
      return;
    }

    // validate
    for (let i = 0; i < 4; i++) {
      if (!this.props.sequence[i]) {
        this.props.seterror('You must populate the entire sequence');
        return;
      }
    }
    this.props.submit(this.props.sequence);
    this.setState({ zero: null, one: null, two: null, three: null});
  }

  onmousedown(e) {
    this.props.dragcolor(e.target.id);
  }

  ontouchstart(e) {
    this.props.dragcolor(e.target.id);
  }

  ontouchset(e) {
    this.ondrop(e);
  }

  renderPicker() {
    return (
      <div className='picker'>
        <div className='color' id='blue' style={{border: this.props.color === 'blue' ? '3px solid black' : null}} draggable onMouseDown={this.onmousedown} onTouchStart={this.ontouchstart.bind(this)}></div>
        <div className='color' id='red' style={{border: this.props.color === 'red' ? '3px solid black' : null}} draggable onMouseDown={this.onmousedown} onTouchStart={this.ontouchstart.bind(this)}></div>
        <div className='color' id='green' style={{border: this.props.color === 'green' ? '3px solid black' : null}} draggable onMouseDown={this.onmousedown} onTouchStart={this.ontouchstart.bind(this)}></div>
        <div className='color' id='yellow' style={{border: this.props.color === 'yellow' ? '3px solid black' : null}} draggable onMouseDown={this.onmousedown} onTouchStart={this.ontouchstart.bind(this)}></div>
        <div className='color' id='purple' style={{border: this.props.color === 'purple' ? '3px solid black' : null}} draggable onMouseDown={this.onmousedown} onTouchStart={this.ontouchstart.bind(this)}></div>
        <div className='color' id='orange' style={{border: this.props.color === 'orange' ? '3px solid black' : null}} draggable onMouseDown={this.onmousedown} onTouchStart={this.ontouchstart.bind(this)}></div>
      </div>
      );
  }

  render() {
    return (
      <div className='sequence'>
        <div>
          <div className={'peg ' + (this.state ? this.state.zero : '')} onDragOver={this.ondragover} onDrop={this.ondrop} value='zero' id='0' onTouchStart={this.ontouchset.bind(this)}></div>
          <div className={'peg ' + (this.state ? this.state.one : '')} onDragOver={this.ondragover} onDrop={this.ondrop} name='one' id='1' onTouchStart={this.ontouchset.bind(this)}></div>
          <div className={'peg ' + (this.state ? this.state.two : '')} onDragOver={this.ondragover} onDrop={this.ondrop} name='two' id='2' onTouchStart={this.ontouchset.bind(this)}></div>
          <div className={'peg ' + (this.state ? this.state.three : '')} onDragOver={this.ondragover} onDrop={this.ondrop} name='three' id='3' onTouchStart={this.ontouchset.bind(this)}></div>
        </div>
        {this.props.message && this.props.message.concluded ? null : this.renderPicker()}
        {this.props.message && this.props.submitting ? <div className="submitting">Waiting...</div> : null}
        <button onClick={this.onclick}>Submit</button>
      </div>)
  }

}

const mapStateToProps = state => {
	return {
    message: state.game.message,
    color: state.game.color,
    sequence: state.game.sequence,
    error: state.game.error,
    submitting: state.game.submitting
	};
}

const mapDispatchToProps = dispatch => ({
  setsequence: (color, position) => dispatch(setSequence(color, position)),
  seterror: (err) => dispatch(setError(err)),
  submit: (sequence) => dispatch(submit(sequence)),
  dragcolor: (color) => dispatch(dragColor(color))
});

export default connect(mapStateToProps, mapDispatchToProps)(Play);
