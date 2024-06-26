/* tslint:disable */
/* eslint-disable */
/**
 * Bckupr
 * Bckupr API
 *
 * The version of the OpenAPI document: latest
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


/**
 * 
 * @export
 */
export const StopModes = {
    All: 'all',
    Labelled: 'labelled',
    Writers: 'writers',
    Attached: 'attached',
    Linked: 'linked'
} as const;
export type StopModes = typeof StopModes[keyof typeof StopModes];


export function instanceOfStopModes(value: any): boolean {
    return Object.values(StopModes).includes(value);
}

export function StopModesFromJSON(json: any): StopModes {
    return StopModesFromJSONTyped(json, false);
}

export function StopModesFromJSONTyped(json: any, ignoreDiscriminator: boolean): StopModes {
    return json as StopModes;
}

export function StopModesToJSON(value?: StopModes | null): any {
    return value as any;
}

