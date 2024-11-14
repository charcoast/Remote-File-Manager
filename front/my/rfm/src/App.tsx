import { ReactNode, useCallback, useEffect, useState } from "react";
import { Element } from "./model/Element.ts";
import "./App.css";
import { SimpleTreeView, TreeItem } from "@mui/x-tree-view";
import { CommandRepository } from "./repository/CommandRepository.ts";
import _, { cloneDeep } from "lodash";

function App() {
  const [element, setElement] = useState<Element>({
    name: "/",
    isDir: true,
    subs: undefined,
  });

  const load = useCallback((path: string) => {
    path = path.replace("//", "/");

    const splittedPath = path.split("/");
    let elementPath = "subs";
    let obj = cloneDeep(element);
    for (let p of splittedPath) {
      if (p === "") continue;

      elementPath += "[" + p + "]";
    }

    CommandRepository.post<Element[]>({
      command: "ls",
      arguments: { path },
    }).then((e) => {
      console.log(splittedPath, elementPath);
      const newValue = cloneDeep(element);
      _.set(newValue, elementPath, e);

      console.log(newValue);
      setElement(newValue);
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
            itemId={x.name}
            label={x.name}
            onClick={() => load(localPath)}
          >
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
    [],
  );

  return (
    <SimpleTreeView>
      {getTree("/", element.isDir ? (element.subs ?? []) : [])}
    </SimpleTreeView>
  );
}

export default App;
