import { Show, createSignal, createEffect, For, Index } from 'solid-js';
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
  // Separate signal for editing - this is what the form actually uses
  const [editingRows, setEditingRows] = createSignal<string[][]>(props.display?.content.rows ?? [['']]);
  const [error, setError] = createSignal<string | null>(null);
  const [isLoading, setIsLoading] = createSignal(false);
  
  let gridRef: HTMLDivElement | undefined;

  // Normalize rows: convert space-only cells to empty strings
  const normalizeRows = (rows: string[][]): string[][] => {
    return rows.map(row => row.map(cell => (cell === ' ' ? '' : cell)));
  };

  // Sync signals when props.display changes (for edit mode)
  createEffect(() => {
    if (props.display) {
      const normalizedRows = normalizeRows(props.display.content.rows);
      setId(props.display.id);
      setRowCount(props.display.config.rowCount);
      setColumnCount(props.display.config.columnCount);
      setRows(normalizedRows);
      setEditingRows(normalizedRows);
      
      // When display changes, reset the form inputs to show new values
      // Use queueMicrotask to wait for DOM updates
      queueMicrotask(() => {
        if (gridRef) {
          const cells = gridRef.querySelectorAll('input[maxLength="1"]');
          const rows = normalizedRows;
          const cols = props.display!.config.columnCount;
          cells.forEach((cell, idx) => {
            const input = cell as HTMLInputElement;
            const rowIdx = Math.floor(idx / cols);
            const colIdx = idx % cols;
            input.value = rows[rowIdx]?.[colIdx] ?? '';
          });
        }
      });
    }
  });

  // When grid size changes, rebuild the grid
  const updateGridSize = (newRowCount: number, newColumnCount: number) => {
    const currentRows = editingRows();

    // Resize rows
    let newRows = currentRows.slice(0, newRowCount);
    while (newRows.length < newRowCount) {
      // Initialize new rows with empty strings (not spaces)
      newRows.push(Array(newColumnCount).fill(''));
    }

    // Resize each row
    newRows = newRows.map((row) => {
      let newRow = row.slice(0, newColumnCount);
      while (newRow.length < newColumnCount) {
        newRow.push('');
      }
      return newRow;
    });

    setEditingRows(newRows);
    setRowCount(newRowCount);
    setColumnCount(newColumnCount);
  };

  const updateCell = (rowIdx: number, colIdx: number, char: string) => {
    const newRows = editingRows().map((r) => [...r]);
    // Only take first character, allow empty string
    newRows[rowIdx][colIdx] = char.slice(0, 1);
    setEditingRows(newRows);
  };

  // Read current grid values from DOM inputs - for uncontrolled input pattern
  const readGridFromDOM = (): string[][] => {
    if (!gridRef) return editingRows();
    
    const cells = gridRef.querySelectorAll('input[maxLength="1"]');
    const gridRows = rowCount();
    const gridCols = columnCount();
    const newRows: string[][] = [];
    
    for (let r = 0; r < gridRows; r++) {
      newRows[r] = [];
      for (let c = 0; c < gridCols; c++) {
        const cellIndex = r * gridCols + c;
        const input = cells[cellIndex] as HTMLInputElement | undefined;
        const value = input?.value ?? ''; // Use empty string, not space
        newRows[r].push(value.slice(0, 1));
      }
    }
    
    return newRows;
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);

    try {
      // Read current values from DOM instead of depending on signals
      let gridContent = readGridFromDOM();
      // Normalize: convert spaces to empty strings
      gridContent = normalizeRows(gridContent);
      
      const display: Display = {
        id: id(),
        content: {
          rows: gridContent,
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
        <div class={styles.grid} ref={gridRef}>
          <For each={editingRows()}>
            {(row, getRowIdx) => {
              return (
                <div class={styles.gridRow}>
                  <For each={row}>
                    {(char) => (
                      <input
                        type="text"
                        maxLength="1"
                        value={char}
                        class={styles.gridCell}
                        placeholder=" "
                      />
                    )}
                  </For>
                </div>
              );
            }}
          </For>
        </div>
      </div>

      {/* Actions */}
      <div class={styles.actions}>
        <button type="button" onClick={props.onCancel} disabled={isLoading()} class={styles.button}>
          Cancel
        </button>
        <button type="submit" disabled={isLoading()} class={`${styles.button} ${styles.primary}`}>
          {isLoading() ? 'Saving...' : isEditMode() ? 'Update' : 'Create'}
        </button>
      </div>
    </form>
  );
}
