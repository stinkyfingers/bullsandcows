// import { api } from '../config';

const api = 'wss://eqsjdaluec.execute-api.us-west-1.amazonaws.com/dev'// TODO

const game = (state = {}, action) => {
	switch (action.type) {
    case 'CONNECT':
      const ws = getWS(action);
      return Object.assign({}, state, {ws});

    case 'DISCONNECT':
  			state.ws.close();
  			return state;

    case 'START_GAME':
      if (action.id) {
        // resume
        state.ws.send(JSON.stringify({option:"game", gameId: action.id}));
      } else {
        // start
        state.ws.send(JSON.stringify({option:"game"}));
      }
      return state;

    case 'BROADCAST_MESSAGE':
      return Object.assign({}, state, {message: action.message, submitting: false});

		case 'DRAG_COLOR':
			return Object.assign({}, state, {color: action.color});

		case 'SET_SEQUENCE':
			let sequence = state.sequence ? state.sequence : {};
			sequence[action.position] = action.color;
		 	return Object.assign({}, state, {sequence});

		case 'SET_ERROR':
			return Object.assign({}, state, {error: action.error});

		case 'SET_GAME_ID':
			return Object.assign({}, state, {gameId: action.id });

		case 'SUBMIT':
			let seq = [];
			for (let i = 0; i < 4; i++) {
				seq[i] = {"color": action.sequence[i]}
			}
			state.ws.send(JSON.stringify({option: 'round', gameId: state.gameId, "sequence": seq}))
			return Object.assign({}, state, {submitting: true});
		case 'SHOW_RULES':
			let showRules = 'show';
			if (action.showRules === 'show') showRules = 'hide';
      return Object.assign({}, state, {showRules});
    default:
    return state;
  }
}

const getWS = (action) => {
  const ws = new window.WebSocket(api);
  ws.addEventListener('open', action.onopen);
  return ws;
}

export default game;
