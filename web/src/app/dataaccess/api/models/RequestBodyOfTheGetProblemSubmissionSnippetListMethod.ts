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
import type { RpcGetProblemSubmissionSnippetListRequest } from './RpcGetProblemSubmissionSnippetListRequest';
import {
    RpcGetProblemSubmissionSnippetListRequestFromJSON,
    RpcGetProblemSubmissionSnippetListRequestFromJSONTyped,
    RpcGetProblemSubmissionSnippetListRequestToJSON,
} from './RpcGetProblemSubmissionSnippetListRequest';

/**
 * 
 * @export
 * @interface RequestBodyOfTheGetProblemSubmissionSnippetListMethod
 */
export interface RequestBodyOfTheGetProblemSubmissionSnippetListMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof RequestBodyOfTheGetProblemSubmissionSnippetListMethod
     */
    jsonrpc?: RequestBodyOfTheGetProblemSubmissionSnippetListMethodJsonrpcEnum;
    /**
     * An identifier established by the Client.
     * @type {string}
     * @memberof RequestBodyOfTheGetProblemSubmissionSnippetListMethod
     */
    id?: string;
    /**
     * A String containing the name of the method to be invoked.
     * @type {string}
     * @memberof RequestBodyOfTheGetProblemSubmissionSnippetListMethod
     */
    method?: RequestBodyOfTheGetProblemSubmissionSnippetListMethodMethodEnum;
    /**
     * 
     * @type {RpcGetProblemSubmissionSnippetListRequest}
     * @memberof RequestBodyOfTheGetProblemSubmissionSnippetListMethod
     */
    params?: RpcGetProblemSubmissionSnippetListRequest;
}


/**
 * @export
 */
export const RequestBodyOfTheGetProblemSubmissionSnippetListMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type RequestBodyOfTheGetProblemSubmissionSnippetListMethodJsonrpcEnum = typeof RequestBodyOfTheGetProblemSubmissionSnippetListMethodJsonrpcEnum[keyof typeof RequestBodyOfTheGetProblemSubmissionSnippetListMethodJsonrpcEnum];

/**
 * @export
 */
export const RequestBodyOfTheGetProblemSubmissionSnippetListMethodMethodEnum = {
    GetProblemSubmissionSnippetList: 'get_problem_submission_snippet_list'
} as const;
export type RequestBodyOfTheGetProblemSubmissionSnippetListMethodMethodEnum = typeof RequestBodyOfTheGetProblemSubmissionSnippetListMethodMethodEnum[keyof typeof RequestBodyOfTheGetProblemSubmissionSnippetListMethodMethodEnum];


/**
 * Check if a given object implements the RequestBodyOfTheGetProblemSubmissionSnippetListMethod interface.
 */
export function instanceOfRequestBodyOfTheGetProblemSubmissionSnippetListMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RequestBodyOfTheGetProblemSubmissionSnippetListMethodFromJSON(json: any): RequestBodyOfTheGetProblemSubmissionSnippetListMethod {
    return RequestBodyOfTheGetProblemSubmissionSnippetListMethodFromJSONTyped(json, false);
}

export function RequestBodyOfTheGetProblemSubmissionSnippetListMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestBodyOfTheGetProblemSubmissionSnippetListMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'method': !exists(json, 'method') ? undefined : json['method'],
        'params': !exists(json, 'params') ? undefined : RpcGetProblemSubmissionSnippetListRequestFromJSON(json['params']),
    };
}

export function RequestBodyOfTheGetProblemSubmissionSnippetListMethodToJSON(value?: RequestBodyOfTheGetProblemSubmissionSnippetListMethod | null): any {
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
        'params': RpcGetProblemSubmissionSnippetListRequestToJSON(value.params),
    };
}

