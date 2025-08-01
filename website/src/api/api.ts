/* tslint:disable */
/* eslint-disable */
/**
 * API V1
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import type { Configuration } from './configuration';
import type { AxiosPromise, AxiosInstance, RawAxiosRequestConfig } from 'axios';
import globalAxios from 'axios';
// Some imports not used depending on template conditions
// @ts-ignore
import { DUMMY_BASE_URL, assertParamExists, setApiKeyToObject, setBasicAuthToObject, setBearerAuthToObject, setOAuthToObject, setSearchParams, serializeDataIfNeeded, toPathString, createRequestFunction } from './common';
import type { RequestArgs } from './base';
// @ts-ignore
import { BASE_PATH, COLLECTION_FORMATS, BaseAPI, RequiredError, operationServerMap } from './base';

/**
 * 
 * @export
 * @interface EchoHTTPError
 */
export interface EchoHTTPError {
    /**
     * 
     * @type {object}
     * @memberof EchoHTTPError
     */
    'message'?: object;
}
/**
 * 
 * @export
 * @interface TimestamppbTimestamp
 */
export interface TimestamppbTimestamp {
    /**
     * Non-negative fractions of a second at nanosecond resolution. Negative second values with fractions must still have non-negative nanos values that count forward in time. Must be from 0 to 999,999,999 inclusive.
     * @type {number}
     * @memberof TimestamppbTimestamp
     */
    'nanos'?: number;
    /**
     * Represents seconds of UTC time since Unix epoch 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive.
     * @type {number}
     * @memberof TimestamppbTimestamp
     */
    'seconds'?: number;
}
/**
 * 
 * @export
 * @interface UserpbListUsersResponse
 */
export interface UserpbListUsersResponse {
    /**
     * 
     * @type {number}
     * @memberof UserpbListUsersResponse
     */
    'total'?: number;
    /**
     * 
     * @type {Array<UserpbUser>}
     * @memberof UserpbListUsersResponse
     */
    'users'?: Array<UserpbUser>;
}
/**
 * 
 * @export
 * @interface UserpbUser
 */
export interface UserpbUser {
    /**
     * 
     * @type {TimestamppbTimestamp}
     * @memberof UserpbUser
     */
    'created_at'?: TimestamppbTimestamp;
    /**
     * 
     * @type {string}
     * @memberof UserpbUser
     */
    'email'?: string;
    /**
     * 
     * @type {string}
     * @memberof UserpbUser
     */
    'github_id'?: string;
    /**
     * 
     * @type {number}
     * @memberof UserpbUser
     */
    'id'?: number;
    /**
     * 
     * @type {string}
     * @memberof UserpbUser
     */
    'name'?: string;
    /**
     * 
     * @type {UserpbUserRole}
     * @memberof UserpbUser
     */
    'role'?: UserpbUserRole;
    /**
     * 
     * @type {TimestamppbTimestamp}
     * @memberof UserpbUser
     */
    'updated_at'?: TimestamppbTimestamp;
}


/**
 * 
 * @export
 * @enum {number}
 */

export const UserpbUserRole = {
    UserRole_USER: 0,
    UserRole_ADMIN: 1
} as const;

export type UserpbUserRole = typeof UserpbUserRole[keyof typeof UserpbUserRole];



/**
 * UserApi - axios parameter creator
 * @export
 */
export const UserApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * Retrieve a paginated list of all users (requires admin privileges)
         * @summary List users
         * @param {number} [page] Page number (default: 1)
         * @param {number} [pageSize] Page size (default: 10)
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        userListGet: async (page?: number, pageSize?: number, options: RawAxiosRequestConfig = {}): Promise<RequestArgs> => {
            const localVarPath = `/user/list`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            if (page !== undefined) {
                localVarQueryParameter['page'] = page;
            }

            if (pageSize !== undefined) {
                localVarQueryParameter['page_size'] = pageSize;
            }


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * Retrieve the information of the currently authenticated user
         * @summary Get current user
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        userMeGet: async (options: RawAxiosRequestConfig = {}): Promise<RequestArgs> => {
            const localVarPath = `/user/me`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
    }
};

/**
 * UserApi - functional programming interface
 * @export
 */
export const UserApiFp = function(configuration?: Configuration) {
    const localVarAxiosParamCreator = UserApiAxiosParamCreator(configuration)
    return {
        /**
         * Retrieve a paginated list of all users (requires admin privileges)
         * @summary List users
         * @param {number} [page] Page number (default: 1)
         * @param {number} [pageSize] Page size (default: 10)
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async userListGet(page?: number, pageSize?: number, options?: RawAxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<UserpbListUsersResponse>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.userListGet(page, pageSize, options);
            const localVarOperationServerIndex = configuration?.serverIndex ?? 0;
            const localVarOperationServerBasePath = operationServerMap['UserApi.userListGet']?.[localVarOperationServerIndex]?.url;
            return (axios, basePath) => createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration)(axios, localVarOperationServerBasePath || basePath);
        },
        /**
         * Retrieve the information of the currently authenticated user
         * @summary Get current user
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async userMeGet(options?: RawAxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<UserpbUser>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.userMeGet(options);
            const localVarOperationServerIndex = configuration?.serverIndex ?? 0;
            const localVarOperationServerBasePath = operationServerMap['UserApi.userMeGet']?.[localVarOperationServerIndex]?.url;
            return (axios, basePath) => createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration)(axios, localVarOperationServerBasePath || basePath);
        },
    }
};

/**
 * UserApi - factory interface
 * @export
 */
export const UserApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    const localVarFp = UserApiFp(configuration)
    return {
        /**
         * Retrieve a paginated list of all users (requires admin privileges)
         * @summary List users
         * @param {number} [page] Page number (default: 1)
         * @param {number} [pageSize] Page size (default: 10)
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        userListGet(page?: number, pageSize?: number, options?: RawAxiosRequestConfig): AxiosPromise<UserpbListUsersResponse> {
            return localVarFp.userListGet(page, pageSize, options).then((request) => request(axios, basePath));
        },
        /**
         * Retrieve the information of the currently authenticated user
         * @summary Get current user
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        userMeGet(options?: RawAxiosRequestConfig): AxiosPromise<UserpbUser> {
            return localVarFp.userMeGet(options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * UserApi - object-oriented interface
 * @export
 * @class UserApi
 * @extends {BaseAPI}
 */
export class UserApi extends BaseAPI {
    /**
     * Retrieve a paginated list of all users (requires admin privileges)
     * @summary List users
     * @param {number} [page] Page number (default: 1)
     * @param {number} [pageSize] Page size (default: 10)
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof UserApi
     */
    public userListGet(page?: number, pageSize?: number, options?: RawAxiosRequestConfig) {
        return UserApiFp(this.configuration).userListGet(page, pageSize, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * Retrieve the information of the currently authenticated user
     * @summary Get current user
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof UserApi
     */
    public userMeGet(options?: RawAxiosRequestConfig) {
        return UserApiFp(this.configuration).userMeGet(options).then((request) => request(this.axios, this.basePath));
    }
}



