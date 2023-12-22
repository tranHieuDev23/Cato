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
import type { RpcGetAccountProblemSnippetListResponse } from './RpcGetAccountProblemSnippetListResponse';
import {
    RpcGetAccountProblemSnippetListResponseFromJSON,
    RpcGetAccountProblemSnippetListResponseFromJSONTyped,
    RpcGetAccountProblemSnippetListResponseToJSON,
} from './RpcGetAccountProblemSnippetListResponse';

/**
 * 
 * @export
 * @interface ResponseBodyOfTheGetAccountProblemSnippetListMethod
 */
export interface ResponseBodyOfTheGetAccountProblemSnippetListMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof ResponseBodyOfTheGetAccountProblemSnippetListMethod
     */
    jsonrpc?: ResponseBodyOfTheGetAccountProblemSnippetListMethodJsonrpcEnum;
    /**
     * It MUST be the same as the value of the id member in the Request.
     * @type {string}
     * @memberof ResponseBodyOfTheGetAccountProblemSnippetListMethod
     */
    id?: string;
    /**
     * 
     * @type {RpcError}
     * @memberof ResponseBodyOfTheGetAccountProblemSnippetListMethod
     */
    error?: RpcError;
    /**
     * 
     * @type {RpcGetAccountProblemSnippetListResponse}
     * @memberof ResponseBodyOfTheGetAccountProblemSnippetListMethod
     */
    result?: RpcGetAccountProblemSnippetListResponse;
}


/**
 * @export
 */
export const ResponseBodyOfTheGetAccountProblemSnippetListMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type ResponseBodyOfTheGetAccountProblemSnippetListMethodJsonrpcEnum = typeof ResponseBodyOfTheGetAccountProblemSnippetListMethodJsonrpcEnum[keyof typeof ResponseBodyOfTheGetAccountProblemSnippetListMethodJsonrpcEnum];


/**
 * Check if a given object implements the ResponseBodyOfTheGetAccountProblemSnippetListMethod interface.
 */
export function instanceOfResponseBodyOfTheGetAccountProblemSnippetListMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ResponseBodyOfTheGetAccountProblemSnippetListMethodFromJSON(json: any): ResponseBodyOfTheGetAccountProblemSnippetListMethod {
    return ResponseBodyOfTheGetAccountProblemSnippetListMethodFromJSONTyped(json, false);
}

export function ResponseBodyOfTheGetAccountProblemSnippetListMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): ResponseBodyOfTheGetAccountProblemSnippetListMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'error': !exists(json, 'error') ? undefined : RpcErrorFromJSON(json['error']),
        'result': !exists(json, 'result') ? undefined : RpcGetAccountProblemSnippetListResponseFromJSON(json['result']),
    };
}

export function ResponseBodyOfTheGetAccountProblemSnippetListMethodToJSON(value?: ResponseBodyOfTheGetAccountProblemSnippetListMethod | null): any {
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
        'result': RpcGetAccountProblemSnippetListResponseToJSON(value.result),
    };
}

