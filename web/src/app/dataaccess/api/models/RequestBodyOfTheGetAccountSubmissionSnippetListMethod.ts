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
import type { RpcGetAccountSubmissionSnippetListRequest } from './RpcGetAccountSubmissionSnippetListRequest';
import {
    RpcGetAccountSubmissionSnippetListRequestFromJSON,
    RpcGetAccountSubmissionSnippetListRequestFromJSONTyped,
    RpcGetAccountSubmissionSnippetListRequestToJSON,
} from './RpcGetAccountSubmissionSnippetListRequest';

/**
 * 
 * @export
 * @interface RequestBodyOfTheGetAccountSubmissionSnippetListMethod
 */
export interface RequestBodyOfTheGetAccountSubmissionSnippetListMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof RequestBodyOfTheGetAccountSubmissionSnippetListMethod
     */
    jsonrpc?: RequestBodyOfTheGetAccountSubmissionSnippetListMethodJsonrpcEnum;
    /**
     * An identifier established by the Client.
     * @type {string}
     * @memberof RequestBodyOfTheGetAccountSubmissionSnippetListMethod
     */
    id?: string;
    /**
     * A String containing the name of the method to be invoked.
     * @type {string}
     * @memberof RequestBodyOfTheGetAccountSubmissionSnippetListMethod
     */
    method?: RequestBodyOfTheGetAccountSubmissionSnippetListMethodMethodEnum;
    /**
     * 
     * @type {RpcGetAccountSubmissionSnippetListRequest}
     * @memberof RequestBodyOfTheGetAccountSubmissionSnippetListMethod
     */
    params?: RpcGetAccountSubmissionSnippetListRequest;
}


/**
 * @export
 */
export const RequestBodyOfTheGetAccountSubmissionSnippetListMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type RequestBodyOfTheGetAccountSubmissionSnippetListMethodJsonrpcEnum = typeof RequestBodyOfTheGetAccountSubmissionSnippetListMethodJsonrpcEnum[keyof typeof RequestBodyOfTheGetAccountSubmissionSnippetListMethodJsonrpcEnum];

/**
 * @export
 */
export const RequestBodyOfTheGetAccountSubmissionSnippetListMethodMethodEnum = {
    GetAccountSubmissionSnippetList: 'get_account_submission_snippet_list'
} as const;
export type RequestBodyOfTheGetAccountSubmissionSnippetListMethodMethodEnum = typeof RequestBodyOfTheGetAccountSubmissionSnippetListMethodMethodEnum[keyof typeof RequestBodyOfTheGetAccountSubmissionSnippetListMethodMethodEnum];


/**
 * Check if a given object implements the RequestBodyOfTheGetAccountSubmissionSnippetListMethod interface.
 */
export function instanceOfRequestBodyOfTheGetAccountSubmissionSnippetListMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RequestBodyOfTheGetAccountSubmissionSnippetListMethodFromJSON(json: any): RequestBodyOfTheGetAccountSubmissionSnippetListMethod {
    return RequestBodyOfTheGetAccountSubmissionSnippetListMethodFromJSONTyped(json, false);
}

export function RequestBodyOfTheGetAccountSubmissionSnippetListMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestBodyOfTheGetAccountSubmissionSnippetListMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'method': !exists(json, 'method') ? undefined : json['method'],
        'params': !exists(json, 'params') ? undefined : RpcGetAccountSubmissionSnippetListRequestFromJSON(json['params']),
    };
}

export function RequestBodyOfTheGetAccountSubmissionSnippetListMethodToJSON(value?: RequestBodyOfTheGetAccountSubmissionSnippetListMethod | null): any {
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
        'params': RpcGetAccountSubmissionSnippetListRequestToJSON(value.params),
    };
}

