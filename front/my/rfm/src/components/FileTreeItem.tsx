import { Box, IconButton, Typography } from "@mui/material";
import {
  CreateNewFolderOutlined,
  DeleteForeverOutlined,
  NoteAddOutlined,
} from "@mui/icons-material";

import { Element } from "../model/Element";
import React from "react";

interface FileTreeItemProps {
  element: Element;
  localPath: string;
  setAdding: (
    item: { name: string; path: string } & (
      | { type: "dir" }
      | {
          type: "file";
          content: string;
        }
      | undefined
    )
  ) => void;
}

export const FileTreeItem: React.FC<FileTreeItemProps> = ({
  element,
  localPath,
  setAdding,
}) => {
  return (
    <Box display="flex" alignItems="center" justifyContent="space-between">
      <Typography>{element.name}</Typography>
      <Box>
        <Box flexGrow={1}>
          {element.isDir && (
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
                <CreateNewFolderOutlined />
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
                <NoteAddOutlined />
              </IconButton>
            </>
          )}
          <IconButton
            onClick={(e) => {
              e.stopPropagation();
            }}
          >
            <DeleteForeverOutlined />
          </IconButton>
        </Box>
      </Box>
    </Box>
  );
};
