import React, { Component } from 'react';
import { connect } from 'react-redux';
import { startGame,  wsconnect, broadcastMessage, setError, setGameId } from '../actions/game';
import '../css/header.css';

class Header extends Component {
  constructor() {
    super();
    this.handleClick = this.handleClick.bind(this);
    this.onopen = this.onopen.bind(this);
    this.onmessage = this.onmessage.bind(this);
  }

  componentWillMount() {
    this.props.connect(this.onopen);
  }

  handleClick() {
    this.props.startGame(null);
  }

  onopen() {
    if (this.props.ws) {
			this.props.ws.addEventListener('message', this.onmessage);
      const id = this.getID();
      if (id !== '') {
        this.props.startGame(id);
      }
		}
  }

  onmessage(msg, err) {
    if (err) { this.props.seterror(err); return }
    const data = JSON.parse(msg.data);
    if (data.error) {
      this.props.seterror(data.error);
      return;
    }
    this.props.seterror(null);
    const id = this.getID();
    if (id !== data.id) {
      window.location.hash = data.id;
    }
    this.props.broadcastMessage(data);
    this.props.setGameId(data.id);
  }

  getID() {
    return window.location.hash.replace('#','');
  }

  render() {
    return (
      <div className="header">
        <button onClick={this.handleClick}>New Game</button>
      </div>
    );
  }

}

const mapStateToProps = state => {
	return {
    gameId: state.game.gameId,
    ws: state.game.ws
	};
}

const mapDispatchToProps = dispatch => ({
  startGame: (id) => dispatch(startGame(id)),
  setGameId: (id) => dispatch(setGameId(id)),
  connect: (onopen) => dispatch(wsconnect(onopen)),
  broadcastMessage: (message) => dispatch(broadcastMessage(message)),
  seterror: (error) => dispatch(setError(error))
});

export default connect(mapStateToProps, mapDispatchToProps)(Header);
