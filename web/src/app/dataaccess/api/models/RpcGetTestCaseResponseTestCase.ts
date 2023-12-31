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

import type { RpcTestCase } from './RpcTestCase';
import {
    instanceOfRpcTestCase,
    RpcTestCaseFromJSON,
    RpcTestCaseFromJSONTyped,
    RpcTestCaseToJSON,
} from './RpcTestCase';

/**
 * @type RpcGetTestCaseResponseTestCase
 * 
 * @export
 */
export type RpcGetTestCaseResponseTestCase = RpcTestCase;

export function RpcGetTestCaseResponseTestCaseFromJSON(json: any): RpcGetTestCaseResponseTestCase {
    return RpcGetTestCaseResponseTestCaseFromJSONTyped(json, false);
}

export function RpcGetTestCaseResponseTestCaseFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcGetTestCaseResponseTestCase {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return { ...RpcTestCaseFromJSONTyped(json, true) };
}

export function RpcGetTestCaseResponseTestCaseToJSON(value?: RpcGetTestCaseResponseTestCase | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }

    if (instanceOfRpcTestCase(value)) {
        return RpcTestCaseToJSON(value as RpcTestCase);
    }

    return {};
}

