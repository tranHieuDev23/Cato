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
import type { RpcUpdateSettingResponse } from './RpcUpdateSettingResponse';
import {
    RpcUpdateSettingResponseFromJSON,
    RpcUpdateSettingResponseFromJSONTyped,
    RpcUpdateSettingResponseToJSON,
} from './RpcUpdateSettingResponse';

/**
 * 
 * @export
 * @interface ResponseBodyOfTheUpdateSettingMethod
 */
export interface ResponseBodyOfTheUpdateSettingMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof ResponseBodyOfTheUpdateSettingMethod
     */
    jsonrpc?: ResponseBodyOfTheUpdateSettingMethodJsonrpcEnum;
    /**
     * It MUST be the same as the value of the id member in the Request.
     * @type {string}
     * @memberof ResponseBodyOfTheUpdateSettingMethod
     */
    id?: string;
    /**
     * 
     * @type {RpcError}
     * @memberof ResponseBodyOfTheUpdateSettingMethod
     */
    error?: RpcError;
    /**
     * 
     * @type {RpcUpdateSettingResponse}
     * @memberof ResponseBodyOfTheUpdateSettingMethod
     */
    result?: RpcUpdateSettingResponse;
}


/**
 * @export
 */
export const ResponseBodyOfTheUpdateSettingMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type ResponseBodyOfTheUpdateSettingMethodJsonrpcEnum = typeof ResponseBodyOfTheUpdateSettingMethodJsonrpcEnum[keyof typeof ResponseBodyOfTheUpdateSettingMethodJsonrpcEnum];


/**
 * Check if a given object implements the ResponseBodyOfTheUpdateSettingMethod interface.
 */
export function instanceOfResponseBodyOfTheUpdateSettingMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ResponseBodyOfTheUpdateSettingMethodFromJSON(json: any): ResponseBodyOfTheUpdateSettingMethod {
    return ResponseBodyOfTheUpdateSettingMethodFromJSONTyped(json, false);
}

export function ResponseBodyOfTheUpdateSettingMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): ResponseBodyOfTheUpdateSettingMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'error': !exists(json, 'error') ? undefined : RpcErrorFromJSON(json['error']),
        'result': !exists(json, 'result') ? undefined : RpcUpdateSettingResponseFromJSON(json['result']),
    };
}

export function ResponseBodyOfTheUpdateSettingMethodToJSON(value?: ResponseBodyOfTheUpdateSettingMethod | null): any {
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
        'result': RpcUpdateSettingResponseToJSON(value.result),
    };
}

