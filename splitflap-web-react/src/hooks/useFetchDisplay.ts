/**
 * Custom hook for fetching display data from the API
 */

import { useState, useEffect } from 'react';
import type { Display, ApiError } from '../types/api';
import { ApiErrorType } from '../types/api';
import { BASE_URL } from '../config/api';

interface UseFetchDisplayResult {
  data: Display | null;
  loading: boolean;
  error: ApiError | null;
}

export function useFetchDisplay(displayId: string): UseFetchDisplayResult {
  const [data, setData] = useState<Display | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<ApiError | null>(null);

  useEffect(() => {
    let isMounted = true;

    const fetchDisplay = async () => {
      try {
        setLoading(true);
        setError(null);

        const response = await fetch(`${BASE_URL}/api/v1/displays/${displayId}`);

        if (!response.ok) {
          if (response.status === 404) {
            throw {
              type: ApiErrorType.NOT_FOUND,
              message: `Display with ID "${displayId}" not found`,
              statusCode: 404,
            } as ApiError;
          }
          
          throw {
            type: ApiErrorType.UNKNOWN_ERROR,
            message: `Failed to fetch display: ${response.statusText}`,
            statusCode: response.status,
          } as ApiError;
        }

        const displayData: Display = await response.json();
        
        if (isMounted) {
          setData(displayData);
        }
      } catch (err) {
        if (isMounted) {
          // If it's already an ApiError, use it directly
          if (err && typeof err === 'object' && 'type' in err) {
            setError(err as ApiError);
          } else if (err instanceof TypeError) {
            // Network errors typically throw TypeError
            setError({
              type: ApiErrorType.NETWORK_ERROR,
              message: 'Network error: Unable to connect to the API',
              originalError: err,
            });
          } else if (err instanceof SyntaxError) {
            // JSON parsing errors
            setError({
              type: ApiErrorType.PARSE_ERROR,
              message: 'Failed to parse response data',
              originalError: err,
            });
          } else {
            // Generic fallback
            setError({
              type: ApiErrorType.UNKNOWN_ERROR,
              message: err instanceof Error ? err.message : 'An unexpected error occurred',
              originalError: err instanceof Error ? err : undefined,
            });
          }
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    fetchDisplay();

    // Cleanup function to prevent state updates on unmounted component
    return () => {
      isMounted = false;
    };
  }, [displayId]);

  return { data, loading, error };
}
