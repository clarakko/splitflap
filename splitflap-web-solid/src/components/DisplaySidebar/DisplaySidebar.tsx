import { For, Show, createSignal } from 'solid-js';
import type { Display } from '../../types/api';
import styles from './DisplaySidebar.module.css';

interface DisplaySidebarProps {
  displays: Display[];
  selectedId?: string;
  onSelect: (id: string) => void;
  onCreateClick: () => void;
  onEditClick: (display: Display) => void;
  onDeleteClick: (display: Display) => void;
}

export function DisplaySidebar(props: DisplaySidebarProps) {
  const [hoverId, setHoverId] = createSignal<string | null>(null);

  return (
    <div class={styles.sidebar}>
      <div class={styles.header}>
        <h1>Displays</h1>
        <button onClick={props.onCreateClick} class={styles.createBtn} title="Create new display">
          +
        </button>
      </div>

      <Show
        when={props.displays.length > 0}
        fallback={
          <div class={styles.empty}>
            <p>No displays yet</p>
            <button onClick={props.onCreateClick}>Create First Display</button>
          </div>
        }
      >
        <div class={styles.list}>
          <For each={props.displays}>
            {(display) => (
              <div
                class={`${styles.item} ${props.selectedId === display.id ? styles.selected : ''}`}
                onMouseEnter={() => setHoverId(display.id)}
                onMouseLeave={() => setHoverId(null)}
              >
                <div class={styles.displayInfo} onClick={() => props.onSelect(display.id)}>
                  <div class={styles.displayId}>{display.id}</div>
                  <div class={styles.displaySize}>
                    {display.config.rowCount}×{display.config.columnCount}
                  </div>
                </div>
                <Show when={hoverId() === display.id}>
                  <div class={styles.actions}>
                    <button
                      onClick={() => props.onEditClick(display)}
                      class={styles.iconBtn}
                      title="Edit"
                    >
                      ✏️
                    </button>
                    <button
                      onClick={() => props.onDeleteClick(display)}
                      class={styles.iconBtn}
                      title="Delete"
                    >
                      🗑️
                    </button>
                  </div>
                </Show>
              </div>
            )}
          </For>
        </div>
      </Show>
    </div>
  );
}
