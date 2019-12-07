import React, { Component } from 'react';
import { connect } from 'react-redux';
import Round from './round';
import Result from './result';
import Play from './play';
import { dragColor } from '../actions/game';
import '../css/board.css'

class Board extends Component {
  constructor () {
    super();
    this.onmousedown = this.onmousedown.bind(this);
  }

  onmousedown(e) {
    this.props.dragcolor(e.target.id)
  }

  renderResultRounds() {
    let arr = [];
    this.props.message.rounds.forEach((e, i) => {
      arr.push(<div className='roundResult' key={i}>
        <label className='roundLabel'>Round {i + 1}</label>
        <Result bulls={e.bulls} cows={e.cows}/>
        <Round sequence={e.sequence}/>
      </div>)
    });
    return (<div className='roundResults'>{arr}</div>);
  }

  renderGameState() {
    if (!this.props.message) {
      return (<Play />);
    }
    if (this.props.message.won) {
      return (<div className='winner'>You're a Winner!</div>);
    }
    if (this.props.message.concluded) {
      return (
        <div className='loser'>
          <div className='loserText'>Game over. You're a loser!</div>
          <div className='color' id={`${this.props.message.code[0].color}`}></div>
          <div className='color' id={`${this.props.message.code[1].color}`}></div>
          <div className='color' id={`${this.props.message.code[2].color}`}></div>
          <div className='color' id={`${this.props.message.code[3].color}`}></div>
        </div>
      );
    }
    return (<Play />);
  }

  render() {
    if (!this.props) return null;
    return (
      <div>
        {this.props.error ? <div className='error'>{this.props.error}</div> : null}
        {this.props && this.props.message && this.props.message.rounds ? this.renderResultRounds() : null}
        {this.props.gameId ? <div className='play'>
          {this.renderGameState()}
        </div> : null}
      </div>)
  }
}

const mapStateToProps = state => {
	return {
    gameId: state.game.gameId,
    message: state.game.message,
    error: state.game.error
	};
}

const mapDispatchToProps = dispatch => ({
  dragcolor: (color) => dispatch(dragColor(color))
});

export default connect(mapStateToProps, mapDispatchToProps)(Board);
