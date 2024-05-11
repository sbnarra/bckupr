# BackupApi

All URIs are relative to *http://localhost:8000*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**deleteBackup**](BackupApi.md#deleteBackup) | **DELETE** /api/backups/{id} | Deletes backup |
| [**getBackup**](BackupApi.md#getBackup) | **GET** /api/backups/{id} | Gets backup by id |
| [**getRestore**](BackupApi.md#getRestore) | **GET** /api/backups/{id}/restore |  |
| [**getRotate**](BackupApi.md#getRotate) | **GET** /api/rotate |  |
| [**listBackups**](BackupApi.md#listBackups) | **GET** /api/backups |  |
| [**startBackup**](BackupApi.md#startBackup) | **POST** /api/backups | Creates new backup |
| [**startBackupWithId**](BackupApi.md#startBackupWithId) | **PUT** /api/backups/{id} |  |
| [**startRestore**](BackupApi.md#startRestore) | **POST** /api/backups/{id}/restore |  |
| [**startRotate**](BackupApi.md#startRotate) | **POST** /api/rotate | Retrieves application version |


<a name="deleteBackup"></a>
# **deleteBackup**
> deleteBackup(id)

Deletes backup

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **String**| backup id | [default to null] |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="getBackup"></a>
# **getBackup**
> Backup getBackup(id)

Gets backup by id

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **String**| backup id | [default to null] |

### Return type

[**Backup**](../Models/Backup.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getRestore"></a>
# **getRestore**
> Restore getRestore(id)



### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **String**| backup id | [default to null] |

### Return type

[**Restore**](../Models/Restore.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getRotate"></a>
# **getRotate**
> Rotate getRotate()



### Parameters
This endpoint does not need any parameter.

### Return type

[**Rotate**](../Models/Rotate.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="listBackups"></a>
# **listBackups**
> List listBackups()



### Parameters
This endpoint does not need any parameter.

### Return type

[**List**](../Models/Backup.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="startBackup"></a>
# **startBackup**
> Backup startBackup(TaskInput)

Creates new backup

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **TaskInput** | [**TaskInput**](../Models/TaskInput.md)| Backup to be created | |

### Return type

[**Backup**](../Models/Backup.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="startBackupWithId"></a>
# **startBackupWithId**
> Backup startBackupWithId(id, TaskInput)



### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **String**| backup id | [default to null] |
| **TaskInput** | [**TaskInput**](../Models/TaskInput.md)| Backup to be created | [optional] |

### Return type

[**Backup**](../Models/Backup.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="startRestore"></a>
# **startRestore**
> Restore startRestore(id, TaskInput)



### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **String**| backup id | [default to null] |
| **TaskInput** | [**TaskInput**](../Models/TaskInput.md)| Backup to be created | |

### Return type

[**Restore**](../Models/Restore.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="startRotate"></a>
# **startRotate**
> Rotate startRotate(RotateInput)

Retrieves application version

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **RotateInput** | [**RotateInput**](../Models/RotateInput.md)| Backup to be created | [optional] |

### Return type

[**Rotate**](../Models/Rotate.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

