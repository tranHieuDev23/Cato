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
import type { RpcUpdateSettingRequest } from './RpcUpdateSettingRequest';
import {
    RpcUpdateSettingRequestFromJSON,
    RpcUpdateSettingRequestFromJSONTyped,
    RpcUpdateSettingRequestToJSON,
} from './RpcUpdateSettingRequest';

/**
 * 
 * @export
 * @interface RequestBodyOfTheUpdateSettingMethod
 */
export interface RequestBodyOfTheUpdateSettingMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof RequestBodyOfTheUpdateSettingMethod
     */
    jsonrpc?: RequestBodyOfTheUpdateSettingMethodJsonrpcEnum;
    /**
     * An identifier established by the Client.
     * @type {string}
     * @memberof RequestBodyOfTheUpdateSettingMethod
     */
    id?: string;
    /**
     * A String containing the name of the method to be invoked.
     * @type {string}
     * @memberof RequestBodyOfTheUpdateSettingMethod
     */
    method?: RequestBodyOfTheUpdateSettingMethodMethodEnum;
    /**
     * 
     * @type {RpcUpdateSettingRequest}
     * @memberof RequestBodyOfTheUpdateSettingMethod
     */
    params?: RpcUpdateSettingRequest;
}


/**
 * @export
 */
export const RequestBodyOfTheUpdateSettingMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type RequestBodyOfTheUpdateSettingMethodJsonrpcEnum = typeof RequestBodyOfTheUpdateSettingMethodJsonrpcEnum[keyof typeof RequestBodyOfTheUpdateSettingMethodJsonrpcEnum];

/**
 * @export
 */
export const RequestBodyOfTheUpdateSettingMethodMethodEnum = {
    UpdateSetting: 'update_setting'
} as const;
export type RequestBodyOfTheUpdateSettingMethodMethodEnum = typeof RequestBodyOfTheUpdateSettingMethodMethodEnum[keyof typeof RequestBodyOfTheUpdateSettingMethodMethodEnum];


/**
 * Check if a given object implements the RequestBodyOfTheUpdateSettingMethod interface.
 */
export function instanceOfRequestBodyOfTheUpdateSettingMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RequestBodyOfTheUpdateSettingMethodFromJSON(json: any): RequestBodyOfTheUpdateSettingMethod {
    return RequestBodyOfTheUpdateSettingMethodFromJSONTyped(json, false);
}

export function RequestBodyOfTheUpdateSettingMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestBodyOfTheUpdateSettingMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'method': !exists(json, 'method') ? undefined : json['method'],
        'params': !exists(json, 'params') ? undefined : RpcUpdateSettingRequestFromJSON(json['params']),
    };
}

export function RequestBodyOfTheUpdateSettingMethodToJSON(value?: RequestBodyOfTheUpdateSettingMethod | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'jsonrpc': value.jsonrpc,
        'id': value.id,
        'method': value.method,
        'params': RpcUpdateSettingRequestToJSON(value.params),
    };
}
