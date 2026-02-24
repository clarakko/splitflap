/**
 * API Configuration
 * 
 * Centralized configuration for API communication.
 * BASE_URL can be extended to use environment variables in future phases.
 */

export const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
