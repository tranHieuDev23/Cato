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
 * @interface RpcCreateSubmissionRequest
 */
export interface RpcCreateSubmissionRequest {
    /**
     * 
     * @type {string}
     * @memberof RpcCreateSubmissionRequest
     */
    problemUUID: string;
    /**
     * 
     * @type {string}
     * @memberof RpcCreateSubmissionRequest
     */
    content: string;
    /**
     * 
     * @type {string}
     * @memberof RpcCreateSubmissionRequest
     */
    language: string;
}

/**
 * Check if a given object implements the RpcCreateSubmissionRequest interface.
 */
export function instanceOfRpcCreateSubmissionRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "problemUUID" in value;
    isInstance = isInstance && "content" in value;
    isInstance = isInstance && "language" in value;

    return isInstance;
}

export function RpcCreateSubmissionRequestFromJSON(json: any): RpcCreateSubmissionRequest {
    return RpcCreateSubmissionRequestFromJSONTyped(json, false);
}

export function RpcCreateSubmissionRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcCreateSubmissionRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'problemUUID': json['ProblemUUID'],
        'content': json['Content'],
        'language': json['Language'],
    };
}

export function RpcCreateSubmissionRequestToJSON(value?: RpcCreateSubmissionRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'ProblemUUID': value.problemUUID,
        'Content': value.content,
        'Language': value.language,
    };
}

