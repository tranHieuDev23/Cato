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
 * @interface RpcCreateTestCaseListRequest
 */
export interface RpcCreateTestCaseListRequest {
    /**
     * 
     * @type {string}
     * @memberof RpcCreateTestCaseListRequest
     */
    problemUUID: string;
    /**
     * 
     * @type {string}
     * @memberof RpcCreateTestCaseListRequest
     */
    zippedTestData: string;
}

/**
 * Check if a given object implements the RpcCreateTestCaseListRequest interface.
 */
export function instanceOfRpcCreateTestCaseListRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "problemUUID" in value;
    isInstance = isInstance && "zippedTestData" in value;

    return isInstance;
}

export function RpcCreateTestCaseListRequestFromJSON(json: any): RpcCreateTestCaseListRequest {
    return RpcCreateTestCaseListRequestFromJSONTyped(json, false);
}

export function RpcCreateTestCaseListRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcCreateTestCaseListRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'problemUUID': json['ProblemUUID'],
        'zippedTestData': json['ZippedTestData'],
    };
}

export function RpcCreateTestCaseListRequestToJSON(value?: RpcCreateTestCaseListRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'ProblemUUID': value.problemUUID,
        'ZippedTestData': value.zippedTestData,
    };
}

