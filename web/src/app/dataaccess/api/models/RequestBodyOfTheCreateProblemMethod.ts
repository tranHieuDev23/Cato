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
import type { RpcCreateProblemRequest } from './RpcCreateProblemRequest';
import {
    RpcCreateProblemRequestFromJSON,
    RpcCreateProblemRequestFromJSONTyped,
    RpcCreateProblemRequestToJSON,
} from './RpcCreateProblemRequest';

/**
 * 
 * @export
 * @interface RequestBodyOfTheCreateProblemMethod
 */
export interface RequestBodyOfTheCreateProblemMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof RequestBodyOfTheCreateProblemMethod
     */
    jsonrpc?: RequestBodyOfTheCreateProblemMethodJsonrpcEnum;
    /**
     * An identifier established by the Client.
     * @type {string}
     * @memberof RequestBodyOfTheCreateProblemMethod
     */
    id?: string;
    /**
     * A String containing the name of the method to be invoked.
     * @type {string}
     * @memberof RequestBodyOfTheCreateProblemMethod
     */
    method?: RequestBodyOfTheCreateProblemMethodMethodEnum;
    /**
     * 
     * @type {RpcCreateProblemRequest}
     * @memberof RequestBodyOfTheCreateProblemMethod
     */
    params?: RpcCreateProblemRequest;
}


/**
 * @export
 */
export const RequestBodyOfTheCreateProblemMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type RequestBodyOfTheCreateProblemMethodJsonrpcEnum = typeof RequestBodyOfTheCreateProblemMethodJsonrpcEnum[keyof typeof RequestBodyOfTheCreateProblemMethodJsonrpcEnum];

/**
 * @export
 */
export const RequestBodyOfTheCreateProblemMethodMethodEnum = {
    CreateProblem: 'create_problem'
} as const;
export type RequestBodyOfTheCreateProblemMethodMethodEnum = typeof RequestBodyOfTheCreateProblemMethodMethodEnum[keyof typeof RequestBodyOfTheCreateProblemMethodMethodEnum];


/**
 * Check if a given object implements the RequestBodyOfTheCreateProblemMethod interface.
 */
export function instanceOfRequestBodyOfTheCreateProblemMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RequestBodyOfTheCreateProblemMethodFromJSON(json: any): RequestBodyOfTheCreateProblemMethod {
    return RequestBodyOfTheCreateProblemMethodFromJSONTyped(json, false);
}

export function RequestBodyOfTheCreateProblemMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestBodyOfTheCreateProblemMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'method': !exists(json, 'method') ? undefined : json['method'],
        'params': !exists(json, 'params') ? undefined : RpcCreateProblemRequestFromJSON(json['params']),
    };
}

export function RequestBodyOfTheCreateProblemMethodToJSON(value?: RequestBodyOfTheCreateProblemMethod | null): any {
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
        'params': RpcCreateProblemRequestToJSON(value.params),
    };
}

