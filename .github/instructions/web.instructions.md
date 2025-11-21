# React/TypeScript Web Coding Standards

## General Principles

- Functional components with hooks (no class components)
- TypeScript for type safety
- Keep components small and focused
- Colocate related files (component + styles + tests)
- Follow React best practices

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
│   └── useFetchDisplay.ts
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
  onFlipComplete 
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

## React Patterns

### Functional Components

```typescript
// ✅ Good: Arrow function for named export
export const SplitFlapCell = ({ targetChar }: SplitFlapCellProps) => {
  const [currentChar, setCurrentChar] = useState(' ');
  
  // Implementation
  
  return <div className={styles.cell}>{currentChar}</div>;
};
```

### Custom Hooks

```typescript
// ✅ Good: Extract reusable logic
export function useFlipAnimation(targetChar: string, duration: number) {
  const [currentChar, setCurrentChar] = useState(' ');
  const [isFlipping, setIsFlipping] = useState(false);
  
  useEffect(() => {
    // Animation logic
  }, [targetChar, duration]);
  
  return { currentChar, isFlipping };
}

// Usage in component
const { currentChar, isFlipping } = useFlipAnimation(targetChar, 100);
```

### Data Fetching

```typescript
// ✅ Good: Custom hook for API calls
export function useFetchDisplay(displayId: string) {
  const [data, setData] = useState<Display | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  
  useEffect(() => {
    fetch(`/api/v1/displays/${displayId}`)
      .then(res => res.json())
      .then(setData)
      .catch(setError)
      .finally(() => setLoading(false));
  }, [displayId]);
  
  return { data, loading, error };
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

@keyframes flip {
  0% { transform: rotateX(0deg); }
  50% { transform: rotateX(90deg); }
  100% { transform: rotateX(0deg); }
}
```

## Animation Logic

### Character Set

```typescript
// Phase 1: Alphanumeric only
const CHAR_SET = ' ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';

function getCharIndex(char: string): number {
  return CHAR_SET.indexOf(char.toUpperCase());
}

function calculateFlipPath(from: string, to: string): string[] {
  const fromIdx = getCharIndex(from);
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
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
});
```

### Fetch Wrapper

```typescript
// utils/api.ts
const API_BASE = '/api/v1';

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
import { render, screen, waitFor } from '@testing-library/react';
import { SplitFlapCell } from './SplitFlapCell';

describe('SplitFlapCell', () => {
  it('renders initial character', () => {
    render(<SplitFlapCell targetChar="A" />);
    expect(screen.getByText('A')).toBeInTheDocument();
  });
  
  it('flips to target character', async () => {
    const { rerender } = render(<SplitFlapCell targetChar="A" />);
    
    rerender(<SplitFlapCell targetChar="C" />);
    
    await waitFor(() => {
      expect(screen.getByText('C')).toBeInTheDocument();
    });
  });
});
```

## Performance

### Memoization

```typescript
// ✅ Good: Memoize expensive calculations
const flipPath = useMemo(
  () => calculateFlipPath(currentChar, targetChar),
  [currentChar, targetChar]
);

// ✅ Good: Memoize callbacks
const handleFlipComplete = useCallback(() => {
  console.log('Flip complete');
}, []);
```

### Avoid Unnecessary Renders

```typescript
// ✅ Good: memo for pure components
export const SplitFlapCell = memo(({ targetChar }: SplitFlapCellProps) => {
  // Component logic
});
```

## Dependencies

### Phase 1 Only Use

- `react`
- `react-dom`
- Vite defaults (no additional libraries)

### Don't Add Yet

- ❌ State management (Redux/Zustand) - Phase 2+
- ❌ Routing (React Router) - Phase 3+
- ❌ UI libraries (MUI/Chakra) - Build custom components
- ❌ Animation libraries (Framer Motion) - Use CSS animations

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
- ✅ Fetch from hardcoded `/api/v1/displays/demo`
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
