import {ReactNode, useCallback, useEffect, useState} from 'react'
import {Element} from "./model/Element.ts"
import './App.css'
import {SimpleTreeView, TreeItem} from "@mui/x-tree-view";
import {CommandRepository} from "./repository/CommandRepository.ts";




function App() {

    const [elements, setElements] = useState<Element[]>();
    const [path, setPath] = useState<string>("./")

    useEffect(() => {
        CommandRepository.post<Element[]>({command: "ls", arguments: {path}})
            .then((elements) => {
                setElements(elements)
            })
    }, [path]);

    const getTree = useCallback((elements: Element[]): ReactNode => {
        return elements?.map((x) => {
            return <TreeItem itemId={x.name} label={x.name}>{x.isDir && getTree(x.subs ?? [])}</TreeItem>
        })
    }, [path])

  return (
    <SimpleTreeView>
        {getTree(elements ?? [])}
    </SimpleTreeView>
  )
}

export default App
