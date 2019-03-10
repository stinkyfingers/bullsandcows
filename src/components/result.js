import React, { Component } from 'react';
import { connect } from 'react-redux';
import '../css/round.css';

class Result extends Component {

  render() {
    return (
      <div className='result'>
        <label className='bulls'>Bulls: <span>{this.props.bulls ? this.props.bulls : 0}</span></label>
        <label className='cows'>Cows: <span>{this.props.cows ? this.props.cows : 0}</span></label>
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

export default connect(mapStateToProps, mapDispatchToProps)(Result);
