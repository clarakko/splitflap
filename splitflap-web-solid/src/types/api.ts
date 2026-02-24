export interface DisplayContent {
  rows: string[][];
}

export interface DisplayConfig {
  rowCount: number;
  columnCount: number;
}

export interface Display {
  id: string;
  content: DisplayContent;
  config: DisplayConfig;
}
