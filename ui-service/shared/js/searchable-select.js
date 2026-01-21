/**
 * Searchable Select Utility
 * Wraps Choices.js for easy searchable dropdown management
 */

window.SearchableSelect = (function() {
    // Store all Choices instances by element ID
    const instances = new Map();

    // Default configuration
    const defaultConfig = {
        searchEnabled: true,
        searchPlaceholderValue: 'Type to search...',
        itemSelectText: '',
        shouldSort: false,
        searchResultLimit: 50,
        noResultsText: 'No results found',
        noChoicesText: 'No options available',
        removeItemButton: false,
        allowHTML: false,
        classNames: {
            containerOuter: 'choices',
            containerInner: 'choices__inner',
            input: 'choices__input',
            inputCloned: 'choices__input--cloned',
            list: 'choices__list',
            listItems: 'choices__list--multiple',
            listSingle: 'choices__list--single',
            listDropdown: 'choices__list--dropdown',
            item: 'choices__item',
            itemSelectable: 'choices__item--selectable',
            itemDisabled: 'choices__item--disabled',
            itemChoice: 'choices__item--choice',
            placeholder: 'choices__placeholder',
            group: 'choices__group',
            groupHeading: 'choices__heading',
            button: 'choices__button',
            activeState: 'is-active',
            focusState: 'is-focused',
            openState: 'is-open',
            disabledState: 'is-disabled',
            highlightedState: 'is-highlighted',
            selectedState: 'is-selected',
            flippedState: 'is-flipped',
            loadingState: 'is-loading',
            noResults: 'has-no-results',
            noChoices: 'has-no-choices'
        }
    };

    /**
     * Initialize a searchable select on an element
     * @param {string|HTMLElement} selector - CSS selector or DOM element
     * @param {Object} config - Optional Choices.js configuration
     * @returns {Choices|null} The Choices instance or null if failed
     */
    function init(selector, config = {}) {
        const element = typeof selector === 'string' 
            ? document.querySelector(selector) 
            : selector;

        if (!element) {
            console.warn('SearchableSelect: Element not found:', selector);
            return null;
        }

        // Destroy existing instance if any
        destroy(element);

        // Merge configs
        const finalConfig = { ...defaultConfig, ...config };

        try {
            const instance = new Choices(element, finalConfig);
            const id = element.id || element.name || Math.random().toString(36).substr(2, 9);
            instances.set(id, { element, instance });
            return instance;
        } catch (error) {
            console.error('SearchableSelect: Failed to initialize:', error);
            return null;
        }
    }

    /**
     * Initialize all selects with data-searchable attribute
     * @param {HTMLElement} container - Container to search within (default: document)
     * @param {Object} config - Optional Choices.js configuration
     */
    function initAll(container = document, config = {}) {
        const selects = container.querySelectorAll('select[data-searchable]');
        selects.forEach(select => init(select, config));
    }

    /**
     * Destroy a Choices instance
     * @param {string|HTMLElement} selector - CSS selector or DOM element
     */
    function destroy(selector) {
        const element = typeof selector === 'string' 
            ? document.querySelector(selector) 
            : selector;

        if (!element) return;

        const id = element.id || element.name;
        
        // Check by ID first
        if (id && instances.has(id)) {
            try {
                instances.get(id).instance.destroy();
            } catch (e) {
                // Instance might already be destroyed
            }
            instances.delete(id);
            return;
        }

        // Check by element reference
        for (const [key, value] of instances) {
            if (value.element === element) {
                try {
                    value.instance.destroy();
                } catch (e) {
                    // Instance might already be destroyed
                }
                instances.delete(key);
                return;
            }
        }
    }

    /**
     * Destroy all Choices instances
     */
    function destroyAll() {
        instances.forEach(({ instance }) => {
            try {
                instance.destroy();
            } catch (e) {
                // Instance might already be destroyed
            }
        });
        instances.clear();
    }

    /**
     * Get a Choices instance by selector
     * @param {string|HTMLElement} selector - CSS selector or DOM element
     * @returns {Choices|null}
     */
    function getInstance(selector) {
        const element = typeof selector === 'string' 
            ? document.querySelector(selector) 
            : selector;

        if (!element) return null;

        const id = element.id || element.name;
        if (id && instances.has(id)) {
            return instances.get(id).instance;
        }

        for (const [key, value] of instances) {
            if (value.element === element) {
                return value.instance;
            }
        }

        return null;
    }

    /**
     * Refresh/reinitialize a select after changing its options
     * @param {string|HTMLElement} selector - CSS selector or DOM element
     * @param {Object} config - Optional Choices.js configuration
     * @returns {Choices|null}
     */
    function refresh(selector, config = {}) {
        return init(selector, config);
    }

    /**
     * Set choices programmatically
     * @param {string|HTMLElement} selector - CSS selector or DOM element
     * @param {Array} choices - Array of {value, label, selected, disabled} objects
     * @param {string} value - Value to select after setting choices
     */
    function setChoices(selector, choices, value = '') {
        const instance = getInstance(selector);
        if (instance) {
            instance.clearStore();
            instance.setChoices(choices, 'value', 'label', true);
            if (value) {
                instance.setChoiceByValue(value);
            }
        }
    }

    /**
     * Set the selected value
     * @param {string|HTMLElement} selector - CSS selector or DOM element
     * @param {string} value - Value to select
     */
    function setValue(selector, value) {
        const instance = getInstance(selector);
        if (instance) {
            instance.setChoiceByValue(value);
        }
    }

    /**
     * Clear the selection
     * @param {string|HTMLElement} selector - CSS selector or DOM element
     */
    function clearValue(selector) {
        const instance = getInstance(selector);
        if (instance) {
            instance.removeActiveItems();
        }
    }

    /**
     * Enable/disable a select
     * @param {string|HTMLElement} selector - CSS selector or DOM element
     * @param {boolean} enabled - Enable or disable
     */
    function setEnabled(selector, enabled) {
        const instance = getInstance(selector);
        if (instance) {
            if (enabled) {
                instance.enable();
            } else {
                instance.disable();
            }
        }
    }

    return {
        init,
        initAll,
        destroy,
        destroyAll,
        getInstance,
        refresh,
        setChoices,
        setValue,
        clearValue,
        setEnabled
    };
})();
