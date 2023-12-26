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
import type { RpcGetAndUpdateFirstSubmittedSubmissionAsExecutingResponseSubmission } from './RpcGetAndUpdateFirstSubmittedSubmissionAsExecutingResponseSubmission';
import {
    RpcGetAndUpdateFirstSubmittedSubmissionAsExecutingResponseSubmissionFromJSON,
    RpcGetAndUpdateFirstSubmittedSubmissionAsExecutingResponseSubmissionFromJSONTyped,
    RpcGetAndUpdateFirstSubmittedSubmissionAsExecutingResponseSubmissionToJSON,
} from './RpcGetAndUpdateFirstSubmittedSubmissionAsExecutingResponseSubmission';

/**
 * 
 * @export
 * @interface RpcGetSubmissionResponse
 */
export interface RpcGetSubmissionResponse {
    /**
     * 
     * @type {RpcGetAndUpdateFirstSubmittedSubmissionAsExecutingResponseSubmission}
     * @memberof RpcGetSubmissionResponse
     */
    submission: RpcGetAndUpdateFirstSubmittedSubmissionAsExecutingResponseSubmission;
}

/**
 * Check if a given object implements the RpcGetSubmissionResponse interface.
 */
export function instanceOfRpcGetSubmissionResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "submission" in value;

    return isInstance;
}

export function RpcGetSubmissionResponseFromJSON(json: any): RpcGetSubmissionResponse {
    return RpcGetSubmissionResponseFromJSONTyped(json, false);
}

export function RpcGetSubmissionResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcGetSubmissionResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'submission': RpcGetAndUpdateFirstSubmittedSubmissionAsExecutingResponseSubmissionFromJSON(json['Submission']),
    };
}

export function RpcGetSubmissionResponseToJSON(value?: RpcGetSubmissionResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'Submission': RpcGetAndUpdateFirstSubmittedSubmissionAsExecutingResponseSubmissionToJSON(value.submission),
    };
}

