import { createContext, useContext, useState } from 'react'
import styles from './popup.module.css'

const PT_ERROR = 'error'

const PopupContext = createContext(null)

function PopupProvider({ children }) {

  const [messages, setMessages] = useState([])

  function show(message) {
    setMessages([...messages, message])
  }

  const popup = {
    error: text => show({ type: PT_ERROR, text: `Error: ${text}` })
  }

  function handleClick(i) {
    setMessages([
      ...messages.slice(0, i),
      ...messages.slice(i + 1)
    ])
  }

  return (
    <PopupContext.Provider value={popup}>
      {children}
      <div className={styles.container}>
        {messages.map((message, i) =>
          <div
            key={i}
            className={`${styles.message} ${styles[message.type]}`}
            onClick={() => handleClick(i)}
          >
            {message.text}
          </div>
        )}
      </div>
    </PopupContext.Provider>
  )
}

function usePopup() {
  return useContext(PopupContext)
}

export { PopupProvider, usePopup }
