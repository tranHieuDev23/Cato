/* tslint:disable */
/* eslint-disable */
/**
 * API
 * Generated by genpjrpc: v0.4.0
 *
 * The version of the OpenAPI document: v0.0.0-unknown
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
import type { RpcError } from './RpcError';
import {
    RpcErrorFromJSON,
    RpcErrorFromJSONTyped,
    RpcErrorToJSON,
} from './RpcError';
import type { RpcGetUserResponse } from './RpcGetUserResponse';
import {
    RpcGetUserResponseFromJSON,
    RpcGetUserResponseFromJSONTyped,
    RpcGetUserResponseToJSON,
} from './RpcGetUserResponse';

/**
 * 
 * @export
 * @interface ResponseBodyOfTheGetUserMethod
 */
export interface ResponseBodyOfTheGetUserMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof ResponseBodyOfTheGetUserMethod
     */
    jsonrpc?: ResponseBodyOfTheGetUserMethodJsonrpcEnum;
    /**
     * It MUST be the same as the value of the id member in the Request.
     * @type {string}
     * @memberof ResponseBodyOfTheGetUserMethod
     */
    id?: string;
    /**
     * 
     * @type {RpcError}
     * @memberof ResponseBodyOfTheGetUserMethod
     */
    error?: RpcError;
    /**
     * 
     * @type {RpcGetUserResponse}
     * @memberof ResponseBodyOfTheGetUserMethod
     */
    result?: RpcGetUserResponse;
}


/**
 * @export
 */
export const ResponseBodyOfTheGetUserMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type ResponseBodyOfTheGetUserMethodJsonrpcEnum = typeof ResponseBodyOfTheGetUserMethodJsonrpcEnum[keyof typeof ResponseBodyOfTheGetUserMethodJsonrpcEnum];


/**
 * Check if a given object implements the ResponseBodyOfTheGetUserMethod interface.
 */
export function instanceOfResponseBodyOfTheGetUserMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ResponseBodyOfTheGetUserMethodFromJSON(json: any): ResponseBodyOfTheGetUserMethod {
    return ResponseBodyOfTheGetUserMethodFromJSONTyped(json, false);
}

export function ResponseBodyOfTheGetUserMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): ResponseBodyOfTheGetUserMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'error': !exists(json, 'error') ? undefined : RpcErrorFromJSON(json['error']),
        'result': !exists(json, 'result') ? undefined : RpcGetUserResponseFromJSON(json['result']),
    };
}

export function ResponseBodyOfTheGetUserMethodToJSON(value?: ResponseBodyOfTheGetUserMethod | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'jsonrpc': value.jsonrpc,
        'id': value.id,
        'error': RpcErrorToJSON(value.error),
        'result': RpcGetUserResponseToJSON(value.result),
    };
}

