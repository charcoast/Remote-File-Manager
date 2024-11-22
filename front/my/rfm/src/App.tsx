import { ReactNode, useCallback, useEffect, useState } from "react";
import { Element } from "./model/Element.ts";
import { SimpleTreeView, TreeItem } from "@mui/x-tree-view";
import { CommandRepository } from "./repository/CommandRepository.ts";
import _, { cloneDeep } from "lodash";
import { Box, IconButton, Paper, TextField } from "@mui/material";
import {
  Check,
  Close,
  CreateNewFolder,
  Delete,
  NoteAdd,
} from "@mui/icons-material";
import { CreateDirRequest, CreateFileRequest } from "./model/Command.ts";

function App() {
  const [element, setElement] = useState<Element>({
    name: "/",
    isDir: true,
    subs: undefined,
  });

  const [adding, setAdding] = useState<
    | ({ name: string; path: string } & (
        | { type: "dir" }
        | { type: "file"; content: string }
      ))
    | undefined
  >();

  useEffect(() => {
    console.log(adding);
  }, [adding]);

  const load = useCallback(
    (path: string, force?: boolean) => {
      path = path.replace("//", "/");

      const splittedPath = path.split("/");

      let obj = cloneDeep(element);
      let holder: Element | undefined = obj;
      for (let p of splittedPath) {
        if (p === "") continue;
        if (holder && holder.isDir) {
          holder = holder.subs?.find((x) => x.name === p);
        }
      }

      if (holder?.isDir && holder.subs && holder.subs.length > 0 && !force) {
        return;
      }

      CommandRepository.post<Element[]>({
        command: "ls",
        arguments: { path },
      }).then((e) => {
        if (holder?.isDir) {
          holder.subs = e;
        }
        setElement(obj);
      });
    },
    [element],
  );

  const createFile = useCallback(
    async (path: string, filename: string, content: string): Promise<void> => {
      await CommandRepository.post({
        command: "mkfile",
        arguments: {
          path,
          name: filename,
          content,
        } satisfies CreateFileRequest,
      });
    },
    [],
  );

  const createDirectory = useCallback(async (path: string): Promise<void> => {
    path = path.replace(RegExp("/{2,}", "g"), "/");
    await CommandRepository.post({
      command: "mkdir",
      arguments: {
        path,
      } satisfies CreateDirRequest,
    });
  }, []);

  useEffect(() => {
    load("/");
  }, []);

  const getTree = useCallback(
    (path: string, elements: Element[]): ReactNode => {
      path = `${path}/`;
      return elements?.map((x) => {
        const localPath = path + x.name;
        return (
          <TreeItem
            itemId={path + x.name}
            label={
              <Box
                display="flex"
                alignItems="center"
                justifyContent="space-between"
              >
                {x.name}

                <Box display="flex" flexGrow={1}>
                  <Box flexGrow={1}>
                    {x.isDir && (
                      <>
                        <IconButton
                          onClick={(e) => {
                            setAdding({
                              type: "dir",
                              name: "",
                              path: localPath,
                            });
                            e.stopPropagation();
                          }}
                        >
                          <CreateNewFolder />
                        </IconButton>
                        <IconButton
                          onClick={(e) => {
                            setAdding({
                              type: "file",
                              name: "",
                              path: localPath,
                              content: "",
                            });
                            e.stopPropagation();
                          }}
                        >
                          <NoteAdd />
                        </IconButton>
                      </>
                    )}
                  </Box>

                  <Box>
                    <IconButton
                      onClick={(e) => {
                        e.stopPropagation();
                      }}
                    >
                      <Delete />
                    </IconButton>
                  </Box>
                </Box>
              </Box>
            }
            onClick={() => load(localPath)}
          >
            {adding && adding.path === localPath && (
              <TreeItem
                itemId={"adding"}
                label={
                  <Box display="flex" alignItems="center">
                    <Box>
                      <TextField
                        size="small"
                        value={adding.name}
                        onClick={(e) => e.stopPropagation()}
                        onKeyDown={(e) => {
                          e.stopPropagation();
                        }}
                        onChange={(e) => {
                          setAdding((p) => ({
                            ...p!,
                            name: e.target.value.replace(RegExp("[/.]"), ""),
                          }));
                          e.stopPropagation();
                        }}
                      />
                      {adding.type === "file" && (
                        <input
                          type="file"
                          onClick={(e) => e.stopPropagation()}
                          onChange={async (e) => {
                            const files = e.target.files;
                            if (files && files.length > 0) {
                              const fileReader = new FileReader();
                              fileReader.onload = () => {
                                const base64 = (
                                  fileReader.result?.toString() ?? ""
                                )
                                  .replace("data:", "")
                                  .replace(/^.+,/, "");
                                setAdding((p) => ({
                                  ...p!,
                                  content: base64,
                                }));
                              };
                              fileReader.readAsDataURL(files[0]);
                            }
                            e.stopPropagation();
                          }}
                        />
                      )}
                    </Box>
                    <IconButton
                      onClick={(e) => {
                        e.stopPropagation();
                        setAdding(undefined);
                      }}
                    >
                      <Close />
                    </IconButton>
                    <IconButton
                      onClick={(e) => {
                        e.stopPropagation();
                        if (!adding) return;
                        let promise: Promise<void>;
                        if (adding.type === "dir") {
                          promise = createDirectory(
                            adding.path + "/" + adding!.name,
                          );
                        } else {
                          promise = createFile(
                            adding.path,
                            adding.name,
                            adding.content,
                          );
                        }
                        promise.then((_) => {
                          setAdding(undefined);
                          load(localPath, true);
                        });
                      }}
                    >
                      <Check />
                    </IconButton>
                  </Box>
                }
              ></TreeItem>
            )}
            {x.isDir &&
              (!x.subs ? (
                <TreeItem itemId={Math.random().toString()} label={""} />
              ) : (
                getTree(localPath, x.subs ?? [])
              ))}
          </TreeItem>
        );
      });
    },
    [element, adding],
  );

  return (
    <Paper sx={{ width: "30vw", height: "100vh" }}>
      <SimpleTreeView>
        {getTree("/", element.isDir ? (element.subs ?? []) : [])}
      </SimpleTreeView>
    </Paper>
  );
}

export default App;
