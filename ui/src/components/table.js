import styles from './table.module.css'

export default function Table({ headers, rows }) {

  let colTemplate = headers.map(row => {
    if (row.expand) {
      return '1fr'
    }
    return 'auto'
  }).join(' ')

  let rowToggle = false

  return (
    <div
      className={styles.table}
      style={{ gridTemplateColumns: colTemplate }}
    >
      {headers.map((header, i) =>
        <div key={i} className={`${styles.cell} ${styles.header}`}>
          {header.title}
        </div>
      )}
      {rows.length ?
        rows.map((row, i) => {
          rowToggle = !rowToggle
          return headers.map((header, j) =>
            <div
              key={[i, j]}
              className={styles.cell + (rowToggle ? ` ${styles.striped}` : '')}>
              {header.render(row)}
            </div>
          )
        }) :
        <div className={`${styles.cell} ${styles.striped} ${styles.empty}`}>
          This table is currently empty
        </div>
      }
    </div>
  )
}
