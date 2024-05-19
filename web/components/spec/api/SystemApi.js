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


import ApiClient from "../ApiClient";
import Version from '../model/Version';

/**
* System service.
* @module api/SystemApi
* @version latest
*/
export default class SystemApi {

    /**
    * Constructs a new SystemApi. 
    * @alias module:api/SystemApi
    * @class
    * @param {module:ApiClient} [apiClient] Optional API client implementation to use,
    * default to {@link module:ApiClient#instance} if unspecified.
    */
    constructor(apiClient) {
        this.apiClient = apiClient || ApiClient.instance;
    }


    /**
     * Callback function to receive the result of the getVersion operation.
     * @callback module:api/SystemApi~getVersionCallback
     * @param {String} error Error message, if any.
     * @param {module:model/Version} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Retrieves application version
     * @param {module:api/SystemApi~getVersionCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/Version}
     */
    getVersion(callback) {
      let postBody = null;

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = [];
      let accepts = ['application/json'];
      let returnType = Version;
      return this.apiClient.callApi(
        '/api/version', 'GET',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null, callback
      );
    }


}