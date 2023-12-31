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
 * @interface RpcUpdateSubmissionRequest
 */
export interface RpcUpdateSubmissionRequest {
    /**
     * 
     * @type {number}
     * @memberof RpcUpdateSubmissionRequest
     */
    iD: number;
    /**
     * 
     * @type {number}
     * @memberof RpcUpdateSubmissionRequest
     */
    status: number;
    /**
     * 
     * @type {number}
     * @memberof RpcUpdateSubmissionRequest
     */
    result: number;
}

/**
 * Check if a given object implements the RpcUpdateSubmissionRequest interface.
 */
export function instanceOfRpcUpdateSubmissionRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "iD" in value;
    isInstance = isInstance && "status" in value;
    isInstance = isInstance && "result" in value;

    return isInstance;
}

export function RpcUpdateSubmissionRequestFromJSON(json: any): RpcUpdateSubmissionRequest {
    return RpcUpdateSubmissionRequestFromJSONTyped(json, false);
}

export function RpcUpdateSubmissionRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcUpdateSubmissionRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'iD': json['ID'],
        'status': json['Status'],
        'result': json['Result'],
    };
}

export function RpcUpdateSubmissionRequestToJSON(value?: RpcUpdateSubmissionRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'ID': value.iD,
        'Status': value.status,
        'Result': value.result,
    };
}

