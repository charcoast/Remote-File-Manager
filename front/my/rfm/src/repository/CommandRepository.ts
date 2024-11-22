import { api } from "../api/api.ts";
import { Command } from "../model/Command.ts";
import { AxiosRequestConfig, AxiosResponse } from "axios";

export class CommandRepository {
  static async postAndGetData<T>(command: Command): Promise<T> {
    return (await this.post<T>(command)).data;
  }

  static async post<T>(
    command: Command,
    config?: AxiosRequestConfig<Command>,
  ): Promise<AxiosResponse<T>> {
    return await api.post("/command", command, config);
  }
}
