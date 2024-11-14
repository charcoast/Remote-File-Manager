export type Element = {
  name: string;
} & ({ isDir: false } | { isDir: true; subs?: Element[] });
