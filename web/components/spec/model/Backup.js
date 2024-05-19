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
 *
 */

import ApiClient from '../ApiClient';
import Status from './Status';
import Volume from './Volume';

/**
 * The Backup model module.
 * @module model/Backup
 * @version latest
 */
class Backup {
    /**
     * Constructs a new <code>Backup</code>.
     * @alias module:model/Backup
     * @param id {String} 
     * @param created {Date} 
     * @param type {String} 
     * @param status {module:model/Status} 
     * @param volumes {Array.<module:model/Volume>} 
     */
    constructor(id, created, type, status, volumes) { 
        
        Backup.initialize(this, id, created, type, status, volumes);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, id, created, type, status, volumes) { 
        obj['id'] = id;
        obj['created'] = created;
        obj['type'] = type;
        obj['status'] = status;
        obj['volumes'] = volumes;
    }

    /**
     * Constructs a <code>Backup</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/Backup} obj Optional instance to populate.
     * @return {module:model/Backup} The populated <code>Backup</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new Backup();

            if (data.hasOwnProperty('id')) {
                obj['id'] = ApiClient.convertToType(data['id'], 'String');
            }
            if (data.hasOwnProperty('created')) {
                obj['created'] = ApiClient.convertToType(data['created'], 'Date');
            }
            if (data.hasOwnProperty('type')) {
                obj['type'] = ApiClient.convertToType(data['type'], 'String');
            }
            if (data.hasOwnProperty('status')) {
                obj['status'] = Status.constructFromObject(data['status']);
            }
            if (data.hasOwnProperty('error')) {
                obj['error'] = ApiClient.convertToType(data['error'], 'String');
            }
            if (data.hasOwnProperty('volumes')) {
                obj['volumes'] = ApiClient.convertToType(data['volumes'], [Volume]);
            }
        }
        return obj;
    }

    /**
     * Validates the JSON data with respect to <code>Backup</code>.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @return {boolean} to indicate whether the JSON data is valid with respect to <code>Backup</code>.
     */
    static validateJSON(data) {
        // check to make sure all required properties are present in the JSON string
        for (const property of Backup.RequiredProperties) {
            if (!data.hasOwnProperty(property)) {
                throw new Error("The required field `" + property + "` is not found in the JSON data: " + JSON.stringify(data));
            }
        }
        // ensure the json data is a string
        if (data['id'] && !(typeof data['id'] === 'string' || data['id'] instanceof String)) {
            throw new Error("Expected the field `id` to be a primitive type in the JSON string but got " + data['id']);
        }
        // ensure the json data is a string
        if (data['type'] && !(typeof data['type'] === 'string' || data['type'] instanceof String)) {
            throw new Error("Expected the field `type` to be a primitive type in the JSON string but got " + data['type']);
        }
        // ensure the json data is a string
        if (data['error'] && !(typeof data['error'] === 'string' || data['error'] instanceof String)) {
            throw new Error("Expected the field `error` to be a primitive type in the JSON string but got " + data['error']);
        }
        if (data['volumes']) { // data not null
            // ensure the json data is an array
            if (!Array.isArray(data['volumes'])) {
                throw new Error("Expected the field `volumes` to be an array in the JSON data but got " + data['volumes']);
            }
            // validate the optional field `volumes` (array)
            for (const item of data['volumes']) {
                Volume.validateJSON(item);
            };
        }

        return true;
    }


}

Backup.RequiredProperties = ["id", "created", "type", "status", "volumes"];

/**
 * @member {String} id
 */
Backup.prototype['id'] = undefined;

/**
 * @member {Date} created
 */
Backup.prototype['created'] = undefined;

/**
 * @member {String} type
 */
Backup.prototype['type'] = undefined;

/**
 * @member {module:model/Status} status
 */
Backup.prototype['status'] = undefined;

/**
 * @member {String} error
 */
Backup.prototype['error'] = undefined;

/**
 * @member {Array.<module:model/Volume>} volumes
 */
Backup.prototype['volumes'] = undefined;






export default Backup;

