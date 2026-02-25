import { Show, Suspense, createEffect, createResource, createSignal } from 'solid-js';
import type { Display } from '../types/api';
import { DisplayForm } from './DisplayForm/DisplayForm';
import { DisplayPreview } from './DisplayPreview/DisplayPreview';
import { DisplaySidebar } from './DisplaySidebar/DisplaySidebar';
import { DeleteConfirmation } from './DeleteConfirmation/DeleteConfirmation';
import * as api from '../utils/api';
import styles from './App.module.css';

type FormMode = 'create' | 'edit';

export function App() {
  const [displayList, { refetch: refetchList }] = createResource(api.fetchDisplays);
  const [selectedId, setSelectedId] = createSignal<string | undefined>();
  const [formMode, setFormMode] = createSignal<FormMode | null>(null);
  const [editingDisplay, setEditingDisplay] = createSignal<Display | undefined>();
  const [deleteTarget, setDeleteTarget] = createSignal<Display | undefined>();
  const [successMessage, setSuccessMessage] = createSignal<string | null>(null);

  // Auto-select first display when list loads
  createEffect(() => {
    const displays = displayList();
    if (displays && displays.length > 0 && !selectedId()) {
      setSelectedId(displays[0].id);
    }
  });

  const selectedDisplay = () => {
    const displays = displayList();
    const id = selectedId();
    if (!displays || !id) return undefined;
    return displays.find((d) => d.id === id);
  };

  const handleCreateClick = () => {
    setFormMode('create');
    setEditingDisplay(undefined);
  };

  const handleEditClick = (display: Display) => {
    setFormMode('edit');
    setEditingDisplay(display);
  };

  const handleDeleteClick = (display: Display) => {
    setDeleteTarget(display);
  };

  const handleFormSubmit = async (display: Display) => {
    try {
      if (formMode() === 'create') {
        await api.createDisplay(display);
        showSuccess('Display created');
      } else {
        await api.updateDisplay(display.id, display);
        showSuccess('Display updated');
      }
      await refetchList();
      setFormMode(null);
      setEditingDisplay(undefined);
    } catch (error) {
      throw error;
    }
  };

  const handleDeleteConfirm = async () => {
    const target = deleteTarget();
    if (!target) return;

    try {
      await api.deleteDisplay(target.id);
      showSuccess('Display deleted');
      setDeleteTarget(undefined);
      await refetchList();
      if (selectedId() === target.id) {
        setSelectedId(undefined);
      }
    } catch (error) {
      throw error;
    }
  };

  const showSuccess = (message: string) => {
    setSuccessMessage(message);
    setTimeout(() => setSuccessMessage(null), 2000);
  };

  return (
    <div class={styles.app}>
      <Suspense fallback={<div class={styles.loading}>Loading...</div>}>
        <DisplaySidebar
          displays={displayList() ?? []}
          selectedId={selectedId()}
          onSelect={setSelectedId}
          onCreateClick={handleCreateClick}
          onEditClick={handleEditClick}
          onDeleteClick={handleDeleteClick}
        />

        <div class={styles.main}>
          <DisplayPreview display={selectedDisplay()} />

          <Show when={successMessage()}>
            <div class={styles.success}>✓ {successMessage()}</div>
          </Show>
        </div>

        {/* Form Modal */}
        <Show when={formMode()}>
          <div class={styles.modalOverlay}>
            <DisplayForm
              display={editingDisplay()}
              onSubmit={handleFormSubmit}
              onCancel={() => setFormMode(null)}
              onSuccess={() => setFormMode(null)}
            />
          </div>
        </Show>

        {/* Delete Confirmation Modal */}
        <Show when={deleteTarget()}>
          {(target) => (
            <DeleteConfirmation
              displayId={target().id}
              onConfirm={handleDeleteConfirm}
              onCancel={() => setDeleteTarget(undefined)}
            />
          )}
        </Show>
      </Suspense>
    </div>
  );
}
