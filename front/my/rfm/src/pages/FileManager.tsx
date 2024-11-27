import { ReactNode, useCallback, useEffect, useState } from "react";
import { Element } from "../model/Element.ts";
import { SimpleTreeView, TreeItem } from "@mui/x-tree-view";
import { CommandRepository } from "../repository/CommandRepository.ts";
import _, { cloneDeep } from "lodash";
import { Box, Grid, IconButton, Paper, styled, TextField } from "@mui/material";
import { Check, Close } from "@mui/icons-material";
import { CreateDirRequest, CreateFileRequest } from "../model/Command.ts";
import { AxiosResponse } from "axios";
import { fileStructure } from "../mock/mock.tsx";
import backgroundImage from "../assets/background.png";
import { FileTreeItem } from "../components/FileTreeItem.tsx";

const pathRegex = RegExp("/{2,}", "g");

function FileManager() {
  const [element, setElement] = useState<Element>({
    name: "/",
    isDir: true,
    subs: undefined,
  });

  const [adding, setAdding] = useState<
    | ({ name: string; path: string } & (
        | { type: "dir" }
        | {
            type: "file";
            content: string;
          }
      ))
    | undefined
  >();

  const addItem = (
    item: { name: string; path: string } & (
      | { type: "dir" }
      | {
          type: "file";
          content: string;
        }
      | undefined
    )
  ) => {
    setAdding(item);
  };

  const [display, setDisplay] = useState<string | undefined>();

  useEffect(() => {
    console.log(adding);
  }, [adding]);

  const readFile = useCallback(async (path: string, filename: string) => {
    path = path.replace(pathRegex, "/");
    const response: AxiosResponse<Blob> = await CommandRepository.post(
      {
        command: "read",
        arguments: { path, name: filename },
      },
      { responseType: "blob" }
    );
    const href = URL.createObjectURL(response.data);
    setDisplay(href);
  }, []);

  const load = useCallback(
    (path: string, force?: boolean) => {
      path = path.replace(pathRegex, "/");

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

      CommandRepository.postAndGetData<Element[]>({
        command: "ls",
        arguments: { path },
      })
        .then((e) => {
          if (holder?.isDir) {
            holder.subs = e;
          }
          setElement(obj);
        })
        .catch(() => setElement(fileStructure));
    },
    [element]
  );

  const createFile = useCallback(
    async (path: string, filename: string, content: string): Promise<void> => {
      await CommandRepository.postAndGetData({
        command: "mkfile",
        arguments: {
          path,
          name: filename,
          content,
        } satisfies CreateFileRequest,
      });
    },
    []
  );

  const createDirectory = useCallback(async (path: string): Promise<void> => {
    path = path.replace(pathRegex, "/");
    await CommandRepository.postAndGetData({
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
      return elements?.map((element) => {
        const localPath = path + element.name;
        return (
          <TreeItem
            itemId={localPath}
            label={
              <FileTreeItem
                element={element}
                localPath={localPath}
                setAdding={addItem}
              />
            }
            onClick={() => {
              if (element.isDir) {
                load(localPath);
              } else {
                readFile(path, element.name);
              }
            }}
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
                            adding.path + "/" + adding!.name
                          );
                        } else {
                          promise = createFile(
                            adding.path,
                            adding.name,
                            adding.content
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
            {element.isDir &&
              (!element.subs ? (
                <TreeItem itemId={Math.random().toString()} label={""} />
              ) : (
                getTree(localPath, element.subs ?? [])
              ))}
          </TreeItem>
        );
      });
    },
    [element, adding]
  );

  const GradientBox = styled(Box)({
    height: "100vh",
    padding: "2rem",
    color: "#27272a",
    backgroundImage: `url(${backgroundImage})`,
  });

  const CircleButton = styled(IconButton)({
    width: "1rem",
    height: "1rem",
    borderRadius: "50%",
    backgroundColor: "#d1d5db",
    "&:hover": {
      backgroundColor: "#e0e0e0",
    },
  });

  return (
    <GradientBox>
      <Paper
        elevation={1}
        sx={{
          height: "100%",
          borderRadius: "1rem",
          overflow: "hidden",
          border: "1px solid rgba(0, 0, 0, 0.2)",
        }}
      >
        <Box sx={{ height: "100%", display: "flex", flexDirection: "row" }}>
          <Box
            sx={{
              height: "100%",
              width: "30%",
              backgroundColor: "#f9fafb",
              borderRight: "1px solid #e5e7eb",
              padding: "1rem",
            }}
          >
            <Grid container spacing={1} sx={{ marginBottom: "5%" }}>
              <Grid item>
                <CircleButton
                  sx={{
                    backgroundColor: "#FF464F",
                    "&:hover": { backgroundColor: "#FF464F" },
                  }}
                />
              </Grid>
              <Grid item>
                <CircleButton
                  sx={{
                    backgroundColor: "#FFB41B",
                    "&:hover": { backgroundColor: "#FFB41B" },
                  }}
                />
              </Grid>
              <Grid item>
                <CircleButton
                  sx={{
                    backgroundColor: "#20CA2D",
                    "&:hover": { backgroundColor: "#20CA2D" },
                  }}
                />
              </Grid>
            </Grid>
            <SimpleTreeView
              sx={{
                maxHeight: "95%",
                overflow: "auto",
              }}
            >
              {getTree("/", element.isDir ? (element.subs ?? []) : [])}
            </SimpleTreeView>
          </Box>
          <Box component="main" sx={{ padding: "1rem", flexGrow: 1 }}>
            {display && (
              <iframe
                src={display}
                seamless={true}
                style={{ overflow: "auto", height: "90%", width: "70%" }}
              ></iframe>
            )}
          </Box>
        </Box>
      </Paper>
    </GradientBox>
  );
}

export default FileManager;
