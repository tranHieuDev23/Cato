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
 * @interface RpcGetProblemTestCaseSnippetListRequest
 */
export interface RpcGetProblemTestCaseSnippetListRequest {
    /**
     * 
     * @type {string}
     * @memberof RpcGetProblemTestCaseSnippetListRequest
     */
    problemUUID: string;
    /**
     * 
     * @type {number}
     * @memberof RpcGetProblemTestCaseSnippetListRequest
     */
    offset: number;
    /**
     * 
     * @type {number}
     * @memberof RpcGetProblemTestCaseSnippetListRequest
     */
    limit: number;
}

/**
 * Check if a given object implements the RpcGetProblemTestCaseSnippetListRequest interface.
 */
export function instanceOfRpcGetProblemTestCaseSnippetListRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "problemUUID" in value;
    isInstance = isInstance && "offset" in value;
    isInstance = isInstance && "limit" in value;

    return isInstance;
}

export function RpcGetProblemTestCaseSnippetListRequestFromJSON(json: any): RpcGetProblemTestCaseSnippetListRequest {
    return RpcGetProblemTestCaseSnippetListRequestFromJSONTyped(json, false);
}

export function RpcGetProblemTestCaseSnippetListRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcGetProblemTestCaseSnippetListRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'problemUUID': json['ProblemUUID'],
        'offset': json['Offset'],
        'limit': json['Limit'],
    };
}

export function RpcGetProblemTestCaseSnippetListRequestToJSON(value?: RpcGetProblemTestCaseSnippetListRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'ProblemUUID': value.problemUUID,
        'Offset': value.offset,
        'Limit': value.limit,
    };
}

