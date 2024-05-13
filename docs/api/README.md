# Documentation for Bckupr

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost:8000*

| Class | Method | HTTP request | Description |
|------------ | ------------- | ------------- | -------------|
| *BackupApi* | [**deleteBackup**](Apis/BackupApi.md#deletebackup) | **DELETE** /api/backups/{id} | Deletes backup |
*BackupApi* | [**getBackup**](Apis/BackupApi.md#getbackup) | **GET** /api/backups/{id} | Gets backup by id |
*BackupApi* | [**getRestore**](Apis/BackupApi.md#getrestore) | **GET** /api/backups/{id}/restore |  |
*BackupApi* | [**getRotate**](Apis/BackupApi.md#getrotate) | **GET** /api/rotate |  |
*BackupApi* | [**listBackups**](Apis/BackupApi.md#listbackups) | **GET** /api/backups |  |
*BackupApi* | [**startBackup**](Apis/BackupApi.md#startbackup) | **POST** /api/backups | Creates new backup |
*BackupApi* | [**startBackupWithId**](Apis/BackupApi.md#startbackupwithid) | **PUT** /api/backups/{id} |  |
*BackupApi* | [**startRestore**](Apis/BackupApi.md#startrestore) | **POST** /api/backups/{id}/restore |  |
*BackupApi* | [**startRotate**](Apis/BackupApi.md#startrotate) | **POST** /api/rotate | Retrieves application version |
| *SystemApi* | [**getVersion**](Apis/SystemApi.md#getversion) | **GET** /api/version | Retrieves application version |


<a name="documentation-for-models"></a>
## Documentation for Models

 - [Backup](./Models/Backup.md)
 - [Error](./Models/Error.md)
 - [Filters](./Models/Filters.md)
 - [Restore](./Models/Restore.md)
 - [Rotate](./Models/Rotate.md)
 - [RotateInput](./Models/RotateInput.md)
 - [Status](./Models/Status.md)
 - [StopModes](./Models/StopModes.md)
 - [TaskInput](./Models/TaskInput.md)
 - [Version](./Models/Version.md)
 - [Volume](./Models/Volume.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
