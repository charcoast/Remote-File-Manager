import { ReactNode, useCallback, useEffect, useState } from "react";
import { Element } from "../model/Element.ts";
import { SimpleTreeView, TreeItem, useTreeViewApiRef } from "@mui/x-tree-view";
import { Box, Grid, IconButton, Paper, styled, TextField } from "@mui/material";
import { Check, Close } from "@mui/icons-material";
import backgroundImage from "../assets/background.png";
import { FileTreeItem } from "../components/FileTreeItem.tsx";
import ApiService from "../services/ApiService.ts";
import { fileStructure } from "../mock/mock.tsx";
import { MacButtons } from "../components/MacButtons.tsx";

function FileManager() {

  const [expanded, setExpanded] = useState<string[]>([]);

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

  const apiRef = useTreeViewApiRef();

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

  const handleToggle = (event: React.ChangeEvent<{}>, itemId: string, isExpanded: boolean) => {
    let state: string[] = [...expanded]
    state.push(itemId)
    setExpanded(state);
  };

  useEffect(() => {
    console.log(adding);
  }, [adding]);

  const readFile = useCallback((path: string, filename: string) => {
    console.log("reading file");
    ApiService.readFile(path, filename).then((href) => setDisplay(href));
  }, []);

  const load = useCallback((path: string, force?: boolean) => {
    console.log("loading.path", path);
    ApiService.load(element, path, force).then(response => {
      if (response != fileStructure) {
        console.log("response", response);
        setElement(response)
      }
    }).catch(() => setElement(fileStructure))
  }, [element])

  const createFile = useCallback(
    (path: string, filename: string, content: string) => {
      console.log("creating file");
      return ApiService.createFile(path, filename, content);
    }, []
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
    console.log("element", element);
    load("/");
  }, []);

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
            onClick={(event) => {
              event.stopPropagation();
              if (element.isDir) {
                load(localPath)
              } else {
                readFile(path, element.name)
              }
            }
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

  const [scale, setScale] = useState<number>(1);

  const adjustIframeScale = useCallback(() => {
    const iframe = document.querySelector<HTMLIFrameElement>("iframe");
    if (!iframe || !iframe.contentWindow) return;

    const iframeDocument = iframe.contentWindow.document;
    const contentWidth = iframeDocument.body.scrollWidth;
    const contentHeight = iframeDocument.body.scrollHeight;

    const container = iframe.parentElement; // Get the parent container (main box)
    if (!container) return;

    const containerWidth = container.offsetWidth;
    const containerHeight = container.offsetHeight;

    // Calculate scaling factor based on container vs content sizes
    const scaleWidth = containerWidth / contentWidth;
    const scaleHeight = containerHeight / contentHeight;

    // Use the smaller of the two to fit content within the container
    const newScale = Math.min(scaleWidth, scaleHeight);

    // Apply scaling and set state
    setScale(newScale);
  }, []);

  useEffect(() => {
    // Adjust scaling when the iframe is loaded or resized
    const iframe = document.querySelector("iframe");
    if (iframe) {
      iframe.onload = adjustIframeScale; // Run scaling adjustment after iframe loads
    }

    // Adjust scaling when the window is resized
    window.addEventListener("resize", adjustIframeScale);

    return () => {
      window.removeEventListener("resize", adjustIframeScale);
    };
  }, [adjustIframeScale]);

  const GradientBox = styled(Box)({
    height: "100vh",
    padding: "2rem",
    color: "#27272a",
    backgroundImage: `url(${backgroundImage})`,
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
            <MacButtons />
            <SimpleTreeView
              apiRef={apiRef}
              expandedItems={expanded}
              onItemExpansionToggle={handleToggle}
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
        seamless
        style={{
          overflow: "auto",
          border: "none",
          transform: `scale(${scale})`, // Apply dynamic scaling
          transformOrigin: "0 0", // Set scaling origin to top-left
          width: `${100 / scale}%`, // Compensate for scaling
          height: `${100 / scale}%`, // Compensate for scaling
        }}
      ></iframe>
            )}
          </Box>
        </Box>
      </Paper>
    </GradientBox>
  );
}

export default FileManager;
