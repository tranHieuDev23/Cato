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
import type { RpcGetAccountProblemSubmissionSnippetListRequest } from './RpcGetAccountProblemSubmissionSnippetListRequest';
import {
    RpcGetAccountProblemSubmissionSnippetListRequestFromJSON,
    RpcGetAccountProblemSubmissionSnippetListRequestFromJSONTyped,
    RpcGetAccountProblemSubmissionSnippetListRequestToJSON,
} from './RpcGetAccountProblemSubmissionSnippetListRequest';

/**
 * 
 * @export
 * @interface RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod
 */
export interface RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod
     */
    jsonrpc?: RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodJsonrpcEnum;
    /**
     * An identifier established by the Client.
     * @type {string}
     * @memberof RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod
     */
    id?: string;
    /**
     * A String containing the name of the method to be invoked.
     * @type {string}
     * @memberof RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod
     */
    method?: RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodMethodEnum;
    /**
     * 
     * @type {RpcGetAccountProblemSubmissionSnippetListRequest}
     * @memberof RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod
     */
    params?: RpcGetAccountProblemSubmissionSnippetListRequest;
}


/**
 * @export
 */
export const RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodJsonrpcEnum = typeof RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodJsonrpcEnum[keyof typeof RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodJsonrpcEnum];

/**
 * @export
 */
export const RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodMethodEnum = {
    GetAccountProblemSubmissionSnippetList: 'get_account_problem_submission_snippet_list'
} as const;
export type RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodMethodEnum = typeof RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodMethodEnum[keyof typeof RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodMethodEnum];


/**
 * Check if a given object implements the RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod interface.
 */
export function instanceOfRequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodFromJSON(json: any): RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod {
    return RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodFromJSONTyped(json, false);
}

export function RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'method': !exists(json, 'method') ? undefined : json['method'],
        'params': !exists(json, 'params') ? undefined : RpcGetAccountProblemSubmissionSnippetListRequestFromJSON(json['params']),
    };
}

export function RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethodToJSON(value?: RequestBodyOfTheGetAccountProblemSubmissionSnippetListMethod | null): any {
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
        'params': RpcGetAccountProblemSubmissionSnippetListRequestToJSON(value.params),
    };
}
