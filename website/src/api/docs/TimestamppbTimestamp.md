# TimestamppbTimestamp


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**nanos** | **number** | Non-negative fractions of a second at nanosecond resolution. Negative second values with fractions must still have non-negative nanos values that count forward in time. Must be from 0 to 999,999,999 inclusive. | [optional] [default to undefined]
**seconds** | **number** | Represents seconds of UTC time since Unix epoch 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive. | [optional] [default to undefined]

## Example

```typescript
import { TimestamppbTimestamp } from './api';

const instance: TimestamppbTimestamp = {
    nanos,
    seconds,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
