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

  const load = useCallback(
    (path: string) => {
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
    [element],
  );

  return (
    <SimpleTreeView>
      {getTree("/", element.isDir ? (element.subs ?? []) : [])}
    </SimpleTreeView>
  );
}

export default App;
