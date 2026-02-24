import './App.css'
import { useFetchDisplay } from './hooks/useFetchDisplay'

function App() {
  const { data, loading, error } = useFetchDisplay('demo')

  if (loading) {
    return (
      <div className="container">
        <h1>SplitFlap Display</h1>
        <p>Loading display data...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container">
        <h1>SplitFlap Display</h1>
        <div className="error">
          <h2>Error</h2>
          <p>{error.message}</p>
          {error.statusCode && <p>Status Code: {error.statusCode}</p>}
        </div>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="container">
        <h1>SplitFlap Display</h1>
        <p>No display data available</p>
      </div>
    )
  }

  return (
    <div className="container">
      <h1>SplitFlap Display</h1>
      
      <div className="display-info">
        <h2>Display ID: {data.id}</h2>
        <p>Rows: {data.config.rowCount} | Columns: {data.config.columnCount}</p>
      </div>

      <div className="display-content">
        <h3>Content:</h3>
        <table>
          <tbody>
            {data.content.rows.map((row, rowIndex) => (
              <tr key={rowIndex}>
                {row.map((cell, cellIndex) => (
                  <td key={cellIndex}>{cell}</td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default App
