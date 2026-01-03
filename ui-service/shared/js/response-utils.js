// === BAR RESTAURANT - RESPONSE UTILITIES ===

/**
 * Check if an API response is successful
 * Handles the centralized response structure: { code, message, data }
 * @param {Response} response - Fetch API response object
 * @returns {Promise<{success: boolean, result: object}>}
 */
async function isSuccessfulResponse(response) {
    try {
        const result = await response.json();
        
        // Check for success based on HTTP status and response code
        const success = response.ok && (result.code === 200 || result.code === 201);
        
        return { success, result };
    } catch (error) {
        console.error('❌ Error parsing response:', error);
        return { 
            success: false, 
            result: { 
                code: response.status, 
                message: 'Failed to parse response' 
            } 
        };
    }
}

/**
 * Extract error message from API response
 * @param {object} result - Parsed JSON response
 * @returns {string}
 */
function getErrorMessage(result) {
    return result.message || result.error || 'An unexpected error occurred';
}

/**
 * Extract data from successful API response
 * @param {object} result - Parsed JSON response
 * @returns {object}
 */
function getResponseData(result) {
    return result.data || result;
}

// Export for global access
window.isSuccessfulResponse = isSuccessfulResponse;
window.getErrorMessage = getErrorMessage;
window.getResponseData = getResponseData;
