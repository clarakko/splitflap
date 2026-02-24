import { createResource } from 'solid-js';
import type { Display } from '../types/api';

export function createFetchDisplay(displayId: () => string) {
  const [display, { refetch }] = createResource(displayId, async (id) => {
    if (!id) return undefined;

    const response = await fetch(`/api/v1/displays/${id}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch display: ${response.status}`);
    }

    return response.json() as Promise<Display>;
  });

  return { display, refetch };
}
