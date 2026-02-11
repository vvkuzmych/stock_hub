import { useState, useEffect, useRef } from 'react'
import './App.css'

function App() {
  const [user, setUser] = useState(null)
  const [username, setUsername] = useState('')
  const [messages, setMessages] = useState([])
  const [inputMessage, setInputMessage] = useState('')
  const [isConnected, setIsConnected] = useState(false)
  const [orders, setOrders] = useState([])
  const [newOrder, setNewOrder] = useState({
    symbol: 'AAPL',
    orderType: 'bid',
    price: '',
    quantity: ''
  })
  
  const wsRef = useRef(null)
  const reconnectAttempts = useRef(0)
  const maxReconnectAttempts = 5

  useEffect(() => {
    connect()
    return () => {
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [])

  const connect = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.hostname}:${window.location.port || '8082'}/ws`
    
    try {
      const ws = new WebSocket(wsUrl)
      wsRef.current = ws

      ws.onopen = () => {
        setIsConnected(true)
        reconnectAttempts.current = 0
      }

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data)
        handleMessage(data)
      }

      ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }

      ws.onclose = () => {
        setIsConnected(false)
        if (reconnectAttempts.current < maxReconnectAttempts) {
          reconnectAttempts.current++
          setTimeout(connect, 2000)
        }
      }
    } catch (error) {
      console.error('Failed to connect:', error)
      setIsConnected(false)
    }
  }

  const handleMessage = (data) => {
    switch (data.type) {
      case 'registered':
        setUser(data.data)
        break
      case 'orders_update':
        setOrders(data.data || [])
        break
      case 'order_created':
        setOrders(prev => {
          const exists = prev.find(o => o.id === data.data.id)
          if (exists) return prev
          return [...prev, data.data].sort((a, b) => {
            if (a.symbol !== b.symbol) return a.symbol.localeCompare(b.symbol)
            if (a.order_type !== b.order_type) return a.order_type === 'bid' ? -1 : 1
            return a.order_type === 'bid' ? b.price - a.price : a.price - b.price
          })
        })
        break
      case 'order_cancelled':
        setOrders(prev => prev.filter(o => o.id !== data.data.order_id))
        break
      case 'message':
        // Always add message - determine color based on sender
        addMessage(data.data.username, data.data.content)
        break
      case 'error':
        alert(data.data)
        break
    }
  }

  const register = () => {
    if (!username.trim()) {
      alert('Please enter a username')
      return
    }

    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({
        type: 'register',
        data: { username: username.trim() }
      }))
    }
  }

  const sendMessage = () => {
    const message = inputMessage.trim()
    if (!message) return

    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({
        type: 'message',
        data: message
      }))
      setInputMessage('')
      // Message will be added when received from server
    }
  }

  const createOrder = () => {
    if (!newOrder.price || !newOrder.quantity) {
      alert('Please enter price and quantity')
      return
    }

    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({
        type: 'order',
        data: {
          symbol: newOrder.symbol,
          order_type: newOrder.orderType,
          price: parseFloat(newOrder.price),
          quantity: parseInt(newOrder.quantity)
        }
      }))
      setNewOrder({ ...newOrder, price: '', quantity: '' })
    }
  }

  const cancelOrder = (orderId) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({
        type: 'cancel_order',
        data: { order_id: orderId }
      }))
    }
  }

  const getMessageColor = (sender) => {
    if (user && sender === user.username) {
      return 'sent' // Your messages - green
    }
    // Generate consistent color for each user based on their name
    const colors = ['received-blue', 'received-purple', 'received-orange', 'received-teal', 'received-pink']
    let hash = 0
    for (let i = 0; i < sender.length; i++) {
      hash = sender.charCodeAt(i) + ((hash << 5) - hash)
    }
    return colors[Math.abs(hash) % colors.length]
  }

  const addMessage = (sender, text, type) => {
    const newMessage = {
      id: Date.now() + Math.random(),
      sender,
      text,
      type: type || getMessageColor(sender),
      time: new Date().toLocaleTimeString(),
    }
    setMessages((prev) => [...prev, newMessage])
  }

  const groupedOrders = orders.reduce((acc, order) => {
    if (!acc[order.symbol]) {
      acc[order.symbol] = { bids: [], asks: [] }
    }
    if (order.order_type === 'bid') {
      acc[order.symbol].bids.push(order)
    } else {
      acc[order.symbol].asks.push(order)
    }
    return acc
  }, {})

  if (!user) {
    return (
      <div className="container">
        <div className="header">
          <h1>🚀 Stock Hub</h1>
          <p>Stock Trading Platform</p>
          <div className={`status ${isConnected ? 'connected' : 'disconnected'}`}>
            {isConnected ? 'Connected' : 'Disconnected'}
          </div>
        </div>
        <div className="register-panel">
          <h2>Register to Start Trading</h2>
          <div className="input-group">
            <label htmlFor="username">Username:</label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && register()}
              placeholder="Enter your username"
              disabled={!isConnected}
            />
          </div>
          <button onClick={register} disabled={!isConnected || !username.trim()}>
            Register
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="container">
      <div className="header">
        <h1>🚀 Stock Hub</h1>
        <p>Welcome, {user.username}!</p>
        <div className={`status ${isConnected ? 'connected' : 'disconnected'}`}>
          {isConnected ? 'Connected' : 'Disconnected'}
        </div>
      </div>

      <div className="content-grid">
        <div className="panel trading-panel">
          <h2>📈 Place Order</h2>
          <div className="order-form">
            <div className="input-group">
              <label>Symbol:</label>
              <input
                type="text"
                value={newOrder.symbol}
                onChange={(e) => setNewOrder({ ...newOrder, symbol: e.target.value.toUpperCase() })}
                placeholder="AAPL"
                maxLength="5"
              />
            </div>
            <div className="input-group">
              <label>Type:</label>
              <select
                value={newOrder.orderType}
                onChange={(e) => setNewOrder({ ...newOrder, orderType: e.target.value })}
              >
                <option value="bid">Bid (Buy)</option>
                <option value="ask">Ask (Sell)</option>
              </select>
            </div>
            <div className="input-group">
              <label>Price:</label>
              <input
                type="number"
                step="0.01"
                value={newOrder.price}
                onChange={(e) => setNewOrder({ ...newOrder, price: e.target.value })}
                placeholder="100.00"
              />
            </div>
            <div className="input-group">
              <label>Quantity:</label>
              <input
                type="number"
                value={newOrder.quantity}
                onChange={(e) => setNewOrder({ ...newOrder, quantity: e.target.value })}
                placeholder="10"
              />
            </div>
            <button onClick={createOrder} className="order-btn">
              Place {newOrder.orderType === 'bid' ? 'Bid' : 'Ask'}
            </button>
          </div>
        </div>

        <div className="panel chat-panel">
          <h2>💬 Chat</h2>
          <div className="messages">
            {messages.length === 0 ? (
              <div className="empty-state">No messages yet</div>
            ) : (
              messages.map((msg) => {
                const messageClass = user && msg.sender === user.username ? 'sent' : getMessageColor(msg.sender)
                return (
                  <div key={msg.id} className={`message ${messageClass}`}>
                    <div className="message-time">[{msg.time}] {msg.sender}</div>
                    <div>{msg.text}</div>
                  </div>
                )
              })
            )}
          </div>
          <div className="chat-input">
            <input
              type="text"
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
              placeholder="Type a message..."
            />
            <button onClick={sendMessage}>Send</button>
          </div>
        </div>
      </div>

      <div className="orders-section">
        <div className="panel orders-panel">
          <h2>📊 Order Book</h2>
          <div className="order-book">
            {Object.keys(groupedOrders).length === 0 ? (
              <div className="empty-state">No open orders</div>
            ) : (
              Object.entries(groupedOrders).map(([symbol, { bids, asks }]) => (
                <div key={symbol} className="symbol-group">
                  <h3>{symbol}</h3>
                  <div className="orders-container">
                    <div className="bids-section">
                      <h4>Bids (Buy)</h4>
                      {bids.length === 0 ? (
                        <div className="empty-orders">No bids</div>
                      ) : (
                        bids.map(order => (
                          <div key={order.id} className={`order-item bid ${order.user_id === user.id ? 'my-order' : ''}`}>
                            <div className="order-info">
                              <span className="order-price">{order.price.toFixed(2)}</span>
                              <span className="order-qty">{order.quantity}</span>
                              <span className="order-user">{order.username}</span>
                            </div>
                            {order.user_id === user.id && (
                              <button className="cancel-btn" onClick={() => cancelOrder(order.id)}>Cancel</button>
                            )}
                          </div>
                        ))
                      )}
                    </div>
                    <div className="asks-section">
                      <h4>Asks (Sell)</h4>
                      {asks.length === 0 ? (
                        <div className="empty-orders">No asks</div>
                      ) : (
                        asks.map(order => (
                          <div key={order.id} className={`order-item ask ${order.user_id === user.id ? 'my-order' : ''}`}>
                            <div className="order-info">
                              <span className="order-price">{order.price.toFixed(2)}</span>
                              <span className="order-qty">{order.quantity}</span>
                              <span className="order-user">{order.username}</span>
                            </div>
                            {order.user_id === user.id && (
                              <button className="cancel-btn" onClick={() => cancelOrder(order.id)}>Cancel</button>
                            )}
                          </div>
                        ))
                      )}
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default App
