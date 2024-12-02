export type Command = {
  command: string;
  arguments: Record<string, string>;
};

export type CreateFileRequest = {
  name: string;
  path: string;
  content: string;
};

export type CreateDirRequest = {
  path: string;
};

export type DeleteDirRequest = {
  path: string;
}

export type DeleteFileRequest = {
  path: string;
  filename: string;
}
