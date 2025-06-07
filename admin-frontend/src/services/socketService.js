let socket = null;
let reconnectTimeout = null;

const listeners = {
  open: [],
  close: [],
  error: [],
  message: [],
};

export function initSocket(token) {
  const url = `${import.meta.env.VITE_WS_URL}?token=${token}`;
  socket = new WebSocket(url);

  socket.onopen = () => {
    console.log('[WS] Connected');
    listeners.open.forEach((cb) => cb());

    const dialogId = localStorage.getItem('dialog_id');
    const clientRaw = localStorage.getItem('client');
    let clientId = null;

    try {
      clientId = clientRaw ? JSON.parse(clientRaw).id : null;
    } catch (e) {
      console.warn('[WS] Failed to parse client from localStorage');
    }

    if (dialogId && clientId) {
      sendEvent('init', {
        room_id: dialogId,
        client_id: clientId,
      });
    }
  };

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      listeners.message.forEach((cb) => cb(data));
    } catch (e) {
      console.warn('[WS] Non-JSON message:', event.data);
    }
  };

  socket.onerror = (e) => {
    console.error('[WS] Error', e);
    listeners.error.forEach((cb) => cb(e));
  };

  socket.onclose = (e) => {
    console.warn('[WS] Closed', e.code);
    listeners.close.forEach((cb) => cb(e));
    reconnectTimeout = setTimeout(() => {
      console.log('[WS] Reconnecting...');
      initSocket(token);
    }, 3000);
  };

  return socket;
}

export function getSocket() {
  return socket;
}

export function sendEvent(event, data) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({ event, data }));
  }
}

export function subscribe(eventType, callback) {
  if (listeners[eventType]) {
    listeners[eventType].push(callback);
  }
}

export function unsubscribe(eventType, callback) {
  if (listeners[eventType]) {
    listeners[eventType] = listeners[eventType].filter((cb) => cb !== callback);
  }
}
