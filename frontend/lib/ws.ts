let ws: WebSocket | null = null
let subscribers: ((msg: string) => void)[] = []

export function connectWS(url: string = "ws://localhost:8080/stream") {
  if (!ws) {
    ws = new WebSocket(url)
    ws.onmessage = (e) => {
      for (const fn of subscribers) fn(e.data)
    }
  }
}

export function subscribeWS(cb: (msg: string) => void) {
  subscribers.push(cb)
}
