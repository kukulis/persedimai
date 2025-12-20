// State variables (need to persist across function calls)
let debounceTimer;
let destinationDebounceTimer;

// Initialize database bounds display
function initializeDatabaseBounds() {
    const databaseSelect = document.getElementById(ID_DATABASE);

    // Load bounds on page load if database is selected
    if (databaseSelect.value) {
        fetchAndDisplayBounds(databaseSelect.value);
        fetchAndDisplayTimeBounds(databaseSelect.value);
    }

    // Load bounds when database changes
    databaseSelect.addEventListener('change', (e) => {
        fetchAndDisplayBounds(e.target.value);
        fetchAndDisplayTimeBounds(e.target.value);
    });
}

// Initialize arrival time validation
function initializeArrivalTimeValidation() {
    const arrivalFromInput = document.getElementById('arrival_from');
    const arrivalFromError = document.getElementById('arrival_from_error');
    setupDateTimeValidation(arrivalFromInput, arrivalFromError);

    const arrivalToInput = document.getElementById('arrival_to');
    const arrivalToError = document.getElementById('arrival_to_error');
    setupDateTimeValidation(arrivalToInput, arrivalToError);
}

// Initialize source point search functionality
function initializeSourcePointSearch() {
    // Clear source point selection
    document.getElementById('source-clear').addEventListener('click', () => {
        const sourceHiddenInput = document.getElementById('source');
        const sourceSearchInput = document.getElementById('source-search');

        sourceHiddenInput.value = '';
        sourceSearchInput.value = '';
        document.getElementById('source-selected').classList.remove('show');
    });

    // Handle search input
    document.getElementById('source-search').addEventListener('input', (e) => {
        const sourceHiddenInput = document.getElementById('source');
        const sourceDropdown = document.getElementById('source-dropdown');
        const value = e.target.value;

        // Clear selection when user types
        sourceHiddenInput.value = '';

        // Clear previous timer
        clearTimeout(debounceTimer);

        if (!value.trim()) {
            sourceDropdown.classList.remove('show');
            return;
        }

        // Show loading
        sourceDropdown.innerHTML = '<div class="loading">Loading...</div>';
        sourceDropdown.classList.add('show');

        // Debounce API call
        debounceTimer = setTimeout(async () => {
            const parsed = parseSearchValue(value);
            if (parsed) {
                const databaseSelect = document.getElementById(ID_DATABASE);
                parsed.database = databaseSelect.value;

                const url = buildApiUrl(parsed);
                const points = await fetchPoints(url);
                displayPoints(points);
            }
        }, 300);
    });

    // Close dropdown when clicking outside
    document.addEventListener('click', (e) => {
        const sourceSearchInput = document.getElementById('source-search');
        const sourceDropdown = document.getElementById('source-dropdown');

        if (!sourceSearchInput.contains(e.target) && !sourceDropdown.contains(e.target)) {
            sourceDropdown.classList.remove('show');
        }
    });

    // Show dropdown when focusing on search input
    document.getElementById('source-search').addEventListener('focus', () => {
        const sourceDropdown = document.getElementById('source-dropdown');

        if (sourceDropdown.children.length > 0) {
            sourceDropdown.classList.add('show');
        }
    });
}

