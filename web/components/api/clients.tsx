import {ApiClient, BackupApi} from "../spec";

export function NewClient(): ApiClient {
    var url = window.location.protocol + "//" + window.location.hostname
    if (window.location.port) {
      url +=  ":" + window.location.port
    }
    var client = new ApiClient(url)
    client.defaultHeaders = [] // remove user-agent, causing browser error
    return client
}

export function NewBackupApi(): BackupApi {
    var client = NewClient()
    return new BackupApi(client)
}
