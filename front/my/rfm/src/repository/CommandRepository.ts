import {api} from "../api/api.ts";
import {Command} from "../model/Command.ts";

export class CommandRepository {
    static async  post<T>(command: Command): Promise<T>{
       return (await api.post("/command", command)).data
    }
}