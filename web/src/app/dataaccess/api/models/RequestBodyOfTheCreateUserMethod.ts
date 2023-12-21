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
import type { RpcCreateUserRequest } from './RpcCreateUserRequest';
import {
    RpcCreateUserRequestFromJSON,
    RpcCreateUserRequestFromJSONTyped,
    RpcCreateUserRequestToJSON,
} from './RpcCreateUserRequest';

/**
 * 
 * @export
 * @interface RequestBodyOfTheCreateUserMethod
 */
export interface RequestBodyOfTheCreateUserMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof RequestBodyOfTheCreateUserMethod
     */
    jsonrpc?: RequestBodyOfTheCreateUserMethodJsonrpcEnum;
    /**
     * An identifier established by the Client.
     * @type {string}
     * @memberof RequestBodyOfTheCreateUserMethod
     */
    id?: string;
    /**
     * A String containing the name of the method to be invoked.
     * @type {string}
     * @memberof RequestBodyOfTheCreateUserMethod
     */
    method?: RequestBodyOfTheCreateUserMethodMethodEnum;
    /**
     * 
     * @type {RpcCreateUserRequest}
     * @memberof RequestBodyOfTheCreateUserMethod
     */
    params?: RpcCreateUserRequest;
}


/**
 * @export
 */
export const RequestBodyOfTheCreateUserMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type RequestBodyOfTheCreateUserMethodJsonrpcEnum = typeof RequestBodyOfTheCreateUserMethodJsonrpcEnum[keyof typeof RequestBodyOfTheCreateUserMethodJsonrpcEnum];

/**
 * @export
 */
export const RequestBodyOfTheCreateUserMethodMethodEnum = {
    CreateUser: 'create_user'
} as const;
export type RequestBodyOfTheCreateUserMethodMethodEnum = typeof RequestBodyOfTheCreateUserMethodMethodEnum[keyof typeof RequestBodyOfTheCreateUserMethodMethodEnum];


/**
 * Check if a given object implements the RequestBodyOfTheCreateUserMethod interface.
 */
export function instanceOfRequestBodyOfTheCreateUserMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RequestBodyOfTheCreateUserMethodFromJSON(json: any): RequestBodyOfTheCreateUserMethod {
    return RequestBodyOfTheCreateUserMethodFromJSONTyped(json, false);
}

export function RequestBodyOfTheCreateUserMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestBodyOfTheCreateUserMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'method': !exists(json, 'method') ? undefined : json['method'],
        'params': !exists(json, 'params') ? undefined : RpcCreateUserRequestFromJSON(json['params']),
    };
}

export function RequestBodyOfTheCreateUserMethodToJSON(value?: RequestBodyOfTheCreateUserMethod | null): any {
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
        'params': RpcCreateUserRequestToJSON(value.params),
    };
}

