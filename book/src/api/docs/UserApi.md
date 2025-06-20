# UserApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**userPost**](#userpost) | **POST** /user | Create a new user|

# **userPost**
> userPost(user)

Create a new user with the provided details

### Example

```typescript
import {
    UserApi,
    Configuration,
    UserpbCreateUserRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new UserApi(configuration);

let user: UserpbCreateUserRequest; //User details

const { status, data } = await apiInstance.userPost(
    user
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **user** | **UserpbCreateUserRequest**| User details | |


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

