exports.wsconnect = (onopen) => ({
  type: 'CONNECT',
  onopen
});

exports.startGame = (id) => ({
  type: 'START_GAME',
  id
});

exports.broadcastMessage = (message) => ({
  type: 'BROADCAST_MESSAGE',
  message
});

exports.dragColor = (color) => ({
  type: 'DRAG_COLOR',
  color
});

exports.setSequence = (color, position) => ({
  type: 'SET_SEQUENCE',
  color,
  position
});

exports.setGameId = (id) => ({
  type: 'SET_GAME_ID',
  id
})

exports.setError = (err) => ({
  type: 'SET_ERROR',
  error: err
});

exports.submit = (sequence) => ({
  type: 'SUBMIT',
  sequence
});
