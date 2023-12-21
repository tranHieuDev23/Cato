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
import type { RpcGetProblemSnippetListResponse } from './RpcGetProblemSnippetListResponse';
import {
    RpcGetProblemSnippetListResponseFromJSON,
    RpcGetProblemSnippetListResponseFromJSONTyped,
    RpcGetProblemSnippetListResponseToJSON,
} from './RpcGetProblemSnippetListResponse';

/**
 * 
 * @export
 * @interface ResponseBodyOfTheGetProblemSnippetListMethod
 */
export interface ResponseBodyOfTheGetProblemSnippetListMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof ResponseBodyOfTheGetProblemSnippetListMethod
     */
    jsonrpc?: ResponseBodyOfTheGetProblemSnippetListMethodJsonrpcEnum;
    /**
     * It MUST be the same as the value of the id member in the Request.
     * @type {string}
     * @memberof ResponseBodyOfTheGetProblemSnippetListMethod
     */
    id?: string;
    /**
     * 
     * @type {RpcError}
     * @memberof ResponseBodyOfTheGetProblemSnippetListMethod
     */
    error?: RpcError;
    /**
     * 
     * @type {RpcGetProblemSnippetListResponse}
     * @memberof ResponseBodyOfTheGetProblemSnippetListMethod
     */
    result?: RpcGetProblemSnippetListResponse;
}


/**
 * @export
 */
export const ResponseBodyOfTheGetProblemSnippetListMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type ResponseBodyOfTheGetProblemSnippetListMethodJsonrpcEnum = typeof ResponseBodyOfTheGetProblemSnippetListMethodJsonrpcEnum[keyof typeof ResponseBodyOfTheGetProblemSnippetListMethodJsonrpcEnum];


/**
 * Check if a given object implements the ResponseBodyOfTheGetProblemSnippetListMethod interface.
 */
export function instanceOfResponseBodyOfTheGetProblemSnippetListMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ResponseBodyOfTheGetProblemSnippetListMethodFromJSON(json: any): ResponseBodyOfTheGetProblemSnippetListMethod {
    return ResponseBodyOfTheGetProblemSnippetListMethodFromJSONTyped(json, false);
}

export function ResponseBodyOfTheGetProblemSnippetListMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): ResponseBodyOfTheGetProblemSnippetListMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'error': !exists(json, 'error') ? undefined : RpcErrorFromJSON(json['error']),
        'result': !exists(json, 'result') ? undefined : RpcGetProblemSnippetListResponseFromJSON(json['result']),
    };
}

export function ResponseBodyOfTheGetProblemSnippetListMethodToJSON(value?: ResponseBodyOfTheGetProblemSnippetListMethod | null): any {
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
        'result': RpcGetProblemSnippetListResponseToJSON(value.result),
    };
}

