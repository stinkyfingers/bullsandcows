import React, { Component } from 'react';
import { connect } from 'react-redux';
import '../css/board.css';

class Round extends Component {

  renderSequence(seq) {
    let arr = [];
    for (let i = 0; i < seq.length; i++) {
      arr.push(
        <div key={'seq-' + i} className='color' id={seq[i].color}>
        </div>
      );
    }
    return (<div className='round'>{arr}</div>);
  }

  render() {
    return (
      <div>
        {this.renderSequence(this.props.sequence)}
      </div>);
  }

}

const mapStateToProps = state => {
	return {
    message: state.game.message
	};
}

const mapDispatchToProps = dispatch => ({
});

export default connect(mapStateToProps, mapDispatchToProps)(Round);
