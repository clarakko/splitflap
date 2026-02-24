# SolidJS/TypeScript Web Coding Standards

## General Principles

- Functional components with signals and effects
- TypeScript for type safety
- Keep components small and focused
- Colocate related files (component + styles + tests)
- Follow SolidJS best practices

## Code Style

### Naming Conventions

- Components: PascalCase (`SplitFlapCell.tsx`, `DisplayPreview.tsx`)
- Hooks: camelCase with `use` prefix (`useFlipAnimation`, `useFetchDisplay`)
- Types/Interfaces: PascalCase (`Display`, `CellState`)
- Constants: UPPER_SNAKE_CASE (`FLIP_DURATION_MS`, `API_BASE_URL`)
- Variables/functions: camelCase (`currentChar`, `calculateFlips`)

### File Organization

```
src/
├── components/
│   ├── SplitFlapCell/
│   │   ├── SplitFlapCell.tsx        # Component
│   │   ├── SplitFlapCell.module.css # Styles
│   │   └── SplitFlapCell.test.tsx   # Tests
│   └── DisplayPreview/
│       ├── DisplayPreview.tsx
│       └── DisplayPreview.module.css
├── hooks/
│   └── createFetchDisplay.ts        # Resource factory
├── types/
│   └── api.ts                       # API response types
└── utils/
    └── flipCalculator.ts            # Animation logic
```

## TypeScript Patterns

### Component Props

```typescript
// ✅ Good: Explicit interface
interface SplitFlapCellProps {
  targetChar: string;
  flipDuration?: number; // Optional with default
  onFlipComplete?: () => void;
}

export function SplitFlapCell({
  targetChar,
  flipDuration = 100,
  onFlipComplete,
}: SplitFlapCellProps) {
  // Implementation
}
```

### API Types

```typescript
// Match backend DTOs exactly
export interface Display {
  id: string;
  content: DisplayContent;
  config: DisplayConfig;
}

export interface DisplayContent {
  rows: string[][];
}

export interface DisplayConfig {
  rowCount: number;
  columnCount: number;
}
```

### State Types

```typescript
interface CellState {
  current: string;
  target: string;
  isFlipping: boolean;
  progress: number; // 0.0 to 1.0
}
```

## SolidJS Patterns

### Functional Components

```typescript
// ✅ Good: SolidJS component function
export const SplitFlapCell = (props: SplitFlapCellProps) => {
  const [currentChar, setCurrentChar] = createSignal(' ');

  // Implementation

  return <div class={styles.cell}>{currentChar()}</div>;
};
```

### Signals and Effects

```typescript
// ✅ Good: Use signals for reactive state
import { createSignal, createEffect } from 'solid-js';

export function SplitFlapCell(props: SplitFlapCellProps) {
  const [currentChar, setCurrentChar] = createSignal(' ');
  const [isFlipping, setIsFlipping] = createSignal(false);

  // Track target changes
  createEffect(() => {
    const target = props.targetChar();
    if (currentChar() === target) return;

    const flipPath = calculateFlipPath(currentChar(), target);
    let step = 0;

    setIsFlipping(true);
    const interval = setInterval(() => {
      if (step >= flipPath.length) {
        clearInterval(interval);
        setIsFlipping(false);
        return;
      }

      setCurrentChar(flipPath[step]);
      step++;
    }, props.flipDuration ?? 100);
  });

  return (
    <div class={styles.cell} data-flipping={isFlipping()}>
      {currentChar()}
    </div>
  );
}
```

### Resource API for Data Fetching

```typescript
// ✅ Good: Use createResource for async operations
import { createResource } from 'solid-js';

export function createFetchDisplay(displayId: () => string) {
  const [display, { refetch }] = createResource(displayId, async (id) => {
    const response = await fetch(`/api/v1/displays/${id}`);
    props: SplitFlapCellProps) => (
  <div class={styles.cell}>
    <div class={styles.flap}>{props.targetChar}etch display: ${response.status}`);
    }

    return response.json() as Promise<Display>;
  });

  return { display, refetch };
}

// Usage in component
export function DisplayPreview(props: { displayId: () => string }) {
  const { display } = createFetchDisplay(props.displayId);

  return (
    <>
      <Show when={display.loading}>
        <div>Loading...</div>
      </Show>
      <Show when={display.error}>
        <div>Error: {display.error.message}</div>
      </Show>
      <Show when={display()}>
        {(data) => <div>{data().id}</div>}
      </Show>
    </>
  );
}
```

## Styling

### CSS Modules

```typescript
// SplitFlapCell.tsx
import styles from './SplitFlapCell.module.css';

export const SplitFlapCell = () => (
  <div className={styles.cell}>
    <div className={styles.flap}>A</div>
  </div>
);
```

```css
/* SplitFlapCell.module.css */
.cell {
  width: 40px;
  height: 60px;
  position: relative;
  overflow: hidden;
}

.flap {
  position: absolute;
  width: 100%;
  height: 100%;
  background: #000;
  color: #fff;
  animation: flip 0.1s ease-in-out;
}
 with createEffect
import { createEffect } from 'solid-js';

