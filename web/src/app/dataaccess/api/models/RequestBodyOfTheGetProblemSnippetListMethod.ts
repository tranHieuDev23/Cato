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
import type { RpcGetProblemSnippetListRequest } from './RpcGetProblemSnippetListRequest';
import {
    RpcGetProblemSnippetListRequestFromJSON,
    RpcGetProblemSnippetListRequestFromJSONTyped,
    RpcGetProblemSnippetListRequestToJSON,
} from './RpcGetProblemSnippetListRequest';

/**
 * 
 * @export
 * @interface RequestBodyOfTheGetProblemSnippetListMethod
 */
export interface RequestBodyOfTheGetProblemSnippetListMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof RequestBodyOfTheGetProblemSnippetListMethod
     */
    jsonrpc?: RequestBodyOfTheGetProblemSnippetListMethodJsonrpcEnum;
    /**
     * An identifier established by the Client.
     * @type {string}
     * @memberof RequestBodyOfTheGetProblemSnippetListMethod
     */
    id?: string;
    /**
     * A String containing the name of the method to be invoked.
     * @type {string}
     * @memberof RequestBodyOfTheGetProblemSnippetListMethod
     */
    method?: RequestBodyOfTheGetProblemSnippetListMethodMethodEnum;
    /**
     * 
     * @type {RpcGetProblemSnippetListRequest}
     * @memberof RequestBodyOfTheGetProblemSnippetListMethod
     */
    params?: RpcGetProblemSnippetListRequest;
}


/**
 * @export
 */
export const RequestBodyOfTheGetProblemSnippetListMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type RequestBodyOfTheGetProblemSnippetListMethodJsonrpcEnum = typeof RequestBodyOfTheGetProblemSnippetListMethodJsonrpcEnum[keyof typeof RequestBodyOfTheGetProblemSnippetListMethodJsonrpcEnum];

/**
 * @export
 */
export const RequestBodyOfTheGetProblemSnippetListMethodMethodEnum = {
    GetProblemSnippetList: 'get_problem_snippet_list'
} as const;
export type RequestBodyOfTheGetProblemSnippetListMethodMethodEnum = typeof RequestBodyOfTheGetProblemSnippetListMethodMethodEnum[keyof typeof RequestBodyOfTheGetProblemSnippetListMethodMethodEnum];


/**
 * Check if a given object implements the RequestBodyOfTheGetProblemSnippetListMethod interface.
 */
export function instanceOfRequestBodyOfTheGetProblemSnippetListMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RequestBodyOfTheGetProblemSnippetListMethodFromJSON(json: any): RequestBodyOfTheGetProblemSnippetListMethod {
    return RequestBodyOfTheGetProblemSnippetListMethodFromJSONTyped(json, false);
}

export function RequestBodyOfTheGetProblemSnippetListMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestBodyOfTheGetProblemSnippetListMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'method': !exists(json, 'method') ? undefined : json['method'],
        'params': !exists(json, 'params') ? undefined : RpcGetProblemSnippetListRequestFromJSON(json['params']),
    };
}

export function RequestBodyOfTheGetProblemSnippetListMethodToJSON(value?: RequestBodyOfTheGetProblemSnippetListMethod | null): any {
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
        'params': RpcGetProblemSnippetListRequestToJSON(value.params),
    };
}

