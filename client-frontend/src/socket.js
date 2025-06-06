let socket = null
let callbacks = {}

export function connectSocket(clientId, roomId) {
  socket = new WebSocket(`${import.meta.env.VITE_WS_URL}?client_id=${clientId}`)

  socket.onopen = () => {
    console.log('WebSocket connected')

    if (roomId && clientId) {
      sendSocketMessage('init', {
        room_id: roomId,
        client_id: clientId,
      })
    }
  }

  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data)
    const handler = callbacks[msg.event]
    if (handler) {
      handler(msg.data)
    }
  }

  socket.onclose = () => {
    console.warn('Socket closed. Reconnecting in 3s...')
    setTimeout(() => connectSocket(clientId), 3000)
  }
}

export function sendSocketMessage(event, data) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({ event, data }))
  } else {
    console.error('Socket not connected')
  }
}

export function subscribe(event, handler) {
  callbacks[event] = handler
}

export function unsubscribe(event) {
  delete callbacks[event]
}
