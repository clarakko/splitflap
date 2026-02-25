import type { Display } from '../types/api';

const API_BASE = '/api/v1';

export async function fetchDisplays(): Promise<Display[]> {
  const response = await fetch(`${API_BASE}/displays`);
  if (!response.ok) {
    throw new Error(`Failed to fetch displays: ${response.status}`);
  }
  return response.json();
}

export async function fetchDisplay(id: string): Promise<Display> {
  const response = await fetch(`${API_BASE}/displays/${id}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch display: ${response.status}`);
  }
  return response.json();
}

export async function createDisplay(display: Display): Promise<Display> {
  const response = await fetch(`${API_BASE}/displays`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(display),
  });
  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.error || `Failed to create display: ${response.status}`);
  }
  return response.json();
}

export async function updateDisplay(id: string, display: Display): Promise<Display> {
  const response = await fetch(`${API_BASE}/displays/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(display),
  });
  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.error || `Failed to update display: ${response.status}`);
  }
  return response.json();
}

export async function deleteDisplay(id: string): Promise<void> {
  const response = await fetch(`${API_BASE}/displays/${id}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error(`Failed to delete display: ${response.status}`);
  }
}
