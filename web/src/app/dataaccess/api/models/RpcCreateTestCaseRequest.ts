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
 * @interface RpcCreateTestCaseRequest
 */
export interface RpcCreateTestCaseRequest {
    /**
     * 
     * @type {string}
     * @memberof RpcCreateTestCaseRequest
     */
    problemUUID: string;
    /**
     * 
     * @type {string}
     * @memberof RpcCreateTestCaseRequest
     */
    input: string;
    /**
     * 
     * @type {string}
     * @memberof RpcCreateTestCaseRequest
     */
    output: string;
    /**
     * 
     * @type {boolean}
     * @memberof RpcCreateTestCaseRequest
     */
    isHidden: boolean;
}

/**
 * Check if a given object implements the RpcCreateTestCaseRequest interface.
 */
export function instanceOfRpcCreateTestCaseRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "problemUUID" in value;
    isInstance = isInstance && "input" in value;
    isInstance = isInstance && "output" in value;
    isInstance = isInstance && "isHidden" in value;

    return isInstance;
}

export function RpcCreateTestCaseRequestFromJSON(json: any): RpcCreateTestCaseRequest {
    return RpcCreateTestCaseRequestFromJSONTyped(json, false);
}

export function RpcCreateTestCaseRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcCreateTestCaseRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'problemUUID': json['ProblemUUID'],
        'input': json['Input'],
        'output': json['Output'],
        'isHidden': json['IsHidden'],
    };
}

export function RpcCreateTestCaseRequestToJSON(value?: RpcCreateTestCaseRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'ProblemUUID': value.problemUUID,
        'Input': value.input,
        'Output': value.output,
        'IsHidden': value.isHidden,
    };
}

