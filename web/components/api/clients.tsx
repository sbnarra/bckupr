import {BackupApi, Configuration} from "../spec";

export function basePath(): string {
    var url = window.location.protocol + "//" + window.location.hostname
    if (window.location.port) {
      url +=  ":" + window.location.port
    }
    return url
}

export function NewBackupApi(): BackupApi {
    var config = new Configuration({
      basePath: basePath()
    })
    return new BackupApi(config)
}