createEffect(() => {
  if (currentChar() === targetChar()) return;

  const flipPath = calculateFlipPath(currentChar(), targetChar());
  let step = 0;
  setIsFlipping(true);

  const interval = setInterval(() => {
    if (step >= flipPath.length) {
      clearInterval(interval);
      setIsFlipping(false);
      return;
    }

    setCurrentChar(flipPath[step]);
    step++;
  }, flipDuration());

  // Cleanup: Clear interval on unmount
  return () => clearInterval(interval);
}
  const toIdx = getCharIndex(to);
  const path: string[] = [];

  // Calculate shortest circular path
  let current = fromIdx;
  while (current !== toIdx) {
    current = (current + 1) % CHAR_SET.length;
    path.push(CHAR_SET[current]);
  }

  return path;
}
```

### Sequential Flipping

```typescript
// ✅ Good: Step through character set
useEffect(() => {
  if (currentChar === targetChar) return;

  const flipPath = calculateFlipPath(currentChar, targetChar);
  let step = 0;

  const interval = setInterval(() => {
    if (step >= flipPath.length) {
      clearInterval(interval);
      setIsFlipping(false);
      return;
    }

    setCurrentChar(flipPath[step]);
    step++;
  }, flipDuration);

  setIsFlipping(true);

  return () => clearInterval(interval);
}, [targetChar, currentChar, flipDuration]);
```

## API Integration

### Vite Proxy Configuration

```typescript
// vite.config.ts
import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";

export default defineConfig({
  plugins: [solidPlugin()],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
```

### Fetch Wrapper

```typescript
// utils/api.ts
const API_BASE = "/api/v1";

export async function fetchDisplay(id: string): Promise<Display> {
  const response = await fetch(`${API_BASE}/displays/${id}`);

  if (!response.ok) {
    throw new Error(`Failed to fetch display: ${response.status}`);
  }

  return response.json();
}
```

## Testing

### Component Tests

```typescript
// SplitFlapCell.test.tsx
import { render, screen, waitFor } from 'solid-testing-library';
import { SplitFlapCell } from './SplitFlapCell';

describe('SplitFlapCell', () => {
  it('renders initial character', () => {
    render(() => <SplitFlapCell targetChar={() => 'A'} flipDuration={100} />);
    expect(screen.getByText('A')).toBeInTheDocument();
  });

  it('flips to target character', async () => {
    const [target, setTarget] = createSignal('A');

    render(() => <SplitFlapCell targetChar={target} flipDuration={50} />);

    setTarget('C');

    await waitFor(() => {
      expect(screen.getByText('C')).toBeInTheDocument();
    });
  });
});
```

it('renders initial character', () => {
render(() => <SplitFlapCell targetChar={() => 'A'} flipDuration={100} />);
expect(screen.getByText('A')).toBeInTheDocument();
});

it('flips to target character', async () => {
const [target, setTarget] = createSignal('A');

    render(() => <SplitFlapCell targetChar={target} flipDuration={50} />);

    setTarget('C'l targetChar="A" />);

    rerender(<SplitFlapCell targetChar="C" />);

    await waitFor(() => {
      expect(screen.getByText('C')).toBeInTheDocument();
    });

});
});

````with Derived Signals

```typescript
// ✅ Good: Derived signals only recalculate when deps change
import { createMemo } from 'solid-js';

const flipPath = createMemo(() =>
  calculateFlipPath(currentChar(), targetChar())
);

// Memoization is automatic in SolidJS - avoid premature optimization
````

### Avoid Unnecessary Reactivity

```typescript
// ✅ Good: Only reactive when needed
export const SplitFlapCell = (props: SplitFlapCellProps) => {
  // SolidJS optimizes renders automatically
  // No need for memo unless component is truly expensive
}`typescript
// ✅ Good: memo for pure components
export const SplitFlapCell = memo(({ targetChar }: SplitFlapCellProps) => {
  // Component logic
});
```

## Dependencies

solid-js`

- `vite-plugin-solid`
- Vite defaults (no additional libraries)

### Don't Add Yet

- ❌ State management (Pinia/Nanostores) - Phase 2+
- ❌ Routing (Solid Router) - Phase 3+
- ❌ UI libraries - Build custom components
- ❌ Animation libraries - Use CSS animations

## Code Comments

```typescript
// ❌ Bad: Obvious
// Set the current character
setCurrentChar(char);

// ✅ Good: Explains mechanical constraint
// Must flip through entire character set (A→B→C)
// to simulate mechanical split-flap behavior
const flipPath = calculateFlipPath(current, target);
```

## Git Commits

Use conventional commits:

```
feat(web): add SplitFlapCell component with flip animation
fix(web): correct character set wraparound logic
test(web): add unit tests for flip path calculation
style(web): improve cell spacing and contrast
refactor(web): extract flip logic to custom hook
docs(web): add animation architecture comments
```

## Phase Discipline

**DO:**

- ✅ Build flip animation for Phase 1
- ✅ Fetch from hardcoded `/api/v1/signals/effects
- ✅ Keep state management simple (useState/useEffect)

**DON'T:**

- ❌ Add builder UI (Phase 2-3)
- ❌ Add routing (Phase 3)
- ❌ Add WebSocket support (Phase 5)
- ❌ Build embed component (Phase 4 - different package)

## Mechanical Simulation Notes

Remember: This simulates a **physical device** with constraints:

1. Characters flip **sequentially** through the full set
2. Cannot "jump" from A to Z without intermediate flips
3. Timing is critical to maintain realistic appearance
4. Multiple cells can animate simultaneously but independently

```typescript
// ✅ Good: Respects mechanical constraints
// A → B → C (3 flips, 300ms total at 100ms/flip)

// ❌ Wrong: Instant transition
// A → C (0 flips, breaks simulation)
```
