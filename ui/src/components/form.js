import styles from './form.module.css'

export default function Form({ children, title, disabled, errorMessage, onSubmit }) {
  return (
    <form className={styles.form} onSubmit={onSubmit}>
      <div className="title">{title}</div>
      {errorMessage !== null &&
        <div className={styles.error}>
          Error:{' '}{errorMessage}
        </div>
      }
      {children}
      <button
        className="button"
        type="submit"
        disabled={disabled}
      >
        Login
      </button>
    </form>
  )
}