// Initialize destination point search functionality
function initializeDestinationPointSearch() {
    // Clear destination point selection
    document.getElementById('destination-clear').addEventListener('click', () => {
        const destinationHiddenInput = document.getElementById('destination');
        const destinationSearchInput = document.getElementById('destination-search');

        destinationHiddenInput.value = '';
        destinationSearchInput.value = '';
        document.getElementById('destination-selected').classList.remove('show');
    });

    // Handle destination search input
    document.getElementById('destination-search').addEventListener('input', (e) => {
        const destinationHiddenInput = document.getElementById('destination');
        const destinationDropdown = document.getElementById('destination-dropdown');
        const value = e.target.value;

        // Clear selection when user types
        destinationHiddenInput.value = '';

        // Clear previous timer
        clearTimeout(destinationDebounceTimer);

        if (!value.trim()) {
            destinationDropdown.classList.remove('show');
            return;
        }

        // Show loading
        destinationDropdown.innerHTML = '<div class="loading">Loading...</div>';
        destinationDropdown.classList.add('show');

        // Debounce API call
        destinationDebounceTimer = setTimeout(async () => {
            const parsed = parseSearchValue(value);
            if (parsed) {
                const databaseSelect = document.getElementById(ID_DATABASE);
                parsed.database = databaseSelect.value;

                const url = buildApiUrl(parsed);
                const points = await fetchPoints(url);
                displayDestinationPoints(points);
            }
        }, 300);
    });

    // Close destination dropdown when clicking outside
    document.addEventListener('click', (e) => {
        const destinationSearchInput = document.getElementById('destination-search');
        const destinationDropdown = document.getElementById('destination-dropdown');

        if (!destinationSearchInput.contains(e.target) && !destinationDropdown.contains(e.target)) {
            destinationDropdown.classList.remove('show');
        }
    });

    // Show destination dropdown when focusing on search input
    document.getElementById('destination-search').addEventListener('focus', () => {
        const destinationDropdown = document.getElementById('destination-dropdown');

        if (destinationDropdown.children.length > 0) {
            destinationDropdown.classList.add('show');
        }
    });
}

// Initialize form submission validation
function initializeFormValidation() {
    document.getElementById('travel-search-form').addEventListener('submit', (e) => {
        const sourceHiddenInput = document.getElementById('source');
        const sourceSearchInput = document.getElementById('source-search');
        const sourceError = document.getElementById('source_error');
        const destinationHiddenInput = document.getElementById('destination');
        const destinationSearchInput = document.getElementById('destination-search');
        const destinationError = document.getElementById('destination_error');
        const arrivalFromInput = document.getElementById('arrival_from');
        const arrivalFromError = document.getElementById('arrival_from_error');
        const arrivalToInput = document.getElementById('arrival_to');
        const arrivalToError = document.getElementById('arrival_to_error');

        let hasError = false;
        let firstErrorField = null;

        // Validate source point
        if (!sourceHiddenInput.value) {
            sourceSearchInput.classList.add('input-error');
            sourceError.textContent = 'Please select a source point from the dropdown';
            hasError = true;
            if (!firstErrorField) firstErrorField = sourceSearchInput;
        } else {
            sourceSearchInput.classList.remove('input-error');
            sourceError.textContent = '';
        }

        // Validate destination point
        if (!destinationHiddenInput.value) {
            destinationSearchInput.classList.add('input-error');
            destinationError.textContent = 'Please select a destination point from the dropdown';
            hasError = true;
            if (!firstErrorField) firstErrorField = destinationSearchInput;
        } else {
            destinationSearchInput.classList.remove('input-error');
            destinationError.textContent = '';
        }

        // Validate arrival_from
        const resultFrom = validateDateTimeFormat(arrivalFromInput.value);
        if (!resultFrom.valid) {
            arrivalFromInput.classList.add('input-error');
            arrivalFromError.textContent = resultFrom.message;
            hasError = true;
            if (!firstErrorField) firstErrorField = arrivalFromInput;
        }

        // Validate arrival_to
        const resultTo = validateDateTimeFormat(arrivalToInput.value);
        if (!resultTo.valid) {
            arrivalToInput.classList.add('input-error');
            arrivalToError.textContent = resultTo.message;
            hasError = true;
            if (!firstErrorField) firstErrorField = arrivalToInput;
        }

        if (hasError) {
            e.preventDefault();
            // Focus on first error field
            if (firstErrorField) {
                firstErrorField.focus();
            }
            return false;
        }
    });
}

// Main initialization function
function initializeTravelSearchForm() {
    initializeDatabaseBounds();
    initializeArrivalTimeValidation();
    initializeSourcePointSearch();
    initializeDestinationPointSearch();
    initializeFormValidation();

    // Initialize pre-selected points
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initializePreselectedPoints);
    } else {
        initializePreselectedPoints();
    }
}
