import { For, ParentProps } from 'solid-js';
import { Display } from '../../types/api';
import { SplitFlapCell } from '../SplitFlapCell/SplitFlapCell';
import styles from './SplitFlapDisplay.module.css';

interface SplitFlapDisplayProps extends ParentProps {
  display: Display;
  flipDuration?: number; // milliseconds per character in flip sequence
}

/**
 * Grid display component that renders multiple split-flap cells.
 *
 * Takes a Display object and renders animated cells for each character.
 */
export function SplitFlapDisplay(props: SplitFlapDisplayProps) {
  const flipDuration = () => props.flipDuration ?? 100;

  return (
    <div class={styles.display}>
      <For each={props.display.content.rows}>
        {(row) => (
          <div class={styles.row}>
            <For each={row}>
              {(char) => (
                <SplitFlapCell
                  targetChar={char || ' '}
                  flipDuration={flipDuration()}
                />
              )}
            </For>
          </div>
        )}
      </For>
    </div>
  );
}
