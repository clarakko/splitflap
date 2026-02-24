/**
 * Character set and flip animation logic
 */

// Standard character set for split-flap displays (MVP: single characters)
// Includes: space, uppercase, lowercase, digits, punctuation
const CHARACTER_SET = ' ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.,:!-?'.split('');


/**
 * Calculate the sequence of characters to display during a flip animation
 * from the current character to the target character.
 *
 * @param current - The current character being displayed
 * @param target - The target character to flip to
 * @returns Array of characters to display in sequence
 */
export function calculateFlipPath(current: string, target: string): string[] {
  const currentIdx = CHARACTER_SET.indexOf(current);
  const targetIdx = CHARACTER_SET.indexOf(target);

  // If character not in set or already at target, return just the target
  if (currentIdx === -1 || targetIdx === -1 || currentIdx === targetIdx) {
    return [target];
  }

  // Find the shorter direction: forward or backward (wrapping)
  let path: string[] = [];

  // Forward distance
  const forwardDist =
    targetIdx >= currentIdx
      ? targetIdx - currentIdx
      : CHARACTER_SET.length - currentIdx + targetIdx;

  // Backward distance
  const backwardDist = CHARACTER_SET.length - forwardDist;

  // Choose the shorter path
  const goForward = forwardDist <= backwardDist;

  if (goForward) {
    // Go forward (or wrap around)
    for (let i = 0; i < forwardDist; i++) {
      const idx = (currentIdx + i) % CHARACTER_SET.length;
      path.push(CHARACTER_SET[idx]);
    }
  } else {
    // Go backward (or wrap around)
    for (let i = 0; i < backwardDist; i++) {
      const idx = (currentIdx - i + CHARACTER_SET.length) % CHARACTER_SET.length;
      path.push(CHARACTER_SET[idx]);
    }
  }

  // Ensure target is at the end
  if (path[path.length - 1] !== target) {
    path.push(target);
  }

  return path;
}

/**
 * Get the character set used for flip animations
 */
export function getCharacterSet(): string[] {
  return [...CHARACTER_SET];
}
