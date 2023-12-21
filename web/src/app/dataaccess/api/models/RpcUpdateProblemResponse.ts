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
import type { RpcCreateSubmissionResponseProblem } from './RpcCreateSubmissionResponseProblem';
import {
    RpcCreateSubmissionResponseProblemFromJSON,
    RpcCreateSubmissionResponseProblemFromJSONTyped,
    RpcCreateSubmissionResponseProblemToJSON,
} from './RpcCreateSubmissionResponseProblem';

/**
 * 
 * @export
 * @interface RpcUpdateProblemResponse
 */
export interface RpcUpdateProblemResponse {
    /**
     * 
     * @type {RpcCreateSubmissionResponseProblem}
     * @memberof RpcUpdateProblemResponse
     */
    problem: RpcCreateSubmissionResponseProblem;
}

/**
 * Check if a given object implements the RpcUpdateProblemResponse interface.
 */
export function instanceOfRpcUpdateProblemResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "problem" in value;

    return isInstance;
}

export function RpcUpdateProblemResponseFromJSON(json: any): RpcUpdateProblemResponse {
    return RpcUpdateProblemResponseFromJSONTyped(json, false);
}

export function RpcUpdateProblemResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcUpdateProblemResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'problem': RpcCreateSubmissionResponseProblemFromJSON(json['Problem']),
    };
}

export function RpcUpdateProblemResponseToJSON(value?: RpcUpdateProblemResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'Problem': RpcCreateSubmissionResponseProblemToJSON(value.problem),
    };
}

