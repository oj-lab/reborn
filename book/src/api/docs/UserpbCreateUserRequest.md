# UserpbCreateUserRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**email** | **string** |  | [optional] [default to undefined]
**name** | **string** |  | [optional] [default to undefined]
**password** | **string** | Password, optional, encrypted storage | [optional] [default to undefined]
**role** | [**UserpbUserRole**](UserpbUserRole.md) |  | [optional] [default to undefined]

## Example

```typescript
import { UserpbCreateUserRequest } from './api';

const instance: UserpbCreateUserRequest = {
    email,
    name,
    password,
    role,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
