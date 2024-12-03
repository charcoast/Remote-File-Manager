import { Box, IconButton, SvgIcon, Typography } from "@mui/material";
import {
  CreateNewFolderOutlined,
  DeleteOutlined,
  NoteAddOutlined,
} from "@mui/icons-material";

import { Element } from "../model/Element";
import React from "react";

interface FileTreeItemProps {
  element: Element;
  isExpanded: boolean;
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
  setDeleting: (
    params: { path: string }
      & (
        | { type: "dir" }
        | {
          type: "file";
          filename: string;
        }
        | undefined
      )
  ) => void;
}

export const FileTreeItem: React.FC<FileTreeItemProps> = ({
  element,
  isExpanded,
  localPath,
  setAdding,
  setDeleting
}) => {
  return (
    <Box display="flex" alignItems="center" justifyContent="space-between">
      <Box display="flex" alignItems="center" justifyContent="flex-start">
        <SvgIcon viewBox="0 0 16 16" sx={{ fontSize: "1.3rem", marginRight: "0.5rem" }}>
          {element.isDir && isExpanded &&
            <svg aria-hidden="true" focusable="false" viewBox="0 0 16 16" width="16" height="16" fill="#54aeff" >
              <path d="M.513 1.513A1.75 1.75 0 0 1 1.75 1h3.5c.55 0 1.07.26 1.4.7l.9 1.2a.25.25 0 0 0 .2.1H13a1 1 0 0 1 1 1v.5H2.75a.75.75 0 0 0 0 1.5h11.978a1 1 0 0 1 .994 1.117L15 13.25A1.75 1.75 0 0 1 13.25 15H1.75A1.75 1.75 0 0 1 0 13.25V2.75c0-.464.184-.91.513-1.237Z">
              </path>
            </svg>
          }
          {element.isDir && !isExpanded &&
            <svg aria-hidden="true" focusable="false" viewBox="0 0 16 16" width="16" height="16" fill="#54aeff" >
              <path d="M1.75 1A1.75 1.75 0 0 0 0 2.75v10.5C0 14.216.784 15 1.75 15h12.5A1.75 1.75 0 0 0 16 13.25v-8.5A1.75 1.75 0 0 0 14.25 3H7.5a.25.25 0 0 1-.2-.1l-.9-1.2C6.07 1.26 5.55 1 5 1H1.75Z">
              </path>
            </svg>
          }
          {!element.isDir &&
            <svg aria-hidden="true" focusable="false" viewBox="0 0 16 16" width="16" height="16" fill="currentColor" >
              <path d="M2 1.75C2 .784 2.784 0 3.75 0h6.586c.464 0 .909.184 1.237.513l2.914 2.914c.329.328.513.773.513 1.237v9.586A1.75 1.75 0 0 1 13.25 16h-9.5A1.75 1.75 0 0 1 2 14.25Zm1.75-.25a.25.25 0 0 0-.25.25v12.5c0 .138.112.25.25.25h9.5a.25.25 0 0 0 .25-.25V6h-2.75A1.75 1.75 0 0 1 9 4.25V1.5Zm6.75.062V4.25c0 .138.112.25.25.25h2.688l-.011-.013-2.914-2.914-.013-.011Z">
              </path>
            </svg>
          }
        </SvgIcon>
        <Typography>{element.name}</Typography>
      </Box>
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
                <CreateNewFolderOutlined sx={{ fontSize: "1.25rem" }} />
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
                <NoteAddOutlined sx={{ fontSize: "1.25rem" }} />
              </IconButton>
            </>
          )}
          <IconButton
            onClick={(e) => {
              console.log("chamou o delete");

              if (element.isDir) {
                setDeleting({
                  path: localPath,
                  type: "dir"
                })
              } else {
                setDeleting({
                  path: localPath,
                  type: "file",
                  filename: element.name
                })
              }
              e.stopPropagation();
            }}
          >
            <DeleteOutlined sx={{ fontSize: "1.25rem" }} />
          </IconButton>
        </Box>
      </Box>
    </Box>
  );
};
