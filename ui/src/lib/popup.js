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

  return (
    <PopupContext.Provider value={popup}>
      {children}
      <div className={styles.container}>
        {messages.map((message, i) =>
          <div key={i} className={`${styles.message} ${styles[message.type]}`}>
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
