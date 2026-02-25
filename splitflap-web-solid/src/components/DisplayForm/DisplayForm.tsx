import { Show, createSignal } from 'solid-js';
import type { Display, DisplayContent, DisplayConfig } from '../../types/api';
import styles from './DisplayForm.module.css';

interface DisplayFormProps {
  display?: Display; // If provided, this is edit mode; if undefined, this is create mode
  onSubmit: (display: Display) => Promise<void>;
  onCancel: () => void;
  onSuccess?: () => void;
}

/**
 * Unified form for creating and editing displays.
 * If `display` prop is provided, renders in edit mode.
 * Otherwise, renders in create mode with empty form.
 */
export function DisplayForm(props: DisplayFormProps) {
  const isEditMode = () => !!props.display;

  const [id, setId] = createSignal(props.display?.id ?? '');
  const [rowCount, setRowCount] = createSignal(props.display?.config.rowCount ?? 1);
  const [columnCount, setColumnCount] = createSignal(props.display?.config.columnCount ?? 1);
  const [rows, setRows] = createSignal<string[][]>(props.display?.content.rows ?? [['']]);
  const [error, setError] = createSignal<string | null>(null);
  const [isLoading, setIsLoading] = createSignal(false);

  // When grid size changes, rebuild the grid
  const updateGridSize = (newRowCount: number, newColumnCount: number) => {
    const currentRows = rows();

    // Resize rows
    let newRows = currentRows.slice(0, newRowCount);
    while (newRows.length < newRowCount) {
      newRows.push([]);
    }

    // Resize each row
    newRows = newRows.map((row) => {
      let newRow = row.slice(0, newColumnCount);
      while (newRow.length < newColumnCount) {
        newRow.push(' ');
      }
      return newRow;
    });

    setRows(newRows);
    setRowCount(newRowCount);
    setColumnCount(newColumnCount);
  };

  const updateCell = (rowIdx: number, colIdx: number, char: string) => {
    const newRows = rows().map((r) => [...r]);
    newRows[rowIdx][colIdx] = char.slice(0, 1) || ' '; // Take only first char
    setRows(newRows);
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      const display: Display = {
        id: id(),
        content: {
          rows: rows(),
        },
        config: {
          rowCount: rowCount(),
          columnCount: columnCount(),
        },
      };

      await props.onSubmit(display);
      props.onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} class={styles.form}>
      <h2>{isEditMode() ? 'Edit Display' : 'Create Display'}</h2>

      <Show when={error()}>
        <div class={styles.error}>{error()}</div>
      </Show>

      {/* ID Field */}
      <div class={styles.field}>
        <label htmlFor="id">Display ID</label>
        <input
          id="id"
          type="text"
          value={id()}
          onInput={(e) => setId(e.currentTarget.value)}
          disabled={isEditMode()}
          required
          placeholder="e.g., my-display"
        />
        {isEditMode() && <small>ID cannot be changed</small>}
      </div>

      {/* Grid Size Fields */}
      <div class={styles.row}>
        <div class={styles.field}>
          <label htmlFor="rows">Rows</label>
          <input
            id="rows"
            type="number"
            min="1"
            max="20"
            value={rowCount()}
            onInput={(e) => updateGridSize(parseInt(e.currentTarget.value), columnCount())}
            required
          />
        </div>
        <div class={styles.field}>
          <label htmlFor="columns">Columns</label>
          <input
            id="columns"
            type="number"
            min="1"
            max="10"
            value={columnCount()}
            onInput={(e) => updateGridSize(rowCount(), parseInt(e.currentTarget.value))}
            required
          />
        </div>
      </div>

      {/* Grid Editor */}
      <div class={styles.field}>
        <label>Content</label>
        <div class={styles.grid}>
          {rows().map((row, rowIdx) => (
            <div key={rowIdx} class={styles.gridRow}>
              {row.map((char, colIdx) => (
                <input
                  key={`${rowIdx}-${colIdx}`}
                  type="text"
                  maxLength="1"
                  value={char}
                  onInput={(e) => updateCell(rowIdx, colIdx, e.currentTarget.value)}
                  class={styles.gridCell}
                  placeholder=" "
                />
              ))}
            </div>
          ))}
        </div>
      </div>

      {/* Actions */}
      <div class={styles.actions}>
        <button type="button" onClick={props.onCancel} disabled={isLoading()}>
          Cancel
        </button>
        <button type="submit" disabled={isLoading()} class={styles.primary}>
          {isLoading() ? 'Saving...' : isEditMode() ? 'Update' : 'Create'}
        </button>
      </div>
    </form>
  );
}
