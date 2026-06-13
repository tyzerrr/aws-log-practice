export function chunk<T>(items: T[], chunkSize: number): T[][] {
  const results: T[][] = [];

  for (let i = 0; i < items.length; i += chunkSize) {
    results.push(items.slice(i, i + chunkSize));
  }

  return results;
}
