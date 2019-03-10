import React, { Component } from 'react';
import { connect } from 'react-redux';
import Board from './board';

class Body extends Component {

  render() {
    return (<div>
        <Board />
      </div>)
  }

}

const mapStateToProps = state => {
	return {
    gameId: state.game.gameId,
    message: state.game.message
	};
}

const mapDispatchToProps = dispatch => ({

});

export default connect(mapStateToProps, mapDispatchToProps)(Body);
