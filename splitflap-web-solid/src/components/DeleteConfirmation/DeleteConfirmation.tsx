import styles from './DeleteConfirmation.module.css';

interface DeleteConfirmationProps {
  displayId: string;
  onConfirm: () => Promise<void>;
  onCancel: () => void;
}

export function DeleteConfirmation(props: DeleteConfirmationProps) {
  let isLoading = false;

  const handleConfirm = async () => {
    isLoading = true;
    try {
      await props.onConfirm();
    } finally {
      isLoading = false;
    }
  };

  return (
    <div class={styles.overlay}>
      <div class={styles.modal}>
        <h2>Delete Display?</h2>
        <p>
          Are you sure you want to delete <strong>{props.displayId}</strong>? This cannot be undone.
        </p>
        <div class={styles.actions}>
          <button onClick={props.onCancel} disabled={isLoading}>
            Cancel
          </button>
          <button onClick={handleConfirm} disabled={isLoading} class={styles.danger}>
            {isLoading ? 'Deleting...' : 'Delete'}
          </button>
        </div>
      </div>
    </div>
  );
}
