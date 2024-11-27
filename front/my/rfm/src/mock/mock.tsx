import { Element } from "../model/Element.ts";

export const fileStructure: Element = {
  name: "/",
  isDir: true,
  subs: [
    {
      name: "home",
      isDir: true,
      subs: [
        {
          name: "user",
          isDir: true,
          subs: [
            {
              name: "documents",
              isDir: true,
              subs: [
                {
                  name: "file1.txt",
                  isDir: false,
                },
                {
                  name: "file2.pdf",
                  isDir: false,
                },
              ],
            },
            {
              name: "photos",
              isDir: true,
              subs: [
                {
                  name: "photo1.jpg",
                  isDir: false,
                },
                {
                  name: "photo2.png",
                  isDir: false,
                },
              ],
            },
          ],
        },
      ],
    },
    {
      name: "var",
      isDir: true,
      subs: [
        {
          name: "log",
          isDir: true,
          subs: [
            {
              name: "syslog",
              isDir: false,
            },
            {
              name: "dmesg",
              isDir: false,
            },
          ],
        },
        {
          name: "www",
          isDir: true,
          subs: [
            {
              name: "index.html",
              isDir: false,
            },
            {
              name: "style.css",
              isDir: false,
            },
          ],
        },
      ],
    },
  ],
};
