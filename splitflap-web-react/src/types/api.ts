/**
 * API Type Definitions
 * 
 * TypeScript interfaces matching the backend DTOs.
 */

export interface DisplayConfig {
  rowCount: number;
  columnCount: number;
}

export interface DisplayContent {
  rows: string[][];
}

export interface Display {
  id: string;
  content: DisplayContent;
  config: DisplayConfig;
}

/**
 * Error types for API communication
 */
export enum ApiErrorType {
  NOT_FOUND = 'NOT_FOUND',
  NETWORK_ERROR = 'NETWORK_ERROR',
  PARSE_ERROR = 'PARSE_ERROR',
  UNKNOWN_ERROR = 'UNKNOWN_ERROR',
}

export interface ApiError {
  type: ApiErrorType;
  message: string;
  statusCode?: number;
  originalError?: Error;
}
