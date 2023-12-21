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
import type { RpcUpdateUserRequest } from './RpcUpdateUserRequest';
import {
    RpcUpdateUserRequestFromJSON,
    RpcUpdateUserRequestFromJSONTyped,
    RpcUpdateUserRequestToJSON,
} from './RpcUpdateUserRequest';

/**
 * 
 * @export
 * @interface RequestBodyOfTheUpdateUserMethod
 */
export interface RequestBodyOfTheUpdateUserMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof RequestBodyOfTheUpdateUserMethod
     */
    jsonrpc?: RequestBodyOfTheUpdateUserMethodJsonrpcEnum;
    /**
     * An identifier established by the Client.
     * @type {string}
     * @memberof RequestBodyOfTheUpdateUserMethod
     */
    id?: string;
    /**
     * A String containing the name of the method to be invoked.
     * @type {string}
     * @memberof RequestBodyOfTheUpdateUserMethod
     */
    method?: RequestBodyOfTheUpdateUserMethodMethodEnum;
    /**
     * 
     * @type {RpcUpdateUserRequest}
     * @memberof RequestBodyOfTheUpdateUserMethod
     */
    params?: RpcUpdateUserRequest;
}


/**
 * @export
 */
export const RequestBodyOfTheUpdateUserMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type RequestBodyOfTheUpdateUserMethodJsonrpcEnum = typeof RequestBodyOfTheUpdateUserMethodJsonrpcEnum[keyof typeof RequestBodyOfTheUpdateUserMethodJsonrpcEnum];

/**
 * @export
 */
export const RequestBodyOfTheUpdateUserMethodMethodEnum = {
    UpdateUser: 'update_user'
} as const;
export type RequestBodyOfTheUpdateUserMethodMethodEnum = typeof RequestBodyOfTheUpdateUserMethodMethodEnum[keyof typeof RequestBodyOfTheUpdateUserMethodMethodEnum];


/**
 * Check if a given object implements the RequestBodyOfTheUpdateUserMethod interface.
 */
export function instanceOfRequestBodyOfTheUpdateUserMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RequestBodyOfTheUpdateUserMethodFromJSON(json: any): RequestBodyOfTheUpdateUserMethod {
    return RequestBodyOfTheUpdateUserMethodFromJSONTyped(json, false);
}

export function RequestBodyOfTheUpdateUserMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestBodyOfTheUpdateUserMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'method': !exists(json, 'method') ? undefined : json['method'],
        'params': !exists(json, 'params') ? undefined : RpcUpdateUserRequestFromJSON(json['params']),
    };
}

export function RequestBodyOfTheUpdateUserMethodToJSON(value?: RequestBodyOfTheUpdateUserMethod | null): any {
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
        'params': RpcUpdateUserRequestToJSON(value.params),
    };
}

