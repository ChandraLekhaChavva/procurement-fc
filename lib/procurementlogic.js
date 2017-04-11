/*
 * IBM Confidential
 * OCO Source Materials
 * Fabric Composer
 * Copyright IBM Corp. 2016
 * The source code for this program is not published or otherwise
 * divested of its trade secrets, irrespective of what has
 * been deposited with the U.S. Copyright Office.
 */

'use strict';

/**
 * Close the bidding for a vehicle listing and choose the
 * highest bid that is over the asking price
 * @param {ibm.itpeople.procurement.UpdateContractorInfo} updateContractorInfo - the updateContractor transaction
 * @transaction
 */
function UpdateContractorInfo(updateContractorInfo) {
    var tempSupplierContractRegistry;
    return getAssetRegistry('ibm.itpeople.procurement.SupplierContract')
        .then(function(supplierContractRegistry) {
            tempSupplierContractRegistry = supplierContractRegistry;
            return supplierContractRegistry.get(updateContractorInfo.PONumber);
        }).then(function(supplierContract) {
            supplierContract.ContractorInfo = updateContractorInfo.ContractorInfo;
            return tempSupplierContractRegistry.update(supplierContract);
        });
}