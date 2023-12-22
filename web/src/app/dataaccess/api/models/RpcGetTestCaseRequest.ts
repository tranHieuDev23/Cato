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
 * @interface RpcGetTestCaseRequest
 */
export interface RpcGetTestCaseRequest {
    /**
     * 
     * @type {number}
     * @memberof RpcGetTestCaseRequest
     */
    iD: number;
}

/**
 * Check if a given object implements the RpcGetTestCaseRequest interface.
 */
export function instanceOfRpcGetTestCaseRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "iD" in value;

    return isInstance;
}

export function RpcGetTestCaseRequestFromJSON(json: any): RpcGetTestCaseRequest {
    return RpcGetTestCaseRequestFromJSONTyped(json, false);
}

export function RpcGetTestCaseRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcGetTestCaseRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'iD': json['ID'],
    };
}

export function RpcGetTestCaseRequestToJSON(value?: RpcGetTestCaseRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'ID': value.iD,
    };
}

