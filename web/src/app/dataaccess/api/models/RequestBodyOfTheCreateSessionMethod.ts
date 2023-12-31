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
import type { RpcCreateSessionRequest } from './RpcCreateSessionRequest';
import {
    RpcCreateSessionRequestFromJSON,
    RpcCreateSessionRequestFromJSONTyped,
    RpcCreateSessionRequestToJSON,
} from './RpcCreateSessionRequest';

/**
 * 
 * @export
 * @interface RequestBodyOfTheCreateSessionMethod
 */
export interface RequestBodyOfTheCreateSessionMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof RequestBodyOfTheCreateSessionMethod
     */
    jsonrpc?: RequestBodyOfTheCreateSessionMethodJsonrpcEnum;
    /**
     * An identifier established by the Client.
     * @type {string}
     * @memberof RequestBodyOfTheCreateSessionMethod
     */
    id?: string;
    /**
     * A String containing the name of the method to be invoked.
     * @type {string}
     * @memberof RequestBodyOfTheCreateSessionMethod
     */
    method?: RequestBodyOfTheCreateSessionMethodMethodEnum;
    /**
     * 
     * @type {RpcCreateSessionRequest}
     * @memberof RequestBodyOfTheCreateSessionMethod
     */
    params?: RpcCreateSessionRequest;
}


/**
 * @export
 */
export const RequestBodyOfTheCreateSessionMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type RequestBodyOfTheCreateSessionMethodJsonrpcEnum = typeof RequestBodyOfTheCreateSessionMethodJsonrpcEnum[keyof typeof RequestBodyOfTheCreateSessionMethodJsonrpcEnum];

/**
 * @export
 */
export const RequestBodyOfTheCreateSessionMethodMethodEnum = {
    CreateSession: 'create_session'
} as const;
export type RequestBodyOfTheCreateSessionMethodMethodEnum = typeof RequestBodyOfTheCreateSessionMethodMethodEnum[keyof typeof RequestBodyOfTheCreateSessionMethodMethodEnum];


/**
 * Check if a given object implements the RequestBodyOfTheCreateSessionMethod interface.
 */
export function instanceOfRequestBodyOfTheCreateSessionMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RequestBodyOfTheCreateSessionMethodFromJSON(json: any): RequestBodyOfTheCreateSessionMethod {
    return RequestBodyOfTheCreateSessionMethodFromJSONTyped(json, false);
}

export function RequestBodyOfTheCreateSessionMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestBodyOfTheCreateSessionMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'method': !exists(json, 'method') ? undefined : json['method'],
        'params': !exists(json, 'params') ? undefined : RpcCreateSessionRequestFromJSON(json['params']),
    };
}

export function RequestBodyOfTheCreateSessionMethodToJSON(value?: RequestBodyOfTheCreateSessionMethod | null): any {
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
        'params': RpcCreateSessionRequestToJSON(value.params),
    };
}

