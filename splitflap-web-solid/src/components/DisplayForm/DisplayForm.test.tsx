import { render, fireEvent, waitFor } from '@solidjs/testing-library';
import { describe, expect, test, vi } from 'vitest';
import { DisplayForm } from './DisplayForm';
import type { Display } from '../../types/api';

describe('<DisplayForm />', () => {
  const mockDisplay: Display = {
    id: 'test-display',
    content: {
      rows: [
        ['H', 'E', 'L', 'L', 'O'],
        ['W', 'O', 'R', 'L', 'D'],
      ],
    },
    config: {
      rowCount: 2,
      columnCount: 5,
    },
  };

  test('renders in create mode with empty grid', () => {
    const onSubmit = vi.fn();
    const onCancel = vi.fn();
    const { getByText } = render(() => (
      <DisplayForm onSubmit={onSubmit} onCancel={onCancel} />
    ));

    expect(getByText('Create Display')).toBeInTheDocument();
  });

  test('renders in edit mode with display data', () => {
    const onSubmit = vi.fn();
    const onCancel = vi.fn();
    const { getByText, container } = render(() => (
      <DisplayForm display={mockDisplay} onSubmit={onSubmit} onCancel={onCancel} />
    ));

    expect(getByText('Edit Display')).toBeInTheDocument();
    const idInput = container.querySelector('#id') as HTMLInputElement;
    expect(idInput).toBeDisabled();
    expect(idInput.value).toBe('test-display');
  });

  test('initializes grid inputs with correct values in edit mode', async () => {
    const onSubmit = vi.fn();
    const onCancel = vi.fn();
    const { container } = render(() => (
      <DisplayForm display={mockDisplay} onSubmit={onSubmit} onCancel={onCancel} />
    ));

    // Wait for DOM to be properly initialized
    await waitFor(() => {
      const inputs = container.querySelectorAll('input[maxLength="1"]');
      expect(inputs.length).toBe(10); // 2 rows × 5 columns
    });

    const inputs = container.querySelectorAll('input[maxLength="1"]');
    // Check first row: HELLO
    expect((inputs[0] as HTMLInputElement).value).toBe('H');
    expect((inputs[1] as HTMLInputElement).value).toBe('E');
    expect((inputs[2] as HTMLInputElement).value).toBe('L');
    expect((inputs[3] as HTMLInputElement).value).toBe('L');
    expect((inputs[4] as HTMLInputElement).value).toBe('O');

    // Check second row: WORLD
    expect((inputs[5] as HTMLInputElement).value).toBe('W');
    expect((inputs[6] as HTMLInputElement).value).toBe('O');
    expect((inputs[7] as HTMLInputElement).value).toBe('R');
    expect((inputs[8] as HTMLInputElement).value).toBe('L');
    expect((inputs[9] as HTMLInputElement).value).toBe('D');
  });

  test('allows typing in uncontrolled inputs', async () => {
    const onSubmit = vi.fn();
    const onCancel = vi.fn();
    const { container } = render(() => (
      <DisplayForm display={mockDisplay} onSubmit={onSubmit} onCancel={onCancel} />
    ));

    await waitFor(() => {
      const inputs = container.querySelectorAll('input[maxLength="1"]');
      expect(inputs.length).toBe(10);
    });

    const inputs = container.querySelectorAll('input[maxLength="1"]');

    // Simulate user typing - directly set input values (uncontrolled)
    const firstInput = inputs[0] as HTMLInputElement;
    firstInput.value = 'A';
    fireEvent.input(firstInput);
    
    // Value should persist in DOM
    expect(firstInput.value).toBe('A');

    // Edit multiple cells
    (inputs[2] as HTMLInputElement).value = 'X';
    (inputs[5] as HTMLInputElement).value = 'Q';

    expect((inputs[2] as HTMLInputElement).value).toBe('X');
    expect((inputs[5] as HTMLInputElement).value).toBe('Q');
  });

  test('submits correct data after editing cells', async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);
    const onCancel = vi.fn();
    const { container, getByText } = render(() => (
      <DisplayForm display={mockDisplay} onSubmit={onSubmit} onCancel={onCancel} />
    ));

    await waitFor(() => {
      const inputs = container.querySelectorAll('input[maxLength="1"]');
      expect(inputs.length).toBe(10);
    });

    const inputs = container.querySelectorAll('input[maxLength="1"]');

    // User types new values directly into inputs
    (inputs[0] as HTMLInputElement).value = 'A';
    (inputs[5] as HTMLInputElement).value = 'X';
    // Leave others unchanged

    // Submit form
    const submitButton = getByText('Update');
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledTimes(1);
    });

    // Verify submitted data reflects the DOM values
    const submittedDisplay = onSubmit.mock.calls[0][0];
    expect(submittedDisplay.content.rows[0][0]).toBe('A'); // Changed from 'H'
    expect(submittedDisplay.content.rows[1][0]).toBe('X'); // Changed from 'W'
    expect(submittedDisplay.content.rows[0][1]).toBe('E'); // Unchanged
    expect(submittedDisplay.content.rows[0][2]).toBe('L'); // Unchanged
  });

  test('creates display with correct grid size and submitted values', async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);
    const onCancel = vi.fn();
    const { container, getByText } = render(() => (
      <DisplayForm onSubmit={onSubmit} onCancel={onCancel} />
    ));

    // Set ID and grid size
    const idInput = container.querySelector('#id') as HTMLInputElement;
    fireEvent.input(idInput, { target: { value: 'new-display' } });

    const rowsInput = container.querySelector('#rows') as HTMLInputElement;
    fireEvent.input(rowsInput, { target: { value: '2' } });

    const colsInput = container.querySelector('#columns') as HTMLInputElement;
    fireEvent.input(colsInput, { target: { value: '3' } });

    await waitFor(() => {
      const inputs = container.querySelectorAll('input[maxLength="1"]');
      expect(inputs.length).toBe(6); // 2 rows × 3 columns
    });

    const inputs = container.querySelectorAll('input[maxLength="1"]');
    // Type some values
    (inputs[0] as HTMLInputElement).value = 'T';
    (inputs[1] as HTMLInputElement).value = 'E';

    // Submit
    const submitButton = getByText('Create');
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledTimes(1);
    });

    const submittedDisplay = onSubmit.mock.calls[0][0];
    expect(submittedDisplay.id).toBe('new-display');
    expect(submittedDisplay.config.rowCount).toBe(2);
    expect(submittedDisplay.config.columnCount).toBe(3);
    expect(submittedDisplay.content.rows.length).toBe(2);
    expect(submittedDisplay.content.rows[0].length).toBe(3);
    expect(submittedDisplay.content.rows[0][0]).toBe('T');
    expect(submittedDisplay.content.rows[0][1]).toBe('E');
    // Unfilled cells should be empty strings
    expect(submittedDisplay.content.rows[0][2]).toBe('');
  });

  test('displays correct form title and button text based on mode', () => {
    const onSubmit = vi.fn();
    const onCancel = vi.fn();

    const { getByText: getByTextCreate } = render(() => (
      <DisplayForm onSubmit={onSubmit} onCancel={onCancel} />
    ));
    expect(getByTextCreate('Create Display')).toBeInTheDocument();
    expect(getByTextCreate('Create')).toBeInTheDocument();

    const { getByText: getByTextEdit } = render(() => (
      <DisplayForm display={mockDisplay} onSubmit={onSubmit} onCancel={onCancel} />
    ));
    expect(getByTextEdit('Edit Display')).toBeInTheDocument();
    expect(getByTextEdit('Update')).toBeInTheDocument();
  });
});
