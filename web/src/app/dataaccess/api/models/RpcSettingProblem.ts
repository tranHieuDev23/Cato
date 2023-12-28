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

import type { RpcProblemSetting } from './RpcProblemSetting';
import {
    instanceOfRpcProblemSetting,
    RpcProblemSettingFromJSON,
    RpcProblemSettingFromJSONTyped,
    RpcProblemSettingToJSON,
} from './RpcProblemSetting';

/**
 * @type RpcSettingProblem
 * 
 * @export
 */
export type RpcSettingProblem = RpcProblemSetting;

export function RpcSettingProblemFromJSON(json: any): RpcSettingProblem {
    return RpcSettingProblemFromJSONTyped(json, false);
}

export function RpcSettingProblemFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcSettingProblem {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return { ...RpcProblemSettingFromJSONTyped(json, true) };
}

export function RpcSettingProblemToJSON(value?: RpcSettingProblem | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }

    if (instanceOfRpcProblemSetting(value)) {
        return RpcProblemSettingToJSON(value as RpcProblemSetting);
    }

    return {};
}
