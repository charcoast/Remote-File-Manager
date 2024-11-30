import { AxiosResponse } from "axios";
import { CommandRepository } from "../repository/CommandRepository";
import _, { cloneDeep } from "lodash";
import { Element } from "../model/Element";
import { fileStructure } from "../mock/mock.tsx";
import { CreateDirRequest, CreateFileRequest } from "../model/Command.ts";

const pathRegex = RegExp("/{2,}", "g");

class ApiService {
  private static instance: ApiService;

  private constructor() {}

  public static getInstance(): ApiService {
    if (!ApiService.instance) {
      ApiService.instance = new ApiService();
    }
    return ApiService.instance;
  }

  public async readFile(path: string, filename: string) {
    path = path.replace(pathRegex, "/");
    const response: AxiosResponse<Blob> = await CommandRepository.post(
      {
        command: "read",
        arguments: { path, name: filename },
      },
      { responseType: "blob" }
    );
    return URL.createObjectURL(response.data);
  }

  public async load(element: Element, path: string, force?: boolean) {
    path = path.replace(pathRegex, "/");

    const splittedPath = path.split("/");

    let obj = cloneDeep(element);
    let holder: Element | undefined = obj;
    for (let p of splittedPath) {
      if (p !== "" && holder && holder.isDir) {
        holder = holder.subs?.find((x) => x.name === p);
      }
    }

    if (holder?.isDir && holder.subs && holder.subs.length > 0 && !force) {
      return fileStructure;
    }

    return CommandRepository.postAndGetData<Element[]>({
      command: "ls",
      arguments: { path },
    })
      .then((e) => {
        if (holder?.isDir) {
          holder.subs = e;
        }
        return element;
      })
      .catch(() => fileStructure);
  }

  public async createFile(
    path: string,
    filename: string,
    content: string
  ): Promise<void> {
    return await CommandRepository.postAndGetData({
      command: "mkfile",
      arguments: {
        path,
        name: filename,
        content,
      } satisfies CreateFileRequest,
    });
  }

  public async createDirectory(path: string): Promise<void> {
    path = path.replace(pathRegex, "/");
    return await CommandRepository.postAndGetData({
      command: "mkdir",
      arguments: {
        path,
      } satisfies CreateDirRequest,
    });
  }

  public async deleteDirectory(path: string) {
    path = path.replace(pathRegex, "/")
    return ;
  }
}

export default ApiService.getInstance();
