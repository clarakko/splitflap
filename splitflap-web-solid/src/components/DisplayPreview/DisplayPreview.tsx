import { Show, Suspense } from 'solid-js';
import { createSignal } from 'solid-js';
import { createFetchDisplay } from '../../hooks/createFetchDisplay';
import { SplitFlapDisplay } from '../SplitFlapDisplay/SplitFlapDisplay';
import styles from './DisplayPreview.module.css';

interface DisplayPreviewProps {
  displayId?: string;
}

export function DisplayPreview(props: DisplayPreviewProps) {
  const [displayId] = createSignal(props.displayId ?? 'demo');
  const { display } = createFetchDisplay(displayId);

  return (
    <div class={styles.container}>
      <h1>Display Preview</h1>

      <Suspense fallback={<div class={styles.loading}>Loading display data...</div>}>
        <Show
          when={display()}
          fallback={
            <Show when={display.error}>
              <div class={styles.error}>Error: {display.error?.message}</div>
            </Show>
          }
        >
          {(data) => (
            <div class={styles.content}>
              <SplitFlapDisplay display={data()} />
            </div>
          )}
        </Show>
      </Suspense>
    </div>
  );
}
