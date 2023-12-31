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
/**
 * 
 * @export
 * @interface RpcGetAccountProblemSnippetListRequest
 */
export interface RpcGetAccountProblemSnippetListRequest {
    /**
     * 
     * @type {number}
     * @memberof RpcGetAccountProblemSnippetListRequest
     */
    accountID: number;
    /**
     * 
     * @type {number}
     * @memberof RpcGetAccountProblemSnippetListRequest
     */
    offset: number;
    /**
     * 
     * @type {number}
     * @memberof RpcGetAccountProblemSnippetListRequest
     */
    limit: number;
}

/**
 * Check if a given object implements the RpcGetAccountProblemSnippetListRequest interface.
 */
export function instanceOfRpcGetAccountProblemSnippetListRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "accountID" in value;
    isInstance = isInstance && "offset" in value;
    isInstance = isInstance && "limit" in value;

    return isInstance;
}

export function RpcGetAccountProblemSnippetListRequestFromJSON(json: any): RpcGetAccountProblemSnippetListRequest {
    return RpcGetAccountProblemSnippetListRequestFromJSONTyped(json, false);
}

export function RpcGetAccountProblemSnippetListRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcGetAccountProblemSnippetListRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'accountID': json['AccountID'],
        'offset': json['Offset'],
        'limit': json['Limit'],
    };
}

export function RpcGetAccountProblemSnippetListRequestToJSON(value?: RpcGetAccountProblemSnippetListRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'AccountID': value.accountID,
        'Offset': value.offset,
        'Limit': value.limit,
    };
}

