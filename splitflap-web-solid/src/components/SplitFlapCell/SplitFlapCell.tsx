import { createSignal, createEffect, createMemo, ParentProps, onCleanup } from 'solid-js';
import { calculateFlipPath } from '../../utils/flipCalculator';
import styles from './SplitFlapCell.module.css';

interface SplitFlapCellProps extends ParentProps {
  targetChar: string;
  flipDuration?: number; // milliseconds per character in the flip sequence
  onFlipComplete?: () => void;
}

/**
 * Animated split-flap cell component.
 *
 * Displays a single character and animates transitions using a flip sequence.
 */
export function SplitFlapCell(props: SplitFlapCellProps) {
  const flipDuration = () => props.flipDuration ?? 100;

  const [displayChar, setDisplayChar] = createSignal(' ');
  const [isFlipping, setIsFlipping] = createSignal(false);

  // Create a memo to track the target char as a reactive value
  const targetChar = createMemo(() => props.targetChar);

  // Trigger flip animation when targetChar changes
  createEffect(() => {
    const target = targetChar();

    if (displayChar() === target) {
      return; // Already at target
    }

    const flipPath = calculateFlipPath(displayChar(), target);

    if (flipPath.length <= 1) {
      // No flipping needed
      setDisplayChar(target);
      return;
    }

    setIsFlipping(true);
    let step = 0;

    const interval = setInterval(() => {
      if (step >= flipPath.length) {
        clearInterval(interval);
        setIsFlipping(false);
        props.onFlipComplete?.();
        return;
      }

      setDisplayChar(flipPath[step]);
      step++;
    }, flipDuration());

    // Cleanup: clear interval if effect re-runs or component unmounts
    onCleanup(() => clearInterval(interval));
  });

  return (
    <div class={styles.cell} classList={{ [styles.flipping]: isFlipping() }} data-char={displayChar()}>
      {displayChar()}
    </div>
  );
}
