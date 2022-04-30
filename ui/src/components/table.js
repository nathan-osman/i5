import styles from './table.module.css'

export default function Table({ headers, rows }) {

  // TODO: offer options for sorting, etc.

  return (
    <table className={styles.table}>
      <thead>
        <tr>
          {headers.map(header =>
            <th key={header.title}>{header.title}</th>
          )}
        </tr>
      </thead>
      <tbody>
        {rows.length ?
          rows.map((row, i) =>
            <tr key={i}>
              {headers.map((header, i) =>
                <td key={i}>
                  {header.render(row)}
                </td>
              )}
            </tr>
          ) :
          <tr>
            <td colSpan={headers.length} className={styles.empty}>
              <span className="secondary">
                Table is currently empty
              </span>
            </td>
          </tr>
        }
      </tbody>
    </table>
  )
}
