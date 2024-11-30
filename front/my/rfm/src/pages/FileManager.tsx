import { ReactNode, useCallback, useEffect, useState } from "react";
import { Element } from "../model/Element.ts";
import { SimpleTreeView, TreeItem } from "@mui/x-tree-view";
import { Box, Grid, IconButton, Paper, styled, TextField } from "@mui/material";
import { Check, Close } from "@mui/icons-material";
import backgroundImage from "../assets/background.png";
import { FileTreeItem } from "../components/FileTreeItem.tsx";
import ApiService from "../services/ApiService.ts";

function FileManager() {
  const [element, setElement] = useState<Element>({
    name: "/",
    isDir: true,
    subs: undefined,
  });

  const [display, setDisplay] = useState<string | undefined>();

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

  const readFile = useCallback((path: string, filename: string) => {
    console.log("reading file");
    ApiService.readFile(path, filename).then((href) => setDisplay(href));
  }, []);

  const load = useCallback((path: string, force?: boolean) => {
    console.log("loading.path", path);
    ApiService.load(element, path, force).then((e) => setElement(e));
  }, []);

  const createFile = useCallback(
    (path: string, filename: string, content: string) => {
      console.log("creating file");
      return ApiService.createFile(path, filename, content);
    },
    []
  );

  const createDirectory = useCallback((path: string) => {
    console.log("creating directory");
    return ApiService.createDirectory(path);
  }, []);

  const deleteDirectory = useCallback((path: string) => {
    console.log("deleting directory");
    return ApiService.deleteDirectory(path);
  }, []);

  useEffect(() => {
    load("/");
  }, [load]);

  const getTree = useCallback(
    (path: string, elements: Element[]): ReactNode => {
      path = `${path}/`;
      return elements?.map((element) => {
        const localPath = path + element.name;
        return (
          <TreeItem
            key={localPath}
            itemId={localPath}
            label={
              <FileTreeItem
                element={element}
                localPath={localPath}
                setAdding={addItem}
              />
            }
            onClick={() =>
              element.isDir ? load(localPath) : readFile(path, element.name)
            }
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
                        onKeyDown={(e) => e.stopPropagation()}
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
                        let promise =
                          adding.type === "dir"
                            ? createDirectory(`${adding.path}/${adding!.name}`)
                            : createFile(
                                adding.path,
                                adding.name,
                                adding.content
                              );
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
