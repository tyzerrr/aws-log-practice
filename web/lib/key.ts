type KeySeed = unknown;

function stableStringify(value: KeySeed): string {
  if (value === null || typeof value !== "object") {
    return String(value);
  }

  if (value instanceof URL) {
    return value.href;
  }

  if (Array.isArray(value)) {
    return `[${value.map(stableStringify).join(",")}]`;
  }

  return `{${Object.entries(value)
    .sort(([left], [right]) => left.localeCompare(right))
    .map(([key, entryValue]) => `${key}:${stableStringify(entryValue)}`)
    .join(",")}}`;
}

export function hashKey(...seeds: KeySeed[]): string {
  const input = seeds.map(stableStringify).join("|");
  let hash = 0x811c9dc5;

  for (let i = 0; i < input.length; i += 1) {
    hash ^= input.charCodeAt(i);
    hash = Math.imul(hash, 0x01000193);
  }

  return hash.toString(36);
}
