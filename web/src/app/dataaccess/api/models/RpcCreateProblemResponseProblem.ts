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

import type { RpcProblem } from './RpcProblem';
import {
    instanceOfRpcProblem,
    RpcProblemFromJSON,
    RpcProblemFromJSONTyped,
    RpcProblemToJSON,
} from './RpcProblem';

/**
 * @type RpcCreateProblemResponseProblem
 * 
 * @export
 */
export type RpcCreateProblemResponseProblem = RpcProblem;

export function RpcCreateProblemResponseProblemFromJSON(json: any): RpcCreateProblemResponseProblem {
    return RpcCreateProblemResponseProblemFromJSONTyped(json, false);
}

export function RpcCreateProblemResponseProblemFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcCreateProblemResponseProblem {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return { ...RpcProblemFromJSONTyped(json, true) };
}

export function RpcCreateProblemResponseProblemToJSON(value?: RpcCreateProblemResponseProblem | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }

    if (instanceOfRpcProblem(value)) {
        return RpcProblemToJSON(value as RpcProblem);
    }

    return {};
}
