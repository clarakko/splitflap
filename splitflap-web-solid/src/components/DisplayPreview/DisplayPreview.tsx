import type { Display } from '../../types/api';
import { SplitFlapDisplay } from '../SplitFlapDisplay/SplitFlapDisplay';
import styles from './DisplayPreview.module.css';

interface DisplayPreviewProps {
  display?: Display;
}

export function DisplayPreview(props: DisplayPreviewProps) {
  return (
    <div class={styles.container}>
      {props.display ? (
        <div class={styles.content}>
          <h2>{props.display.id}</h2>
          <SplitFlapDisplay display={props.display} />
        </div>
      ) : (
        <div class={styles.empty}>
          <p>Select a display to view</p>
        </div>
      )}
    </div>
  );
}
